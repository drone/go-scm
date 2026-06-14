// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestReviewFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/comments/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_comment.json")

	client := NewDefault()
	got, res, err := client.Reviews.Find(context.Background(), "octocat/hello-world", 2, 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReviewList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1/comments").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pr_comments.json")

	client := NewDefault()
	got, res, err := client.Reviews.List(context.Background(), "octocat/hello-world", 1, scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Review{}
	raw, _ := ioutil.ReadFile("testdata/pr_comments.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestReviewCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/pulls/1/comments").
		MatchType("json").
		JSON(map[string]interface{}{
			"body":      "what?",
			"path":      "file1.txt",
			"commit_id": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			"line":      float64(1),
			"side":      "RIGHT",
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_comment.json")

	input := &scm.ReviewInput{
		Body: "what?",
		Line: 1,
		Path: "file1.txt",
		Sha:  "6dcb09b5b57875f334f61aebed695e2e4193db5e",
	}

	client := NewDefault()
	got, res, err := client.Reviews.Create(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReviewCreate_LeftSide(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/pulls/1/comments").
		MatchType("json").
		JSON(map[string]interface{}{
			"body":      "removed this?",
			"path":      "file1.txt",
			"commit_id": "abc123",
			"line":      float64(10),
			"side":      "LEFT",
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_comment.json")

	input := &scm.ReviewInput{
		Body: "removed this?",
		Line: 10,
		Path: "file1.txt",
		Sha:  "abc123",
		Side: scm.SideLeft,
	}

	client := NewDefault()
	_, res, err := client.Reviews.Create(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
}

func TestReviewCreate_MultiLine(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/pulls/1/comments").
		MatchType("json").
		JSON(map[string]interface{}{
			"body":       "this whole block",
			"path":       "file1.txt",
			"commit_id":  "abc123",
			"line":       float64(20),
			"side":       "RIGHT",
			"start_line": float64(15),
			"start_side": "RIGHT",
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_comment.json")

	input := &scm.ReviewInput{
		Body:      "this whole block",
		Line:      20,
		StartLine: 15,
		Path:      "file1.txt",
		Sha:       "abc123",
	}

	client := NewDefault()
	_, res, err := client.Reviews.Create(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
}

func TestReviewCreate_SubjectTypeFile(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/pulls/1/comments").
		MatchType("json").
		JSON(map[string]interface{}{
			"body":         "file-level comment",
			"path":         "file1.txt",
			"commit_id":    "abc123",
			"subject_type": "file",
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_comment.json")

	input := &scm.ReviewInput{
		Body:        "file-level comment",
		Path:        "file1.txt",
		Sha:         "abc123",
		SubjectType: scm.SubjectTypeFile,
	}

	client := NewDefault()
	_, res, err := client.Reviews.Create(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
}

func TestReviewCreate_InReplyTo(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/pulls/1/comments").
		MatchType("json").
		JSON(map[string]interface{}{
			"body":        "reply to this",
			"in_reply_to": float64(42),
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_comment.json")

	input := &scm.ReviewInput{
		Body:      "reply to this",
		InReplyTo: 42,
	}

	client := NewDefault()
	_, res, err := client.Reviews.Create(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
}

func TestReviewCreate_422(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/pulls/1/comments").
		Reply(422).
		Type("application/json").
		SetHeaders(mockHeaders).
		BodyString(`{"message":"Validation Failed","errors":[{"resource":"PullRequestReviewComment","code":"custom","field":"pull_request_review_thread.line","message":"pull_request_review_thread.line must be part of the diff"}]}`)

	input := &scm.ReviewInput{
		Body: "bad line",
		Line: 999,
		Path: "file1.txt",
		Sha:  "abc123",
	}

	client := NewDefault()
	_, res, err := client.Reviews.Create(context.Background(), "octocat/hello-world", 1, input)
	if err == nil {
		t.Errorf("expected error for 422 response")
		return
	}
	if res.Status != 422 {
		t.Errorf("expected status 422, got %d", res.Status)
	}
}

func TestReviewDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Delete("/repos/octocat/hello-world/pulls/comments/1").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Reviews.Delete(context.Background(), "octocat/hello-world", 2, 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
