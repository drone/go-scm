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

// extractWorkspaceFromURL attempts to extract workspace from the client's BaseURL path.
// Expected URL patterns:
// - https://api.bitbucket.org/2.0/repositories/{workspace}/
// - https://api.bitbucket.org/{workspace}/
// Returns empty string if no workspace found in URL.
func (c *wrapper) extractWorkspaceFromURL() string {
	path := strings.Trim(c.BaseURL.Path, "/")
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")

	// Look for "repositories/{workspace}" pattern
	for i, part := range parts {
		if part == "repositories" && i+1 < len(parts) && parts[i+1] != "" {
			return parts[i+1]
		}
	}

	// If path has segments and last one is not "2.0", it might be workspace
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		if last != "" && last != "2.0" {
			return last
		}
	}

	return ""
}

// workspaceAccess represents a single workspace access entry from /2.0/user/workspaces
type workspaceAccess struct {
	Type          string     `json:"type"`
	Administrator bool       `json:"administrator"`
	Workspace     *workspace `json:"workspace"`
}

// workspace represents the nested workspace object in /2.0/user/workspaces response
type workspace struct {
	Type  string `json:"type"`
	UUID  string `json:"uuid"`
	Slug  string `json:"slug"`
	Links struct {
		Avatar link `json:"avatar"`
		Self   link `json:"self"`
		HTML   link `json:"html"`
	} `json:"links"`
}

// workspaceAccessList represents the paginated response from /2.0/user/workspaces
type workspaceAccessList struct {
	pagination
	Values []*workspaceAccess `json:"values"`
}

// fetchAllWorkspaces fetches all workspaces for the authenticated user
// using the /2.0/user/workspaces endpoint.
func (c *wrapper) fetchAllWorkspaces(ctx context.Context) ([]string, error) {
	var workspaceSlugs []string
	page := 1

	for {
		path := fmt.Sprintf("2.0/user/workspaces?page=%d&pagelen=100", page)
		out := new(workspaceAccessList)
		_, err := c.do(ctx, "GET", path, nil, out)
		if err != nil {
			return nil, err
		}

		for _, wa := range out.Values {
			if wa.Workspace != nil && wa.Workspace.Slug != "" {
				workspaceSlugs = append(workspaceSlugs, wa.Workspace.Slug)
			}
		}

		// Check if there are more pages
		if out.Next == "" {
			break
		}
		page++
	}

	return workspaceSlugs, nil
}

// fetchReposFromAllWorkspaces fetches repositories from all user workspaces.
// This is used by List() and ListV2() to aggregate repositories across all
// workspaces after the /2.0/repositories deprecation.
func (c *wrapper) fetchReposFromAllWorkspaces(ctx context.Context, queryParams string) ([]*scm.Repository, *scm.Response, error) {
	// Step 1: Fetch all workspaces
	workspaces, err := c.fetchAllWorkspaces(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Step 2: Aggregate repositories from all workspaces
	var allRepos []*scm.Repository
	var lastRes *scm.Response

	for _, ws := range workspaces {
		path := fmt.Sprintf("2.0/repositories/%s?%s", ws, queryParams)
		out := new(repositories)
		res, err := c.do(ctx, "GET", path, nil, &out)
		if err != nil {
			// Continue to next workspace on error (user may not have access)
			continue
		}
		allRepos = append(allRepos, convertRepositoryList(out)...)
		lastRes = res
	}

	return allRepos, lastRes, nil
}
