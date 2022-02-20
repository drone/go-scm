// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"net/url"
	"strconv"

	"github.com/drone/go-scm/scm"
)

func encodeListOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("per_page", strconv.Itoa(opts.Size))
	}
	return params.Encode()
}

func encodeCommitListOptions(opts scm.CommitListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("per_page", strconv.Itoa(opts.Size))
	}
	if opts.Ref != "" {
		params.Set("sha", opts.Ref)
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
		params.Set("per_page", strconv.Itoa(opts.Size))
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
		params.Set("per_page", strconv.Itoa(opts.Size))
	}
	if opts.Open && opts.Closed {
		params.Set("state", "all")
	} else if opts.Closed {
		params.Set("state", "closed")
	}
	return params.Encode()
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
