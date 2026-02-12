// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"net/url"
	"testing"

	"github.com/drone/go-scm/scm"
)

func Test_encodeListOptions(t *testing.T) {
	opts := scm.ListOptions{
		Page: 10,
		Size: 30,
	}
	want := "page=10&pagelen=30"
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
	want := "page=10&pagelen=30"
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
	want := "page=10&pagelen=30&state=all"
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
	want := "page=10&pagelen=30&state=all"
	got := encodePullRequestListOptions(opts)
	if got != want {
		t.Errorf("Want encoded pr list options %q, got %q", want, got)
	}
}

func Test_copyPagination(t *testing.T) {
	tests := []struct {
		from pagination
		want scm.Page
	}{
		{
			from: pagination{},
			want: scm.Page{},
		},
		{
			from: pagination{
				Next: "https://api.bitbucket.org/2.0/teams?pagelen=10",
			},
			want: scm.Page{
				Next: 0,
			},
		},
		{
			from: pagination{
				Next: "https://api.bitbucket.org/2.0/teams?pagelen=10&page=2",
			},
			want: scm.Page{
				Next: 2,
			},
		},
	}
	for _, test := range tests {
		res := &scm.Response{}
		err := copyPagination(test.from, res)
		if err != nil {
			t.Error(err)
		}
		if got, want := res.Page.Next, test.want.Next; got != want {
			t.Errorf("Want Next page %d, got %d", want, got)
		}
	}
}

func Test_extractWorkspaceFromURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "default api url",
			url:  "https://api.bitbucket.org/",
			want: "",
		},
		{
			name: "api url without trailing slash",
			url:  "https://api.bitbucket.org",
			want: "",
		},
		{
			name: "workspace in repositories path",
			url:  "https://api.bitbucket.org/2.0/repositories/my-workspace/",
			want: "my-workspace",
		},
		{
			name: "workspace in repositories path without trailing slash",
			url:  "https://api.bitbucket.org/2.0/repositories/my-workspace",
			want: "my-workspace",
		},
		{
			name: "workspace as last path segment",
			url:  "https://api.bitbucket.org/my-workspace/",
			want: "my-workspace",
		},
		{
			name: "only 2.0 in path",
			url:  "https://api.bitbucket.org/2.0/",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsed, _ := url.Parse(test.url)
			client := &wrapper{&scm.Client{BaseURL: parsed}}
			got := client.extractWorkspaceFromURL()
			if got != test.want {
				t.Errorf("Want workspace %q, got %q", test.want, got)
			}
		})
	}
}
