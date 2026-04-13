// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/drone/go-scm/scm"
)

// regex for git author fields ("name <name@mail.tld>")
var reGitMail = regexp.MustCompile("<(.*)>")

// isHardError reports whether err should abort pagination entirely rather than
// just skipping the current workspace.
//
// Hard errors (stop everything):
//   - 429 rate limit: the API is throttled for ALL requests; retrying other workspaces makes it worse
//   - 401 unauthorized: credentials are invalid; no workspace will succeed
//
// Soft errors (skip this workspace, continue):
//   - 403 access denied: no permission for this specific workspace, others may work
//   - 404 not found: workspace gone, others may still exist
func isHardError(err error) bool {
	if err == scm.ErrNotAuthorized {
		return true // 401: bad credentials, all workspaces will fail
	}
	if bbErr, ok := err.(*Error); ok {
		return bbErr.StatusCode == 429 // rate limit: stop immediately
	}
	return false
}

// extracts the email from a git commit author string
func extractEmail(gitauthor string) (author string) {
	matches := reGitMail.FindAllStringSubmatch(gitauthor, -1)
	if len(matches) == 1 {
		author = matches[0][1]
	}
	return
}

func encodeBranchListOptions(opts scm.BranchListOptions) string {
	params := url.Values{}
	if opts.SearchTerm != "" {
		var sb strings.Builder
		sb.WriteString("name~\"")
		sb.WriteString(opts.SearchTerm)
		sb.WriteString("\"")
		params.Set("q", sb.String())
	}
	if opts.PageListOptions != (scm.ListOptions{}) {
		if opts.PageListOptions.Page != 0 {
			params.Set("page", strconv.Itoa(opts.PageListOptions.Page))
		}
		if opts.PageListOptions.Size != 0 {
			params.Set("pagelen", strconv.Itoa(opts.PageListOptions.Size))
		}
	}
	return params.Encode()
}

func encodeListOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("pagelen", strconv.Itoa(opts.Size))
	}
	return params.Encode()
}

func encodeListRoleOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("pagelen", strconv.Itoa(opts.Size))
	}
	params.Set("role", "member")
	return params.Encode()
}

func encodeRepoListOptions(opts scm.RepoListOptions) string {
	params := url.Values{}
	if opts.RepoSearchTerm.RepoName != "" {
		var sb strings.Builder
		sb.WriteString("name~\"")
		sb.WriteString(opts.RepoSearchTerm.RepoName)
		sb.WriteString("\"")
		params.Set("q", sb.String())
	}
	if opts.ListOptions != (scm.ListOptions{}) {
		if opts.ListOptions.Page != 0 {
			params.Set("page", strconv.Itoa(opts.ListOptions.Page))
		}
		if opts.ListOptions.Size != 0 {
			params.Set("pagelen", strconv.Itoa(opts.ListOptions.Size))
		}
	}
	params.Set("role", "member")
	return params.Encode()
}

func encodeCommitListOptions(opts scm.CommitListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("pagelen", strconv.Itoa(opts.Size))
	}
	if opts.Path != "" {
		params.Set("path", opts.Path)
	}
	return params.Encode()
}

func encodeIssueListOptions(opts scm.IssueListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("pagelen", strconv.Itoa(opts.Size))
	}
	if opts.Open && opts.Closed {
		params.Set("state", "all")
	} else if opts.Closed {
		params.Set("state", "closed")
	}
	return params.Encode()
}

func encodePullRequestListOptions(opts scm.PullRequestListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("pagelen", strconv.Itoa(opts.Size))
	}
	if opts.Open && opts.Closed {
		params.Set("state", "all")
	} else if opts.Closed {
		params.Set("state", "closed")
	}
	return params.Encode()
}

func copyPagination(from pagination, to *scm.Response) error {
	if to == nil {
		return nil
	}
	to.Page.NextURL = from.Next
	uri, err := url.Parse(from.Next)
	if err != nil {
		return err
	}
	page := uri.Query().Get("page")
	to.Page.First = 1
	to.Page.Next, _ = strconv.Atoi(page)
	return nil
}

// workspaceAccessList represents the paginated response from /2.0/user/workspaces
type workspaceAccessList struct {
	pagination
	Values []*workspaceAccess `json:"values"`
}

// workspaceAccess represents a single workspace access entry from /2.0/user/workspaces
type workspaceAccess struct {
	Workspace *workspace `json:"workspace"`
}

// workspace represents the nested workspace object in /2.0/user/workspaces response
type workspace struct {
	Slug  string `json:"slug"`
	Links struct {
		Avatar link `json:"avatar"`
	} `json:"links"`
}

// fetchAllWorkspaces fetches all workspaces for the authenticated user.
// Workspaces are sorted by slug after fetching to ensure deterministic
// pagination behavior across calls.
func (c *wrapper) fetchAllWorkspaces(ctx context.Context) ([]string, error) {
	var workspaceSlugs []string
	page := 1
	pageLen := 100
	for {
		path := fmt.Sprintf("2.0/user/workspaces?page=%d&pagelen=%d", page, pageLen)
		workspaceStruct := new(workspaceAccessList)
		_, err := c.do(ctx, "GET", path, nil, workspaceStruct)
		if err != nil {
			return nil, err
		}

		for _, workspaceAccessResponse := range workspaceStruct.Values {
			if workspaceAccessResponse.Workspace != nil && workspaceAccessResponse.Workspace.Slug != "" {
				workspaceSlugs = append(workspaceSlugs, workspaceAccessResponse.Workspace.Slug)
			}
		}

		// Check if there are more pages
		if workspaceStruct.Next == "" {
			break
		}
		page++
	}
	sort.Strings(workspaceSlugs)

	return workspaceSlugs, nil
}

// fetchReposWithPagination fetches repositories with workspace-aware pagination.
// Returns exactly 'size' repos (or fewer on the last page) by stitching across workspaces.
//
// The algorithm is lazy: it only queries a workspace when needed and stops
// the moment the requested page is filled. Nothing is cached between calls;
// determinism comes from always sorting workspaces by slug.
//
// Walk-through (Page 2, Size 100, ws1 has 150 repos, ws2 has 120):
//
//	globalOffset = 100
//	1. Get ws1 count → 150. cumulative=150. 150 > 100 → ws1 contains our start.
//	   localOffset in ws1 = 100. need = min(150-100, 100) = 50.
//	   Fetch 50 repos from ws1 at offset 100. remaining = 50.
//	2. Get ws2 count → 120. cumulative=270. Need 50 more from ws2.
//	   localOffset in ws2 = 0. need = min(120, 50) = 50.
//	   Fetch 50 repos from ws2 at offset 0. remaining = 0. Done.
//	hasMore = cumulative(270) > 100+100 → true → Next=3.
func (c *wrapper) fetchReposWithPagination(ctx context.Context, queryParams string, page, size int) ([]*scm.Repository, *scm.Response, error) {
	workspaces, err := c.fetchAllWorkspaces(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(workspaces) == 0 {
		return []*scm.Repository{}, &scm.Response{}, nil
	}

	// Normalize page and size with sensible defaults
	if size == 0 {
		size = 100
	}
	if page == 0 {
		page = 1
	}

	globalOffset := (page - 1) * size
	paginationState := &paginationState{
		result:     make([]*scm.Repository, 0),
		cumulative: 0,
		remaining:  size,
		hasMore:    false,
	}

	for i, workspaceSlug := range workspaces {
		if !c.processPaginationWorkspace(ctx, queryParams, globalOffset, i, len(workspaces), workspaceSlug, paginationState) {
			break
		}
	}

	if paginationState.err != nil {
		return nil, nil, paginationState.err
	}

	res := &scm.Response{
		Page: scm.Page{
			First: 1,
		},
	}
	if paginationState.hasMore {
		res.Page.Next = page + 1
	}

	return paginationState.result, res, nil
}

// paginationState tracks the progress of pagination across workspaces.
type paginationState struct {
	result     []*scm.Repository
	cumulative int   // running total of repos seen across workspaces
	remaining  int   // repos still needed to fill the page
	hasMore    bool  // whether more repos exist beyond this page
	err        error // first hard error (e.g. 429) encountered; stops pagination
}

// processPaginationWorkspace processes a single workspace in pagination.
// Returns false if pagination is complete, true to continue with next workspace.
func (c *wrapper) processPaginationWorkspace(ctx context.Context, queryParams string, globalOffset, workspaceIndex, totalWorkspaces int, workspaceSlug string, state *paginationState) bool {
	if state.remaining <= 0 {
		return c.checkHasMoreRepos(ctx, workspaceSlug, queryParams, state)
	}

	wsCount, err := c.getWorkspaceRepoCount(ctx, workspaceSlug, queryParams)
	if err != nil {
		if isHardError(err) {
			state.err = err // stop pagination: rate-limit or bad credentials
			return false
		}
		return true // soft error (e.g. 403 on one workspace): skip, try next
	}

	wsStart := state.cumulative
	state.cumulative += wsCount

	if state.cumulative <= globalOffset {
		return true
	}

	localOffset := c.calculateLocalOffset(globalOffset, wsStart)
	need := c.calculateReposNeeded(wsCount, localOffset, state.remaining)

	repos, err := c.fetchReposFromWorkspaceWithOffset(ctx, workspaceSlug, queryParams, localOffset, need)
	if err != nil {
		if isHardError(err) {
			state.err = err // stop pagination: rate-limit or bad credentials
			return false
		}
		return true // soft error: skip, try next workspace
	}

	state.result = append(state.result, repos...)
	state.remaining -= len(repos)

	return c.determineHasMoreAndContinue(wsCount, localOffset, need, workspaceIndex, totalWorkspaces, state)
}

// checkHasMoreRepos checks if there are more repos when pagination is complete.
func (c *wrapper) checkHasMoreRepos(ctx context.Context, workspaceSlug, queryParams string, state *paginationState) bool {
	wsCount, err := c.getWorkspaceRepoCount(ctx, workspaceSlug, queryParams)
	if err == nil && wsCount > 0 {
		state.hasMore = true
	}
	return false
}

// calculateLocalOffset returns the starting offset within a workspace.
func (c *wrapper) calculateLocalOffset(globalOffset, wsStart int) int {
	if globalOffset > wsStart {
		return globalOffset - wsStart
	}
	return 0
}

// calculateReposNeeded returns how many repos are needed from this workspace.
func (c *wrapper) calculateReposNeeded(wsCount, localOffset, remaining int) int {
	need := wsCount - localOffset
	if need > remaining {
		return remaining
	}
	return need
}

// determineHasMoreAndContinue determines if more repos exist and whether to continue.
func (c *wrapper) determineHasMoreAndContinue(wsCount, localOffset, need, workspaceIndex, totalWorkspaces int, state *paginationState) bool {
	if state.remaining <= 0 {
		leftoverInWs := wsCount - (localOffset + need)
		if leftoverInWs > 0 {
			state.hasMore = true
			return false
		}
		if workspaceIndex+1 < totalWorkspaces {
			return true
		}
		return false
	}
	return true
}

// getWorkspaceRepoCount returns the total number of repositories in a workspace.
// It makes a single lightweight API call (pagelen=1) and reads the "size" field
// from Bitbucket's response. Falls back to counting pages if "size" is absent.
func (c *wrapper) getWorkspaceRepoCount(ctx context.Context, workspaceSlug, queryParams string) (int, error) {
	// Use pagelen=1 to minimise payload — we only need the "size" field.
	params, _ := url.ParseQuery(queryParams)
	params.Set("pagelen", "1")
	params.Set("page", "1")
	path := fmt.Sprintf("2.0/repositories/%s?%s", workspaceSlug, params.Encode())

	out := new(repositories)
	_, err := c.do(ctx, "GET", path, nil, &out)
	if err != nil {
		return 0, err
	}

	// Bitbucket returns total count in "size"; prefer it when available
	if out.Size > 0 {
		return out.Size, nil
	}

	// Fallback: Bitbucket didn't provide size, count by paging through
	count := len(out.Values)
	next := out.Next
	for next != "" {
		page := new(repositories)
		_, err := c.do(ctx, "GET", next, nil, &page)
		if err != nil {
			return count, nil
		}
		count += len(page.Values)
		next = page.Next
	}
	return count, nil
}

// fetchReposFromWorkspaceWithOffset fetches `limit` repos from a workspace
// starting at the given logical offset. It translates the offset into
// Bitbucket page numbers and skips leading items on the first page.
//
// Example: offset=150, pagelen=100  →  Bitbucket page 2, skip first 50 items.
func (c *wrapper) fetchReposFromWorkspaceWithOffset(ctx context.Context, workspaceSlug, queryParams string, offset, limit int) ([]*scm.Repository, error) {
	pagelen := extractPagelen(queryParams)
	bbPage := (offset / pagelen) + 1
	skip := offset % pagelen

	var result []*scm.Repository
	collected := 0
	isFirstPage := true

	for collected < limit {
		repos, nextURL, err := c.fetchReposPageFromWorkspace(ctx, workspaceSlug, queryParams, bbPage)
		if err != nil {
			return result, err
		}

		if len(repos) == 0 {
			break
		}

		start, shouldContinue := c.getSliceStartIdxIfWithinLimits(repos, skip, isFirstPage)
		if shouldContinue {
			bbPage++
			isFirstPage = false
			continue
		}

		end := c.getSliceEndIdx(len(repos), start, limit, collected)
		result = append(result, repos[start:end]...)
		collected += end - start
		isFirstPage = false

		if nextURL == "" {
			break
		}
		bbPage++
	}

	return result, nil
}

// getSliceStartIdxIfWithinLimits returns the starting index and whether to skip to next page.
func (c *wrapper) getSliceStartIdxIfWithinLimits(repos []*scm.Repository, skip int, isFirstPage bool) (int, bool) {
	start := 0
	if isFirstPage && skip > 0 {
		start = skip
		if start >= len(repos) {
			return start, true
		}
	}
	return start, false
}

// getSliceEndIdx calculates the ending index for the slice to take.
func (c *wrapper) getSliceEndIdx(reposLen, start, limit, collected int) int {
	end := reposLen
	if end-start > limit-collected {
		end = start + (limit - collected)
	}
	return end
}

// extractPagelen extracts the page length from query parameters, defaulting to 100.
func extractPagelen(queryParams string) int {
	pagelen := 100
	if params, err := url.ParseQuery(queryParams); err == nil {
		if v := params.Get("pagelen"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				pagelen = n
			}
		}
	}
	return pagelen
}

// fetchReposPageFromWorkspace fetches a single page of repositories from a workspace.
func (c *wrapper) fetchReposPageFromWorkspace(ctx context.Context, workspaceSlug, queryParams string, page int) ([]*scm.Repository, string, error) {
	params, _ := url.ParseQuery(queryParams)
	params.Set("page", strconv.Itoa(page))
	path := fmt.Sprintf("2.0/repositories/%s?%s", workspaceSlug, params.Encode())

	out := new(repositories)
	_, err := c.do(ctx, "GET", path, nil, &out)
	if err != nil {
		return nil, "", err
	}

	return convertRepositoryList(out), out.Next, nil
}

// extractWorkspaceFromURL attempts to extract workspace from the client's BaseURL path.
func (c *wrapper) extractWorkspaceFromURL() string {
	path := strings.Trim(c.BaseURL.Path, "/")
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")

	for i, part := range parts {
		if part == "repositories" && i+1 < len(parts) && parts[i+1] != "" {
			return parts[i+1]
		}
	}

	if len(parts) > 0 {
		last := parts[len(parts)-1]
		excludedSegments := map[string]bool{
			"2.0":        true,
			"user":       true,
			"workspaces": true,
			"teams":      true,
		}
		if last != "" && !excludedSegments[last] {
			return last
		}
	}

	return ""
}
