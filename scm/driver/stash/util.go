// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"net/url"
	"strconv"

	"github.com/drone/go-scm/scm"
)

func encodeListOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page > 1 {
		params.Set("start", strconv.Itoa(opts.Page*opts.Size))
	}
	if opts.Size != 0 {
		params.Set("limit", strconv.Itoa(opts.Size))
	}
	return params.Encode()
}

func encodeListRoleOptions(opts scm.ListOptions) string {
	params := url.Values{}
	if opts.Page > 1 {
		params.Set("start", strconv.Itoa(opts.Page*opts.Size))
	}
	if opts.Size != 0 {
		params.Set("size", strconv.Itoa(opts.Size))
	}
	params.Set("role", "member")
	return params.Encode()
}

func encodeCommitListOptions(opts scm.CommitListOptions) string {
	params := url.Values{}
	if opts.Page > 1 {
		params.Set("start", strconv.Itoa(opts.Page*opts.Size))
	}
	if opts.Size != 0 {
		params.Set("size", strconv.Itoa(opts.Size))
	}
	return params.Encode()
}

func encodeIssueListOptions(opts scm.IssueListOptions) string {
	params := url.Values{}
	if opts.Page > 1 {
		params.Set("start", strconv.Itoa(opts.Page*opts.Size))
	}
	if opts.Size != 0 {
		params.Set("size", strconv.Itoa(opts.Size))
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
	if opts.Page > 1 {
		params.Set("start", strconv.Itoa(opts.Page*opts.Size))
	}
	if opts.Size != 0 {
		params.Set("size", strconv.Itoa(opts.Size))
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
	uri, err := url.Parse(from.Next)
	if err != nil {
		return err
	}
	page := uri.Query().Get("page")
	to.Page.First = 1
	to.Page.Next, _ = strconv.Atoi(page)
	return nil
}
