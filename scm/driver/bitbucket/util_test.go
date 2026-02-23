// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"net/url"
	"strings"
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

func Test_encodeListRoleOptions(t *testing.T) {
	opts := scm.ListOptions{
		Page: 10,
		Size: 30,
	}
	want := "page=10&pagelen=30&role=member"
	got := encodeListRoleOptions(opts)
	if got != want {
		t.Errorf("Want encoded list role options %q, got %q", want, got)
	}
}

func Test_encodeRepoListOptions(t *testing.T) {
	tests := []struct {
		name string
		opts scm.RepoListOptions
		want string
	}{
		{
			name: "with repo search term",
			opts: scm.RepoListOptions{
				RepoSearchTerm: scm.RepoSearchTerm{RepoName: "test-repo"},
				ListOptions:    scm.ListOptions{Page: 1, Size: 10},
			},
			want: "page=1&pagelen=10&q=name~%22test-repo%22&role=member",
		},
		{
			name: "without search term",
			opts: scm.RepoListOptions{
				ListOptions: scm.ListOptions{Page: 2, Size: 20},
			},
			want: "page=2&pagelen=20&role=member",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeRepoListOptions(tt.opts)
			if got != tt.want {
				t.Errorf("Want encoded repo list options %q, got %q", tt.want, got)
			}
		})
	}
}

func Test_encodeBranchListOptions(t *testing.T) {
	tests := []struct {
		name string
		opts scm.BranchListOptions
		want string
	}{
		{
			name: "with search term",
			opts: scm.BranchListOptions{
				SearchTerm:      "feature",
				PageListOptions: scm.ListOptions{Page: 1, Size: 10},
			},
			want: "page=1&pagelen=10&q=name~%22feature%22",
		},
		{
			name: "without search term",
			opts: scm.BranchListOptions{
				PageListOptions: scm.ListOptions{Page: 2, Size: 20},
			},
			want: "page=2&pagelen=20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeBranchListOptions(tt.opts)
			if got != tt.want {
				t.Errorf("Want encoded branch list options %q, got %q", tt.want, got)
			}
		})
	}
}

func Test_extractWorkspaceFromURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
	}{
		{
			name:    "URL with repositories path",
			baseURL: "https://api.bitbucket.org/repositories/my-workspace",
			want:    "my-workspace",
		},
		{
			name:    "URL with repositories and repo",
			baseURL: "https://api.bitbucket.org/repositories/my-workspace/my-repo",
			want:    "my-workspace",
		},
		{
			name:    "URL with trailing slash",
			baseURL: "https://api.bitbucket.org/my-workspace/",
			want:    "my-workspace",
		},
		{
			name:    "URL ending with 2.0",
			baseURL: "https://api.bitbucket.org/2.0",
			want:    "",
		},
		{
			name:    "Empty path",
			baseURL: "https://api.bitbucket.org",
			want:    "",
		},
		{
			name:    "URL with workspace slug at end",
			baseURL: "https://api.bitbucket.org/some-workspace",
			want:    "some-workspace",
		},
		{
			name:    "URL ending with user (excluded segment)",
			baseURL: "https://api.bitbucket.org/2.0/user",
			want:    "",
		},
		{
			name:    "URL ending with workspaces (excluded segment)",
			baseURL: "https://api.bitbucket.org/2.0/workspaces",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("Failed to parse URL %q: %v", tt.baseURL, err)
			}
			if !strings.HasSuffix(u.Path, "/") {
				u.Path = u.Path + "/"
			}
			w := &wrapper{
				Client:    &scm.Client{BaseURL: u},
				workspace: "",
				repo:      "",
			}
			got := w.extractWorkspaceFromURL()
			if got != tt.want {
				t.Errorf("extractWorkspaceFromURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_extractWorkspaceFromURL_WithSetWorkspace(t *testing.T) {
	tests := []struct {
		name             string
		baseURL          string
		explicitWorkspace string
		want             string
	}{
		{
			name:             "Explicit workspace set, should use it instead of URL",
			baseURL:          "https://api.bitbucket.org",
			explicitWorkspace: "harness-io",
			want:             "harness-io",
		},
		{
			name:             "Explicit workspace overrides URL workspace",
			baseURL:          "https://api.bitbucket.org/repositories/old-workspace",
			explicitWorkspace: "new-workspace",
			want:             "new-workspace",
		},
		{
			name:             "No explicit workspace, falls back to URL parsing",
			baseURL:          "https://api.bitbucket.org/repositories/url-workspace",
			explicitWorkspace: "",
			want:             "url-workspace",
		},
		{
			name:             "Explicit workspace with standard Bitbucket Cloud URL",
			baseURL:          "https://api.bitbucket.org/2.0",
			explicitWorkspace: "my-workspace",
			want:             "my-workspace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("Failed to parse URL %q: %v", tt.baseURL, err)
			}
			if !strings.HasSuffix(u.Path, "/") {
				u.Path = u.Path + "/"
			}
			w := &wrapper{
				Client:    &scm.Client{BaseURL: u},
				workspace: tt.explicitWorkspace,
				repo:      "",
			}
			got := w.extractWorkspaceFromURL()
			if got != tt.want {
				t.Errorf("extractWorkspaceFromURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
