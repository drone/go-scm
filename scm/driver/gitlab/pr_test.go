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

func TestPullFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.PullRequests.Find(context.Background(), "gitlab-org/testme", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testPullRequest(result))
}

func TestPullFindComment(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.PullRequests.FindComment(context.Background(), "diaspora/diaspora", 1, 301)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testPullRequestComment(result))
}

func TestPullList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.PullRequests.List(context.Background(), "gitlab-org/testme", scm.PullRequestListOptions{Page: 1, Size: 30, Open: true, Closed: false})
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

// func TestPullListChanges(t *testing.T) {
// 	server := fixtures.NewServer()
// 	defer server.Close()

// 	client, _ := New(server.URL)
// 	result, res, err := client.PullRequests.ListChanges(context.Background(), "octocat/hello-world", 1347, scm.ListOptions{Page: 1, Size: 30})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if got, want := len(result), 1; got != want {
// 		t.Errorf("Want %d pull request changes, got %d", want, got)
// 		return
// 	}
// 	t.Run("Request", testRequest(res))
// 	t.Run("Rate", testRate(res))
// 	t.Run("Page", testPage(res))
// 	t.Run("Fields", testChange(result[0]))
// }

func TestPullListComments(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.PullRequests.ListComments(context.Background(), "diaspora/diaspora", 1, scm.ListOptions{Page: 1, Size: 30})
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
	t.Run("Fields", testPullRequestComment(result[0]))
}

func TestPullMerge(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.PullRequests.Close(context.Background(), "gitlab-org/testme", 1)
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
	res, err := client.PullRequests.Close(context.Background(), "gitlab-org/testme", 1)
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

func TestPullCreateComment(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	input := &scm.CommentInput{
		Body: "Comment for MR",
	}

	client, _ := New(server.URL)
	result, res, err := client.PullRequests.CreateComment(context.Background(), "diaspora/diaspora", 1, input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testPullRequestComment(result))
}

func TestPullDeleteComment(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	res, err := client.PullRequests.DeleteComment(context.Background(), "diaspora/diaspora", 1, 1)
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

func testPullRequest(pr *scm.PullRequest) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := pr.Number, 1; got != want {
			t.Errorf("Want pr Number %d, got %d", want, got)
		}
		if got, want := pr.Title, "JS fix"; got != want {
			t.Errorf("Want pr Body %q, got %q", want, got)
		}
		if got, want := pr.Body, "Signed-off-by: Dmitriy Zaporozhets <dmitriy.zaporozhets@gmail.com>"; got != want {
			t.Errorf("Want pr Body %q, got %q", want, got)
		}
		if got, want := pr.Sha, "12d65c8dd2b2676fa3ac47d955accc085a37a9c1"; got != want {
			t.Errorf("Want pr Sha %v, got %v", want, got)
		}
		if got, want := pr.Ref, "refs/merge-requests/1/head"; got != want {
			t.Errorf("Want pr Ref %v, got %v", want, got)
		}
		if got, want := pr.Source, "fix"; got != want {
			t.Errorf("Want pr Source %v, got %v", want, got)
		}
		if got, want := pr.Target, "master"; got != want {
			t.Errorf("Want pr Target %v, got %v", want, got)
		}
		if got, want := pr.Merged, false; got != want {
			t.Errorf("Want pr Merged %v, got %v", want, got)
		}
		if got, want := pr.Closed, true; got != want {
			t.Errorf("Want pr Closed %v, got %v", want, got)
		}
		if got, want := pr.Author.Login, "dblessing"; got != want {
			t.Errorf("Want pr Author Login %v, got %v", want, got)
		}
		if got, want := pr.Author.Avatar, "https://secure.gravatar.com/avatar/b5bf44866b4eeafa2d8114bfe15da02f?s=80&d=identicon"; got != want {
			t.Errorf("Want pr Author Avatar %v, got %v", want, got)
		}
		if got, want := pr.Created.Unix(), int64(1450463393); got != want {
			t.Errorf("Want pr Created %d, got %d", want, got)
		}
		if got, want := pr.Updated.Unix(), int64(1450463422); got != want {
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
		if got, want := change.Path, "doc/update/5.4-to-6.0.md"; got != want {
			t.Errorf("Want file Path %q, got %q", want, got)
		}
	}
}

func testPullRequestComment(comment *scm.Comment) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := comment.ID, 301; got != want {
			t.Errorf("Want issue comment ID %d, got %d", want, got)
		}
		if got, want := comment.Body, "Comment for MR"; got != want {
			t.Errorf("Want issue comment Body %q, got %q", want, got)
		}
		if got, want := comment.Author.Login, "pipin"; got != want {
			t.Errorf("Want issue comment author Login %q, got %q", want, got)
		}
		if got, want := comment.Created.Unix(), int64(1380704234); got != want {
			t.Errorf("Want issue comment Created %d, got %d", want, got)
		}
		if got, want := comment.Updated.Unix(), int64(1380704234); got != want {
			t.Errorf("Want issue comment Updated %d, got %d", want, got)
		}
	}
}
