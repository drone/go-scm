// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"context"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
)

//
// git sub-tests
//

func testGit(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Branches", testBranches(client))
		t.Run("Commits", testCommits(client))
		t.Run("Tags", testTags(client))
	}
}

//
// branch sub-tests
//

func testBranches(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testBranchFind(client))
		t.Run("List", testBranchList(client))
	}
}

func testBranchFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Git.FindBranch(context.Background(), "octocat/Hello-World", "master")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("Branch", testBranch(result))
	}
}

func testBranchList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.ListOptions{}
		result, _, err := client.Git.ListBranches(context.Background(), "octocat/Hello-World", opts)
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty branch list")
		}
		for _, branch := range result {
			if branch.Name == "master" {
				t.Run("Branch", testBranch(branch))
			}
		}
	}
}

//
// branch sub-tests
//

func testTags(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testTagFind(client))
		t.Run("List", testTagList(client))
	}
}

func testTagFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Skipf("Not Supported")
	}
}

func testTagList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.ListOptions{}
		result, _, err := client.Git.ListTags(context.Background(), "octocat/linguist", opts)
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty tag list")
		}
		for _, tag := range result {
			if tag.Name == "v4.8.8" {
				t.Run("Tag", testTag(tag))
			}
		}
	}
}

//
// commit sub-tests
//

func testCommits(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testCommitFind(client))
		t.Run("List", testCommitList(client))
		t.Run("BranchList", testBranchCommitList(client))
	}
}

func testCommitFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Git.FindCommit(context.Background(), "octocat/Hello-World", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("Commit", testCommit(result))
	}
}

func testCommitList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.CommitListOptions{
			Ref: "master",
		}
		result, _, err := client.Git.ListCommits(context.Background(), "octocat/Hello-World", opts)
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty commit list")
		}
		for _, commit := range result {
			if commit.Sha == "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d" {
				t.Run("Commit", testCommit(commit))
			}
		}
	}
}

func testBranchCommitList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		opts := scm.CommitListOptions{
			Sha: "test",
		}
		result, _, err := client.Git.ListCommits(context.Background(), "octocat/Hello-World", opts)
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty commit list")
		}

		if result[0].Sha != "b3cbd5bbd7e81436d2eee04537ea2b4c0cad4cdf" {
			t.Errorf("Unexpected commit")
		}
	}
}

//
// change sub-tests
//

func testChangeList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Git.ListChanges(context.Background(), "ossu/computer-science", "a92b5077b4b0796b680d2a41472c594351ad3f35", scm.ListOptions{})
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty change list")
		}
		for _, change := range result {
			t.Run("Change", testCommitChange(change))
		}
	}
}

func testCompareCommits(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Git.CompareCommits(context.Background(), "ossu/computer-science", "f3e6b8608c05b6c2c21384de2c5dcca43f336ed0", "a92b5077b4b0796b680d2a41472c594351ad3f35", scm.ListOptions{})
		if err != nil {
			t.Error(err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Want a non-empty change list")
		}
		for _, change := range result {
			t.Run("Change", testCommitChange(change))
		}
	}
}

//
// struct sub-tests
//

func testBranch(branch *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := branch.Name, "master"; got != want {
			t.Errorf("Want branch Name %q, got %q", want, got)
		}
		if got, want := branch.Sha, "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"; got != want {
			t.Errorf("Want branch Avatar %q, got %q", want, got)
		}
	}
}

func testTag(tag *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := tag.Name, "v4.8.8"; got != want {
			t.Errorf("Want tag Name %q, got %q", want, got)
		}
		if got, want := tag.Sha, "3f4b8368e81430e3353cb5ad8b781cd044697347"; got != want {
			t.Errorf("Want tag Avatar %q, got %q", want, got)
		}
	}
}

func testCommit(commit *scm.Commit) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := commit.Message, "Merge pull request #6 from Spaceghost/patch-1\n\nNew line at end of file."; got != want {
			t.Errorf("Want commit Message %q, got %q", want, got)
		}
		if got, want := commit.Sha, "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"; got != want {
			t.Errorf("Want commit Sha %q, got %q", want, got)
		}
		if got, want := commit.Author.Name, "The Octocat"; got != want {
			t.Errorf("Want commit author Name %q, got %q", want, got)
		}
		if got, want := commit.Author.Email, "octocat@nowhere.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Author.Date.Unix(), int64(1331075210); got != want {
			t.Errorf("Want commit author Date %d, got %d", want, got)
		}
		if got, want := commit.Committer.Name, "The Octocat"; got != want {
			t.Errorf("Want commit author Name %q, got %q", want, got)
		}
		if got, want := commit.Committer.Email, "octocat@nowhere.com"; got != want {
			t.Errorf("Want commit author Email %q, got %q", want, got)
		}
		if got, want := commit.Committer.Date.Unix(), int64(1331075210); got != want {
			t.Errorf("Want commit author Date %d, got %d", want, got)
		}
		if got, want := commit.Link, "https://github.com/octocat/Hello-World/commit/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"; got != want {
			t.Errorf("Want commit link %q, got %q", want, got)
		}
	}
}

func testCommitChange(change *scm.Change) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := change.Path, "README.md"; got != want {
			t.Errorf("Want commit Path %q, got %q", want, got)
		}
		if got, want := change.PreviousPath, ""; got != want {
			t.Errorf("Want commit PreviousPath %q, got %q", want, got)
		}
		if got, want := change.Added, false; got != want {
			t.Errorf("Want commit Added %t, got %t", want, got)
		}
		if got, want := change.Renamed, false; got != want {
			t.Errorf("Want commit Renamed %t, got %t", want, got)
		}
		if got, want := change.Deleted, false; got != want {
			t.Errorf("Want commit Deleted %t, got %t", want, got)
		}
		if got, want := change.Additions, 1; got != want {
			t.Errorf("Want commit Additions %d, got %d", want, got)
		}
		if got, want := change.Deletions, 1; got != want {
			t.Errorf("Want commit Deletions %d, got %d", want, got)
		}
		if got, want := change.Changes, 2; got != want {
			t.Errorf("Want commit Changes %d, got %d", want, got)
		}
		if got, want := change.BlobURL, "https://github.com/ossu/computer-science/blob/a92b5077b4b0796b680d2a41472c594351ad3f35/README.md"; got != want {
			t.Errorf("Want commit BlobURL %q, got %q", want, got)
		}
		if got, want := change.Sha, "487f3788c10aa8367e2a299dbdc06da48e709baa"; got != want {
			t.Errorf("Want commit Sha %q, got %q", want, got)
		}
	}
}
