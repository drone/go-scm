// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testIssues(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testIssueFind(client))
		t.Run("List", testIssueList(client))
		t.Run("Create", testIssueCreate(client))
		t.Run("Close", testIssueClose(client))
		t.Run("Lock", testIssueLock(client))
		t.Run("Unlock", testIssueUnlock(client))
		t.Run("Comments", testIssueComments(client))
	}
}

//
// issue sub-tests
//

func testIssueFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Issues.Find(context.Background(), "gogits/gogs", 1)
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Issue", testIssue(result))
		}
	}
}

func testIssueList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Issues.List(context.Background(), "gogits/gogs", scm.IssueListOptions{})
		if err != nil {
			t.Error(err)
		} else if got, want := len(result), 1; got != want {
			t.Errorf("Want %d issues, got %d", want, got)
		} else {
			t.Run("Issue", testIssue(result[0]))
		}
	}
}

func testIssueCreate(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		input := scm.IssueInput{
			Title: "Bug found",
			Body:  "I'm having a problem with this.",
		}
		result, _, err := client.Issues.Create(context.Background(), "gogits/gogs", &input)
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Issue", testIssue(result))
		}
	}
}

func testIssueClose(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Issues.Close(context.Background(), "gogits/go-gogs-client", 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testIssueLock(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Issues.Lock(context.Background(), "gogits/go-gogs-client", 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testIssueUnlock(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Issues.Unlock(context.Background(), "gogits/go-gogs-client", 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

//
// issue comment sub-tests
//

func testIssueComments(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testIssueCommentFind(client))
		t.Run("List", testIssueCommentList(client))
		t.Run("Create", testIssueCommentCreate(client))
		t.Run("Delete", testIssueCommentDelete(client))
	}
}

func testIssueCommentFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Issues.FindComment(context.Background(), "gogits/go-gogs-client", 1, 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testIssueCommentList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Issues.ListComments(context.Background(), "gogits/gogs", 1, scm.ListOptions{})
		if err != nil {
			t.Error(err)
		} else if got, want := len(result), 1; got != want {
			t.Errorf("Want %d comments, got %d", want, got)
		} else {
			t.Run("Comment", testIssueComment(result[0]))
		}
	}
}

func testIssueCommentCreate(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		input := &scm.CommentInput{Body: "what?"}
		result, _, err := client.Issues.CreateComment(context.Background(), "gogits/gogs", 1, input)
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Comment", testIssueComment(result))
		}
	}
}

func testIssueCommentDelete(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Issues.DeleteComment(context.Background(), "gogits/gogs", 1, 1)
		if err != nil {
			t.Error(err)
		}
	}
}

//
// struct value sub-tests
//

func testIssueComment(comment *scm.Comment) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := comment.ID, 74; got != want {
			t.Errorf("Want issue comment ID %d, got %d", want, got)
		}
		if got, want := comment.Body, "what?"; got != want {
			t.Errorf("Want issue comment Body %q, got %q", want, got)
		}
		if got, want := comment.Author.Login, "unknwon"; got != want {
			t.Errorf("Want issue comment author Login %q, got %q", want, got)
		}
		if got, want := comment.Author.Avatar, "http://localhost:3000/avatars/1"; got != want {
			t.Errorf("Want issue comment author Avatar %q, got %q", want, got)
		}
		if got, want := comment.Created.Unix(), int64(1472237898); got != want {
			t.Errorf("Want issue comment Created %d, got %d", want, got)
		}
		if got, want := comment.Updated.Unix(), int64(1472237898); got != want {
			t.Errorf("Want issue comment Updated %d, got %d", want, got)
		}
	}
}

func testIssue(issue *scm.Issue) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := issue.Number, 1; got != want {
			t.Errorf("Want issue Number %d, got %d", want, got)
		}
		if got, want := issue.Title, "Bug found"; got != want {
			t.Errorf("Want issue Title %q, got %q", want, got)
		}
		if got, want := issue.Body, "I'm having a problem with this."; got != want {
			t.Errorf("Want issue Title %q, got %q", want, got)
		}
		if got, want := issue.Closed, false; got != want {
			t.Errorf("Want issue Title %v, got %v", want, got)
		}
		if got, want := issue.Author.Login, "janedoe"; got != want {
			t.Errorf("Want issue author Login %q, got %q", want, got)
		}
		if got, want := issue.Author.Avatar, "https://secure.gravatar.com/avatar/8c58a0be77ee441bb8f8595b7f1b4e87"; got != want {
			t.Errorf("Want issue author Avatar %q, got %q", want, got)
		}
		if got, want := issue.Created.Unix(), int64(1506194641); got != want {
			t.Errorf("Want issue Created %d, got %d", want, got)
		}
		if got, want := issue.Updated.Unix(), int64(1506194641); got != want {
			t.Errorf("Want issue Created %d, got %d", want, got)
		}
	}
}
