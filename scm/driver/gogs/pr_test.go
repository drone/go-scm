// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testPullRequests(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testPullRequestFind(client))
		t.Run("List", testPullRequestList(client))
		t.Run("Close", testPullRequestClose(client))
		t.Run("Merge", testPullRequestMerge(client))
		t.Run("Changes", testPullRequestChanges(client))
		t.Run("Comments", testPullRequestComments(client))
	}
}

//
// pull request sub-tests
//

func testPullRequestFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.PullRequests.Find(context.Background(), "gogits/gogs", 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testPullRequestList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.PullRequests.List(context.Background(), "gogits/gogs", scm.PullRequestListOptions{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testPullRequestClose(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.PullRequests.Close(context.Background(), "gogits/gogs", 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testPullRequestMerge(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.PullRequests.Merge(context.Background(), "gogits/gogs", 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

//
// pull request change sub-tests
//

func testPullRequestChanges(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.PullRequests.ListChanges(context.Background(), "gogits/gogs", 1, scm.ListOptions{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

//
// pull request comment sub-tests
//

func testPullRequestComments(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testPullRequestCommentFind(client))
		t.Run("List", testPullRequestCommentList(client))
		t.Run("Create", testPullRequestCommentCreate(client))
		t.Run("Delete", testPullRequestCommentDelete(client))
	}
}

func testPullRequestCommentFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.PullRequests.FindComment(context.Background(), "gogits/gogs", 1, 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testPullRequestCommentList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.PullRequests.ListComments(context.Background(), "gogits/gogs", 1, scm.ListOptions{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testPullRequestCommentCreate(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.PullRequests.CreateComment(context.Background(), "gogits/gogs", 1, &scm.CommentInput{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testPullRequestCommentDelete(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.PullRequests.DeleteComment(context.Background(), "gogits/gogs", 1, 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}
