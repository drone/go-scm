// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
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
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	gock.New("http://example.com:7990").
		Put("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1").
		File("testdata/pr_update.json").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	input := &scm.PullRequestInput{
		Title: "A new title",
		Body:  "A new description",
	}
	got, _, err := client.PullRequests.Update(context.Background(), "PRJ/my-repo", 1, input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, &want)

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
	json.Unmarshal(raw, &want)

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
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullDeleteComment(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comment.json")

	gock.New("http://example.com:7990").
		Delete("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments/1").
		MatchParam("version", "5").
		Reply(204)

	client, _ := New("http://example.com:7990")
	_, err := client.PullRequests.DeleteComment(context.Background(), "PRJ/my-repo", 1, 1)
	if err != nil {
		t.Error(err)
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
	json.Unmarshal(raw, &want)

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
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullMerge(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/merge").
		MatchParam("version", "0").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	_, err := client.PullRequests.Merge(context.Background(), "PRJ/my-repo", 1, nil)
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

func TestPullReopen(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/reopen").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	_, err := client.PullRequests.Reopen(context.Background(), "PRJ/my-repo", 1)
	if err != nil {
		t.Error(err)
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
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullEditComment(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comment.json")

	gock.New("http://example.com:7990").
		Put("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comment.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.EditComment(context.Background(), "PRJ/my-repo", 1, 1, &scm.CommentInput{
		Body: "LGTM",
	})
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests").
		File("testdata/create_pr.json").
		Reply(201).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")

	input := &scm.PullRequestInput{
		Title: "Updated Files",
		Body:  "* added LICENSE\n* update files\n* update files",
		Base:  "master",
		Head:  "feature/x",
	}
	got, _, err := client.PullRequests.Create(context.Background(), "PRJ/my-repo", input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullAddLabel(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments").
		JSON(pullRequestCommentInput{Text: "/jx-label test"}).
		Reply(200).
		Type("application/json").
		File("testdata/pr_label_comment.json")

	client, _ := New("http://example.com:7990")

	res, err := client.PullRequests.AddLabel(context.Background(), "PRJ/my-repo", 1, "test")

	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, res.Status, 200, "Should be a success status in response")
}

func TestPullDeleteLabel(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments").
		JSON(pullRequestCommentInput{Text: "/jx-label test remove"}).
		Reply(200).
		Type("application/json").
		File("testdata/pr_label_comment.json")

	client, _ := New("http://example.com:7990")

	res, err := client.PullRequests.DeleteLabel(context.Background(), "PRJ/my-repo", 1, "test")

	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, res.Status, 200, "Should be a success status in response")
}

func TestPullListLabels(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/activities").
		Reply(200).
		Type("application/json").
		File("testdata/pr_label_comments.json")

	client, _ := New("http://example.com:7990")

	got, _, err := client.PullRequests.ListLabels(context.Background(), "PRJ/my-repo", 1, scm.ListOptions{})

	if err != nil {
		t.Error(err)
	}

	var want []*scm.Label
	want = append(want, &scm.Label{Name: "test"})
	raw, _ := ioutil.ReadFile("testdata/pr_labels.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
