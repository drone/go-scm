// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"reflect"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitlab/fixtures"
)

func TestIssueFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Issues.Find(context.Background(), "diaspora/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testIssue(result))
}

func TestIssueCommentFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Issues.FindComment(context.Background(), "diaspora/diaspora", 1, 302)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testIssueComment(result))
}

func TestIssueList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Issues.List(context.Background(), "diaspora/diaspora", scm.IssueListOptions{Page: 1, Size: 30, Open: true, Closed: false})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d issues, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Fields", testIssue(result[0]))
}

func TestIssueListComments(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Issues.ListComments(context.Background(), "diaspora/diaspora", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d comments, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Fields", testIssueComment(result[0]))
}

func TestIssueCreate(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	input := scm.IssueInput{
		Title: "Found a bug",
		Body:  "I'm having a problem with this.",
	}

	client, _ := New(server.URL)
	result, res, err := client.Issues.Create(context.Background(), "diaspora/diaspora", &input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testIssue(result))
}

func TestIssueCreateComment(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	input := &scm.CommentInput{
		Body: "what?",
	}

	client, _ := New(server.URL)
	result, res, err := client.Issues.CreateComment(context.Background(), "diaspora/diaspora", 1, input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testIssueComment(result))
}

func TestIssueCommentDelete(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.Issues.DeleteComment(context.Background(), "diaspora/diaspora", 1, 1)
	if err != nil {
		t.Error(err)
		return
	}
	if want, got := res.Status, 200; want != got {
		t.Errorf("Want status code %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueClose(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.Issues.Close(context.Background(), "diaspora/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}
	if want, got := res.Status, 200; want != got {
		t.Errorf("Want status code %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueLock(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.Issues.Lock(context.Background(), "diaspora/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}
	if want, got := res.Status, 200; want != got {
		t.Errorf("Want status code %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueUnlock(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.Issues.Unlock(context.Background(), "diaspora/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}
	if want, got := res.Status, 200; want != got {
		t.Errorf("Want status code %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func testIssue(issue *scm.Issue) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := issue.Number, 1; got != want {
			t.Errorf("Want issue Number %d, got %d", want, got)
		}
		if got, want := issue.Title, "Ut commodi ullam eos dolores perferendis nihil sunt."; got != want {
			t.Errorf("Want issue Title %q, got %q", want, got)
		}
		if got, want := issue.Body, "Omnis vero earum sunt corporis dolor et placeat."; got != want {
			t.Errorf("Want issue Body %q, got %q", want, got)
		}
		if got, want := issue.Labels, []string{}; !reflect.DeepEqual(got, want) {
			t.Errorf("Want issue Labels %v, got %v", want, got)
		}
		if got, want := issue.Closed, true; got != want {
			t.Errorf("Want issue Closed %v, got %v", want, got)
		}
		if got, want := issue.Author.Login, "root"; got != want {
			t.Errorf("Want issue author Login %q, got %q", want, got)
		}
		if got, want := issue.Author.Avatar, ""; got != want {
			t.Errorf("Want issue author Avatar %q, got %q", want, got)
		}
		if got, want := issue.Created.Unix(), int64(1451921506); got != want {
			t.Errorf("Want issue Created %d, got %d", want, got)
		}
		if got, want := issue.Updated.Unix(), int64(1451921506); got != want {
			t.Errorf("Want issue Created %d, got %d", want, got)
		}
	}
}

func testIssueComment(comment *scm.Comment) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := comment.ID, 302; got != want {
			t.Errorf("Want issue comment ID %d, got %d", want, got)
		}
		if got, want := comment.Body, "closed"; got != want {
			t.Errorf("Want issue comment Body %q, got %q", want, got)
		}
		if got, want := comment.Author.Login, "pipin"; got != want {
			t.Errorf("Want issue comment author Login %q, got %q", want, got)
		}
		if got, want := comment.Created.Unix(), int64(1380705765); got != want {
			t.Errorf("Want issue comment Created %d, got %d", want, got)
		}
		if got, want := comment.Updated.Unix(), int64(1380709365); got != want {
			t.Errorf("Want issue comment Updated %d, got %d", want, got)
		}
	}
}
