// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"reflect"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github/fixtures"
)

func TestIssueFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Issues.Find(context.Background(), "octocat/hello-world", 1)
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
	result, res, err := client.Issues.FindComment(context.Background(), "octocat/hello-world", 1, 1)
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
	result, res, err := client.Issues.List(context.Background(), "octocat/hello-world", scm.IssueListOptions{Page: 1, Size: 30, Open: true, Closed: true})
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
	result, res, err := client.Issues.ListComments(context.Background(), "octocat/hello-world", 1, scm.ListOptions{Size: 30, Page: 1})
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
	result, res, err := client.Issues.Create(context.Background(), "octocat/hello-world", &input)
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
	result, res, err := client.Issues.CreateComment(context.Background(), "octocat/hello-world", 1, input)
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
	res, err := client.Issues.DeleteComment(context.Background(), "octocat/hello-world", 1, 1)
	if err != nil {
		t.Error(err)
		return
	}
	if want, got := res.Status, 204; want != got {
		t.Errorf("Want status code %d, got %d", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueClose(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.Issues.Close(context.Background(), "octocat/hello-world", 1)
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
		if got, want := issue.Number, 1347; got != want {
			t.Errorf("Want issue Number %d, got %d", want, got)
		}
		if got, want := issue.Title, "Found a bug"; got != want {
			t.Errorf("Want issue Title %q, got %q", want, got)
		}
		if got, want := issue.Body, "I'm having a problem with this."; got != want {
			t.Errorf("Want issue Title %q, got %q", want, got)
		}
		if got, want := issue.Labels, []string{"bug"}; !reflect.DeepEqual(got, want) {
			t.Errorf("Want issue Created %v, got %v", want, got)
		}
		if got, want := issue.Closed, false; got != want {
			t.Errorf("Want issue Title %v, got %v", want, got)
		}
		if got, want := issue.Author.Login, "octocat"; got != want {
			t.Errorf("Want issue author Login %q, got %q", want, got)
		}
		if got, want := issue.Author.Avatar, "https://github.com/images/error/octocat_happy.gif"; got != want {
			t.Errorf("Want issue author Avatar %q, got %q", want, got)
		}
		if got, want := issue.Created.Unix(), int64(1303479228); got != want {
			t.Errorf("Want issue Created %d, got %d", want, got)
		}
		if got, want := issue.Updated.Unix(), int64(1303479228); got != want {
			t.Errorf("Want issue Created %d, got %d", want, got)
		}
	}
}

func testIssueComment(comment *scm.Comment) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := comment.ID, 1; got != want {
			t.Errorf("Want issue comment ID %d, got %d", want, got)
		}
		if got, want := comment.Body, "Me too"; got != want {
			t.Errorf("Want issue comment Body %q, got %q", want, got)
		}
		if got, want := comment.Author.Login, "octocat"; got != want {
			t.Errorf("Want issue comment author Login %q, got %q", want, got)
		}
		if got, want := comment.Author.Avatar, "https://github.com/images/error/octocat_happy.gif"; got != want {
			t.Errorf("Want issue comment author Avatar %q, got %q", want, got)
		}
		if got, want := comment.Created.Unix(), int64(1302796849); got != want {
			t.Errorf("Want issue comment Created %d, got %d", want, got)
		}
		if got, want := comment.Updated.Unix(), int64(1302796849); got != want {
			t.Errorf("Want issue comment Updated %d, got %d", want, got)
		}
	}
}
