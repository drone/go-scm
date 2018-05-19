// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github/fixtures"
)

func TestGitFindCommit(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Git.FindCommit(context.Background(), "octocat/hello-world", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testCommit(result))
}

func TestGitFindBranch(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Git.FindBranch(context.Background(), "octocat/hello-world", "master")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testBranch(result))
}

func TestGitFindTag(t *testing.T) {
	git := new(gitService)
	_, _, err := git.FindTag(context.Background(), "octocat/hello-world", "v1.0")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestGitListCommits(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Git.ListCommits(context.Background(), "octocat/hello-world", scm.CommitListOptions{Ref: "master"})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d commits, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Fields", testCommit(result[0]))
}

func TestGitListBranches(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()
	client, _ := New(server.URL)
	result, res, err := client.Git.ListBranches(context.Background(), "octocat/hello-world", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d branches, got %d", want, got)
		return
	}
	if got, want := result[0].Name, "master"; got != want {
		t.Errorf("Want branch Name %q, got %q", want, got)
	}
	if got, want := result[0].Sha, "6dcb09b5b57875f334f61aebed695e2e4193db5e"; got != want {
		t.Errorf("Want branch Sha %q, got %q", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListTags(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()
	client, _ := New(server.URL)
	result, res, err := client.Git.ListTags(context.Background(), "octocat/hello-world", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
		return
	}
	if got, want := result[0].Name, "v0.1"; got != want {
		t.Errorf("Want tag Name %q, got %q", want, got)
	}
	if got, want := result[0].Sha, "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"; got != want {
		t.Errorf("Want tag Sha %q, got %q", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListChanges(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Git.ListChanges(context.Background(), "octocat/hello-world", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d changes, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testChange(result[0]))
}

func testCommit(commit *scm.Commit) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := commit.Sha, "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"; got != want {
			t.Errorf("Want commit Sha %q, got %q", want, got)
		}
		if got, want := commit.Message, "Merge pull request #6 from Spaceghost/patch-1\n\nNew line at end of file."; got != want {
			t.Errorf("Want commit Message %q, got %q", want, got)
		}
		if got, want := commit.Author.Login, "octocat"; got != want {
			t.Errorf("Want commit author Login %q, got %q", want, got)
		}
		if got, want := commit.Author.Email, "octocat@nowhere.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Author.Date.Unix(), int64(1331075210); got != want {
			t.Errorf("Want commit Timestamp %d, got %d", want, got)
		}
		if got, want := commit.Committer.Login, "octocat"; got != want {
			t.Errorf("Want commit author Login %q, got %q", want, got)
		}
		if got, want := commit.Committer.Email, "octocat@nowhere.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Committer.Date.Unix(), int64(1331075210); got != want {
			t.Errorf("Want commit Timestamp %d, got %d", want, got)
		}
	}
}

func testBranch(branch *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := branch.Name, "master"; got != want {
			t.Errorf("Want branch Name %q, got %q", want, got)
		}
		if got, want := branch.Sha, "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"; got != want {
			t.Errorf("Want branch Sha %q, got %q", want, got)
		}
	}
}
