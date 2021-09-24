// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/drone/go-scm/scm"
)

func encodeListOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("limit", strconv.Itoa(opts.Size))
	}
	return params.Encode()
}

func encodeIssueListOptions(opts scm.IssueListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("limit", strconv.Itoa(opts.Size))
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
		params.Set("limit", strconv.Itoa(opts.Size))
	}
	if opts.Open && opts.Closed {
		params.Set("state", "all")
	} else if opts.Closed {
		params.Set("state", "closed")
	}
	return params.Encode()
}

// ConvertAPIURLToHTMLURL converts an release API endpoint into a html endpoint
func ConvertAPIURLToHTMLURL(apiURL string, tagName string) string {
	// "url": "https://try.gitea.com/api/v1/repos/octocat/Hello-World/123",
	// "html_url": "https://try.gitea.com/octocat/Hello-World/releases/tag/v1.0.0",
	// the url field is the API url, not the html url, so until go-sdk v0.13.3, build it ourselves
	link, err := url.Parse(apiURL)
	if err != nil {
		return ""
	}

	pathParts := strings.Split(link.Path, "/")
	if len(pathParts) != 7 {
		return ""
	}
	link.Path = fmt.Sprintf("/%s/%s/releases/tag/%s", pathParts[4], pathParts[5], tagName)
	return link.String()
}

func encodeMilestoneListOptions(opts scm.MilestoneListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("per_page", strconv.Itoa(opts.Size))
	}
	if opts.Open && opts.Closed {
		params.Set("state", "all")
	} else if opts.Closed {
		params.Set("state", "closed")
	}
	return params.Encode()
}

func encodeReleaseListOptions(opts scm.ReleaseListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("per_page", strconv.Itoa(opts.Size))
	}
	if opts.Open && opts.Closed {
		params.Set("state", "all")
	} else if opts.Closed {
		params.Set("state", "closed")
	}
	return params.Encode()
}