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

// fetchAllWorkspaces fetches all workspaces for the authenticated user
func (c *wrapper) fetchAllWorkspaces(ctx context.Context) ([]string, error) {
	var workspaceSlugs []string
	page := 1
	pageLen := 100
	for {
		path := fmt.Sprintf("2.0/user/workspaces?page=%d&pageLen=%d", page, pageLen)
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

// fetchReposFromAllWorkspaces fetches repositories from all user workspaces.
// It continues processing other workspaces even if one fails, returning partial results.
// An error is only returned if fetching workspaces fails or if no repositories could be fetched at all.
func (c *wrapper) fetchReposFromAllWorkspaces(ctx context.Context, queryParams string) ([]*scm.Repository, *scm.Response, error) {
	workspaces, err := c.fetchAllWorkspaces(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(workspaces) == 0 {
		return []*scm.Repository{}, nil, nil
	}

	var (
		allRepos []*scm.Repository
		lastRes  *scm.Response
	)

	for _, workspaceSlug := range workspaces {

		path := fmt.Sprintf("2.0/repositories/%s?%s", workspaceSlug, queryParams)

		for path != "" {
			out := new(repositories)

			res, err := c.do(ctx, "GET", path, nil, &out)
			if err != nil {
				break
			}

			allRepos = append(allRepos, convertRepositoryList(out)...)
			lastRes = res

			path = out.Next
		}
	}

	return allRepos, lastRes, nil
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
