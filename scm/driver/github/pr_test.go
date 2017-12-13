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

func TestPullFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.PullRequests.Find(context.Background(), "octocat/hello-world", 1347)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testPullRequest(result))
}

// func TestPullFindComment(t *testing.T) {
// 	service := new(pullService)
// 	_, _, err := service.FindComment(context.Background(), "gogits/gogs", 1, 1)
// 	if err == nil || err.Error() != "Not Supported" {
// 		t.Errorf("Expect Not Supported error")
// 	}
// }

func TestPullList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.PullRequests.List(context.Background(), "octocat/hello-world", scm.PullRequestListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d pull requests, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Fields", testPullRequest(result[0]))
}

func TestPullListChanges(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.PullRequests.ListChanges(context.Background(), "octocat/hello-world", 1347, scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d pull request changes, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Fields", testChange(result[0]))
}

// func TestPullListComments(t *testing.T) {
// 	service := new(pullService)
// 	_, _, err := service.ListComments(context.Background(), "gogits/gogs", 1, scm.ListOptions{Page: 1, Size: 30})
// 	if err == nil || err.Error() != "Not Supported" {
// 		t.Errorf("Expect Not Supported error")
// 	}
// }

func TestPullMerge(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.PullRequests.Close(context.Background(), "octocat/hello-world", 1347)
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

func TestPullClose(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.PullRequests.Close(context.Background(), "octocat/hello-world", 1347)
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

// func TestCreateComment(t *testing.T) {
// 	service := new(pullService)
// 	_, _, err := service.CreateComment(context.Background(), "gogits/gogs", 1, &scm.CommentInput{})
// 	if err == nil || err.Error() != "Not Supported" {
// 		t.Errorf("Expect Not Supported error")
// 	}
// }

// func TestDeleteComment(t *testing.T) {
// 	service := new(pullService)
// 	_, err := service.DeleteComment(context.Background(), "gogits/gogs", 1, 1)
// 	if err == nil || err.Error() != "Not Supported" {
// 		t.Errorf("Expect Not Supported error")
// 	}
// }

func testPullRequest(pr *scm.PullRequest) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := pr.Number, 1347; got != want {
			t.Errorf("Want pr Number %d, got %d", want, got)
		}
		if got, want := pr.Title, "new-feature"; got != want {
			t.Errorf("Want pr Body %q, got %q", want, got)
		}
		if got, want := pr.Body, "Please pull these awesome changes"; got != want {
			t.Errorf("Want pr Body %q, got %q", want, got)
		}
		if got, want := pr.Sha, "6dcb09b5b57875f334f61aebed695e2e4193db5e"; got != want {
			t.Errorf("Want pr Sha %v, got %v", want, got)
		}
		if got, want := pr.Ref, "refs/pull/1347/head"; got != want {
			t.Errorf("Want pr Ref %v, got %v", want, got)
		}
		if got, want := pr.Source, "new-topic"; got != want {
			t.Errorf("Want pr Source %v, got %v", want, got)
		}
		if got, want := pr.Target, "master"; got != want {
			t.Errorf("Want pr Target %v, got %v", want, got)
		}
		if got, want := pr.Merged, true; got != want {
			t.Errorf("Want pr Merged %v, got %v", want, got)
		}
		if got, want := pr.Closed, false; got != want {
			t.Errorf("Want pr Closed %v, got %v", want, got)
		}
		if got, want := pr.Author.Login, "octocat"; got != want {
			t.Errorf("Want pr Author Login %v, got %v", want, got)
		}
		if got, want := pr.Author.Avatar, "https://github.com/images/error/octocat_happy.gif"; got != want {
			t.Errorf("Want pr Author Avatar %v, got %v", want, got)
		}
		if got, want := pr.Created.Unix(), int64(1296068472); got != want {
			t.Errorf("Want pr Created %d, got %d", want, got)
		}
		if got, want := pr.Updated.Unix(), int64(1296068472); got != want {
			t.Errorf("Want pr Updated %d, got %d", want, got)
		}
	}
}

func testChange(change *scm.Change) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := change.Added, true; got != want {
			t.Errorf("Want file Added %v, got %v", want, got)
		}
		if got, want := change.Deleted, false; got != want {
			t.Errorf("Want file Deleted %v, got %v", want, got)
		}
		if got, want := change.Renamed, false; got != want {
			t.Errorf("Want file Renamed %v, got %v", want, got)
		}
		if got, want := change.Path, "file1.txt"; got != want {
			t.Errorf("Want file Path %q, got %q", want, got)
		}
	}
}
