// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"testing"

	"github.com/drone/go-scm/scm"
)

func Test_encodeListOptions(t *testing.T) {
	opts := scm.ListOptions{
		Page: 10,
		Size: 30,
	}
	want := "page=10&per_page=30"
	got := encodeListOptions(opts)
	if got != want {
		t.Errorf("Want encoded list options %q, got %q", want, got)
	}
}

func Test_encodeCommitListOptions(t *testing.T) {
	opts := scm.CommitListOptions{
		Page: 10,
		Size: 30,
		Ref:  "master",
	}
	want := "page=10&per_page=30&sha=master"
	got := encodeCommitListOptions(opts)
	if got != want {
		t.Errorf("Want encoded commit list options %q, got %q", want, got)
	}
}

func Test_encodeIssueListOptions(t *testing.T) {
	opts := scm.IssueListOptions{
		Page:   10,
		Size:   30,
		Open:   true,
		Closed: true,
	}
	want := "page=10&per_page=30&state=all"
	got := encodeIssueListOptions(opts)
	if got != want {
		t.Errorf("Want encoded issue list options %q, got %q", want, got)
	}
}

func Test_encodeIssueListOptions_Closed(t *testing.T) {
	opts := scm.IssueListOptions{
		Page:   10,
		Size:   30,
		Open:   false,
		Closed: true,
	}
	want := "page=10&per_page=30&state=closed"
	got := encodeIssueListOptions(opts)
	if got != want {
		t.Errorf("Want encoded issue list options %q, got %q", want, got)
	}
}

func Test_encodePullRequestListOptions(t *testing.T) {
	t.Parallel()
	opts := scm.PullRequestListOptions{
		Page:   10,
		Size:   30,
		Open:   true,
		Closed: true,
	}
	want := "page=10&per_page=30&state=all"
	got := encodePullRequestListOptions(opts)
	if got != want {
		t.Errorf("Want encoded pr list options %q, got %q", want, got)
	}
}

func Test_encodePullRequestListOptions_Closed(t *testing.T) {
	t.Parallel()
	opts := scm.PullRequestListOptions{
		Page:   10,
		Size:   30,
		Open:   false,
		Closed: true,
	}
	want := "page=10&per_page=30&state=closed"
	got := encodePullRequestListOptions(opts)
	if got != want {
		t.Errorf("Want encoded pr list options %q, got %q", want, got)
	}
}
