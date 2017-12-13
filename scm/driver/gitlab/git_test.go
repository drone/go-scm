// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitlab/fixtures"
)

func TestGitFindCommit(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Git.FindCommit(context.Background(), "diaspora/diaspora", "6104942438c14ec7bd21c6cd5bd995272b3faff6")
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
	result, res, err := client.Git.FindBranch(context.Background(), "diaspora/diaspora", "master")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testBranch(result))
}

func TestGitFindTag(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Git.FindTag(context.Background(), "diaspora/diaspora", "v1.0.0")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testTag(result))
}

func TestGitListCommits(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Git.ListCommits(context.Background(), "diaspora/diaspora", scm.CommitListOptions{Page: 1, Size: 30, Ref: "master"})
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
	result, res, err := client.Git.ListBranches(context.Background(), "diaspora/diaspora", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d branches, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Fields", testBranch(result[0]))
}

func TestGitListTags(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()
	client, _ := New(server.URL)
	result, res, err := client.Git.ListTags(context.Background(), "diaspora/diaspora", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Fields", testTag(result[0]))
}

func testCommit(commit *scm.Commit) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := commit.Sha, "6104942438c14ec7bd21c6cd5bd995272b3faff6"; got != want {
			t.Errorf("Want commit Sha %q, got %q", want, got)
		}
		if got, want := commit.Message, "Sanitize for network graph"; got != want {
			t.Errorf("Want commit Message %q, got %q", want, got)
		}
		if got, want := commit.Author.Login, ""; got != want {
			t.Errorf("Want commit author Login %q, got %q", want, got)
		}
		if got, want := commit.Author.Email, "dmitriy.zaporozhets@gmail.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Author.Name, "randx"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Author.Date.Unix(), int64(1340880260); got != want {
			t.Errorf("Want commit Timestamp %d, got %d", want, got)
		}
		if got, want := commit.Committer.Login, ""; got != want {
			t.Errorf("Want commit committer Login %q, got %q", want, got)
		}
		if got, want := commit.Committer.Name, "Dmitriy"; got != want {
			t.Errorf("Want commit committer Email %q, got %q", want, got)
		}
		if got, want := commit.Committer.Email, "dmitriy.zaporozhets@gmail.com"; got != want {
			t.Errorf("Want commit committer Email %q, got %q", want, got)
		}
		if got, want := commit.Committer.Date.Unix(), int64(1340880260); got != want {
			t.Errorf("Want commit Timestamp %d, got %d", want, got)
		}
	}
}

func testBranch(branch *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := branch.Name, "master"; got != want {
			t.Errorf("Want branch Name %q, got %q", want, got)
		}
		if got, want := branch.Sha, "7b5c3cc8be40ee161ae89a06bba6229da1032a0c"; got != want {
			t.Errorf("Want branch Sha %q, got %q", want, got)
		}
	}
}

func testTag(tag *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := tag.Name, "v1.0.0"; got != want {
			t.Errorf("Want tag Name %q, got %q", want, got)
		}
		if got, want := tag.Sha, "2695effb5807a22ff3d138d593fd856244e155e7"; got != want {
			t.Errorf("Want tag Sha %q, got %q", want, got)
		}
	}
}
