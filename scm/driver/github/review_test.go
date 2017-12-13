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

func TestReviewFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Reviews.Find(context.Background(), "octocat/hello-world", 1, 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testReview(result))
}

func TestReviewList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Reviews.List(context.Background(), "octocat/hello-world", 1, scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d review comments, got %d", want, got)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Fields", testReview(result[0]))
}

func TestReviewCreate(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	input := &scm.ReviewInput{
		Body: "what?",
		Line: 1,
		Path: "file1.txt",
		Sha:  "6dcb09b5b57875f334f61aebed695e2e4193db5e",
	}

	client, _ := New(server.URL)
	result, res, err := client.Reviews.Create(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}
	if want, got := res.Status, 201; want != got {
		t.Errorf("Want status code %d, got %d", want, got)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testReview(result))
}

func TestReviewDelete(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.Reviews.Delete(context.Background(), "octocat/hello-world", 1, 1)
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

func testReview(comment *scm.Review) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := comment.ID, 10; got != want {
			t.Errorf("Want issue comment ID %d, got %d", want, got)
		}
		if got, want := comment.Line, 1; got != want {
			t.Errorf("Want issue comment Line %d, got %d", want, got)
		}
		if got, want := comment.Path, "file1.txt"; got != want {
			t.Errorf("Want issue Path Body %q, got %q", want, got)
		}
		if got, want := comment.Sha, "6dcb09b5b57875f334f61aebed695e2e4193db5e"; got != want {
			t.Errorf("Want issue comment Body %q, got %q", want, got)
		}
		if got, want := comment.Body, "Great stuff"; got != want {
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
