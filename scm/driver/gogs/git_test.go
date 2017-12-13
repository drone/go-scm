// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testGit(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Branches", testBranches(client))
		t.Run("Commits", testCommits(client))
		t.Run("Tags", testTags(client))
	}
}

//
// commit sub-tests
//

func testCommits(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testCommitFind(client))
		t.Run("List", testCommitList(client))
	}
}

func testCommitFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Git.FindCommit(
			context.Background(),
			"gogits/gogs",
			"f05f642b892d59a0a9ef6a31f6c905a24b5db13a",
		)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testCommitList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Git.ListCommits(context.Background(), "gogits/gogs", scm.CommitListOptions{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

//
// branch sub-tests
//

func testBranches(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testBranchFind(client))
		t.Run("List", testBranchList(client))
	}
}

func testBranchFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Git.FindBranch(context.Background(), "gogits/gogs", "master")
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Branch", testBranch(result))
		}
	}
}

func testBranchList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Git.ListBranches(context.Background(), "gogits/gogs", scm.ListOptions{})
		if err != nil {
			t.Error(err)
		}
		if got, want := len(result), 1; got != want {
			t.Errorf("Want %d branches, got %d", want, got)
		} else {
			t.Run("Branch", testBranch(result[0]))
		}
	}
}

//
// tag sub-tests
//

func testTags(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testTagFind(client))
		t.Run("List", testTagList(client))
	}
}

func testTagFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Git.FindTag(context.Background(), "gogits/gogs", "v1.0.0")
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testTagList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Git.ListTags(context.Background(), "gogits/gogs", scm.ListOptions{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

//
// struct value sub-tests
//

func testBranch(branch *scm.Reference) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := branch.Name, "master"; got != want {
			t.Errorf("Want branch Name %q, got %q", want, got)
		}
		if got, want := branch.Sha, "f05f642b892d59a0a9ef6a31f6c905a24b5db13a"; got != want {
			t.Errorf("Want branch Sha %q, got %q", want, got)
		}
	}
}
