// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullFind(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.Find(context.Background(), "PRJ/my-repo", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullFindComment(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comment.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.FindComment(context.Background(), "PRJ/my-repo", 1, 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullList(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests").
		Reply(200).
		Type("application/json").
		File("testdata/prs.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.List(context.Background(), "PRJ/my-repo", scm.PullRequestListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.PullRequest{}
	raw, _ := ioutil.ReadFile("testdata/prs.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/changes").
		// MatchParam("pagelen", "30").
		// MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_change.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.ListChanges(context.Background(), "PRJ/my-repo", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/pr_change.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullMerge(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/merge").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	_, err := client.PullRequests.Merge(context.Background(), "PRJ/my-repo", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestPullClose(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/decline").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	_, err := client.PullRequests.Close(context.Background(), "PRJ/my-repo", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests").
		Reply(201).
		Type("application/json").
		File("testdata/pr.json")

	input := scm.PullRequestInput{
		Title:  "Updated Files",
		Body:   `* added LICENSE\r\n* update files\r\n* update files`,
		Source: "feature/x",
		Target: "master",
	}

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.Create(context.Background(), "PRJ/my-repo", &input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullCreateComment(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comment.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.CreateComment(context.Background(), "PRJ/my-repo", 1, &scm.CommentInput{
		Body: "LGTM",
	})
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullListComments(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/activities").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comments.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.ListComments(context.Background(), "PRJ/my-repo", 1, scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Comment{}
	raw, _ := ioutil.ReadFile("testdata/pr_comments.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/commits").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.ListCommits(context.Background(), "PRJ/my-repo", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
