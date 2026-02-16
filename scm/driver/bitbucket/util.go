// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/drone/go-scm/scm"
)

// regex for git author fields ("name <name@mail.tld>")
var reGitMail = regexp.MustCompile("<(.*)>")

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
// Workspaces are returned in sorted order by slug to ensure deterministic pagination behavior.
func (c *wrapper) fetchAllWorkspaces(ctx context.Context) ([]string, error) {
	var workspaceSlugs []string
	page := 1
	pageLen := 100
	for {
		// Sort by workspace.slug to ensure consistent ordering across requests
		// This is critical for deterministic pagination across workspaces
		path := fmt.Sprintf("2.0/user/workspaces?page=%d&pagelen=%d&sort=workspace.slug", page, pageLen)
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
//   globalOffset = 100
//   1. Get ws1 count → 150. cumulative=150. 150 > 100 → ws1 contains our start.
//      localOffset in ws1 = 100. need = min(150-100, 100) = 50.
//      Fetch 50 repos from ws1 at offset 100. remaining = 50.
//   2. Get ws2 count → 120. cumulative=270. Need 50 more from ws2.
//      localOffset in ws2 = 0. need = min(120, 50) = 50.
//      Fetch 50 repos from ws2 at offset 0. remaining = 0. Done.
//   hasMore = cumulative(270) > 100+100 → true → Next=3.
func (c *wrapper) fetchReposWithPagination(ctx context.Context, queryParams string, page, size int) ([]*scm.Repository, *scm.Response, error) {
	workspaces, err := c.fetchAllWorkspaces(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(workspaces) == 0 {
		return []*scm.Repository{}, &scm.Response{}, nil
	}

	if size == 0 {
		size = 100
	}
	if page == 0 {
		page = 1
	}

	globalOffset := (page - 1) * size

	var (
		result     []*scm.Repository
		cumulative int  // running total of repos seen across workspaces
		remaining  = size
	)

	for _, workspaceSlug := range workspaces {
		if remaining <= 0 {
			// Page is already full.
			// We still need to know if more repos exist (for hasMore).
			// One cheap count call is enough to answer that.
			wsCount, err := c.getWorkspaceRepoCount(ctx, workspaceSlug, queryParams)
			if err == nil && wsCount > 0 {
				cumulative += wsCount
			}
			break // we only need one more workspace to confirm hasMore
		}

		// Get this workspace's total repo count (single API call using Bitbucket's "size" field)
		wsCount, err := c.getWorkspaceRepoCount(ctx, workspaceSlug, queryParams)
		if err != nil {
			// Skip inaccessible workspace, continue with next
			continue
		}

		wsStart := cumulative
		wsEnd := cumulative + wsCount
		cumulative = wsEnd

		// This workspace ends before our target range — skip it entirely, no repo fetch needed
		if wsEnd <= globalOffset {
			continue
		}

		// Calculate how many repos we need from this workspace
		localOffset := 0
		if globalOffset > wsStart {
			localOffset = globalOffset - wsStart
		}
		need := wsCount - localOffset
		if need > remaining {
			need = remaining
		}

		// Fetch only the slice we need using Bitbucket's page/pagelen math
		repos, err := c.fetchReposFromWorkspaceWithOffset(ctx, workspaceSlug, queryParams, localOffset, need)
		if err != nil {
			continue
		}

		result = append(result, repos...)
		remaining -= len(repos)
	}

	// Build response
	res := &scm.Response{
		Page: scm.Page{
			First: 1,
		},
	}

	// There is a next page if total repos across workspaces exceed current page's end
	if cumulative > globalOffset+size {
		res.Page.Next = page + 1
	}

	return result, res, nil
}

// getWorkspaceRepoCount returns the total number of repositories in a workspace.
// It makes a single lightweight API call (pagelen=1) and reads the "size" field
// from Bitbucket's response. Falls back to counting pages if "size" is absent.
func (c *wrapper) getWorkspaceRepoCount(ctx context.Context, workspaceSlug, queryParams string) (int, error) {
	// Use pagelen=1 to minimise payload — we only need the "size" field
	params, _ := url.ParseQuery(queryParams)
	params.Set("pagelen", "1")
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
	// Resolve the pagelen that Bitbucket will use for this workspace query
	pagelen := 100
	if params, err := url.ParseQuery(queryParams); err == nil {
		if v := params.Get("pagelen"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				pagelen = n
			}
		}
	}

	// Translate logical offset into Bitbucket page number + skip count
	bbPage := (offset / pagelen) + 1
	skip := offset % pagelen

	var result []*scm.Repository
	collected := 0

	for collected < limit {
		params, _ := url.ParseQuery(queryParams)
		params.Set("page", strconv.Itoa(bbPage))
		path := fmt.Sprintf("2.0/repositories/%s?%s", workspaceSlug, params.Encode())

		out := new(repositories)
		_, err := c.do(ctx, "GET", path, nil, &out)
		if err != nil {
			return result, err
		}

		repos := convertRepositoryList(out)
		if len(repos) == 0 {
			break
		}

		// On the first fetched page, skip items before our logical offset
		start := 0
		if bbPage == (offset/pagelen)+1 && skip > 0 {
			start = skip
			if start >= len(repos) {
				bbPage++
				continue
			}
		}

		// Take only what we still need
		end := len(repos)
		if end-start > limit-collected {
			end = start + (limit - collected)
		}

		result = append(result, repos[start:end]...)
		collected += end - start

		if out.Next == "" {
			break // workspace exhausted
		}
		bbPage++
	}

	return result, nil
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
