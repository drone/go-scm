// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/bitbucket/fixtures"
)

func TestGitFindCommit(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Git.FindCommit(context.Background(), "atlassian/stash-example-plugin", "a6e5e7d797edf751cbd839d6bd4aef86c941eec9")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Fields", testCommit(result))
}

func TestGitFindBranch(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Git.FindBranch(context.Background(), "atlassian/stash-example-plugin", "master")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Fields", testBranch(result))
}

func TestGitFindTag(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Git.FindTag(context.Background(), "atlassian/atlaskit", "@atlaskit/activity@1.0.3")
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := result.Name, "@atlaskit/activity@1.0.3"; got != want {
		t.Errorf("Want tag Name %q, got %q", want, got)
	}
	if got, want := result.Sha, "ceb01356c3f062579bdfeb15bc53fe151b9e00f0"; got != want {
		t.Errorf("Want tag Sha %q, got %q", want, got)
	}
}

func TestGitListCommits(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Git.ListCommits(context.Background(), "atlassian/stash-example-plugin", scm.CommitListOptions{Ref: "master"})
	if err != nil {
		t.Error(err)
		return
	}
	if len(result) == 0 {
		t.Errorf("Want non-empty commit list")
		return
	}
	t.Run("Page", testPage(res))
	t.Run("Fields", testCommit(result[0]))
}

func TestGitListBranches(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()
	client, _ := New(server.URL)
	result, res, err := client.Git.ListBranches(context.Background(), "atlassian/stash-example-plugin", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if len(result) == 0 {
		t.Errorf("Want non-empty branch list")
		return
	}
	if got, want := result[0].Name, "master"; got != want {
		t.Errorf("Want branch Name %q, got %q", want, got)
	}
	if got, want := result[0].Sha, "a6e5e7d797edf751cbd839d6bd4aef86c941eec9"; got != want {
		t.Errorf("Want branch Sha %q, got %q", want, got)
	}
	t.Run("Page", testPage(res))
}

func TestGitListTags(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()
	client, _ := New(server.URL)
	result, res, err := client.Git.ListTags(context.Background(), "atlassian/atlaskit", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if len(result) == 0 {
		t.Errorf("Want non-empty tag list")
		return
	}
	if got, want := result[0].Name, "@atlaskit/activity@1.0.3"; got != want {
		t.Errorf("Want tag Name %q, got %q", want, got)
	}
	if got, want := result[0].Sha, "ceb01356c3f062579bdfeb15bc53fe151b9e00f0"; got != want {
		t.Errorf("Want tag Sha %q, got %q", want, got)
	}
	t.Run("Page", testPage(res))
}

func TestGitListChanges(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Git.ListChanges(context.Background(), "atlassian/atlaskit", "425863f9dbe56d70c8dcdbf2e4e0805e85591fcc", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	if len(result) == 0 {
		t.Errorf("Want non-empty diff")
		return
	}
	if got, want := result[0].Path, "CONTRIBUTING.md"; got != want {
		t.Errorf("Want file path %q, got %q", want, got)
	}
}

func testCommit(commit *scm.Commit) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := commit.Sha, "a6e5e7d797edf751cbd839d6bd4aef86c941eec9"; got != want {
			t.Errorf("Want commit Sha %q, got %q", want, got)
		}
		if got, want := commit.Message, "Add Apache 2.0 License\n"; got != want {
			t.Errorf("Want commit Message %q, got %q", want, got)
		}
		if got, want := commit.Author.Login, "aahmed"; got != want {
			t.Errorf("Want commit author Login %q, got %q", want, got)
		}
		if got, want := commit.Author.Email, "aahmed@atlassian.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Author.Date.Unix(), int64(1440645904); got != want {
			t.Errorf("Want commit Timestamp %d, got %d", want, got)
		}
		if got, want := commit.Committer.Login, "aahmed"; got != want {
			t.Errorf("Want commit author Login %q, got %q", want, got)
		}
		if got, want := commit.Committer.Email, "aahmed@atlassian.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Committer.Date.Unix(), int64(1440645904); got != want {
			t.Errorf("Want commit Timestamp %d, got %d", want, got)
		}
	}
}

func testBranch(branch *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := branch.Name, "master"; got != want {
			t.Errorf("Want branch Name %q, got %q", want, got)
		}
		if got, want := branch.Sha, "a6e5e7d797edf751cbd839d6bd4aef86c941eec9"; got != want {
			t.Errorf("Want branch Sha %q, got %q", want, got)
		}
	}
}
