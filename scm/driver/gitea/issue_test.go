// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

//
// issue sub-tests
//

func TestIssueFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/issues/1").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Issues.Find(context.Background(), "go-gitea/gitea", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Issue)
	raw, _ := os.ReadFile("testdata/issue.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestIssueList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/issues").
		MatchParam("type", "issues").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/issues.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Issues.List(context.Background(), "go-gitea/gitea", scm.IssueListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Issue{}
	raw, _ := os.ReadFile("testdata/issues.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestIssueCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Post("/api/v1/repos/go-gitea/gitea/issues").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	input := scm.IssueInput{
		Title: "Bug found",
		Body:  "I'm having a problem with this.",
	}

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Issues.Create(context.Background(), "go-gitea/gitea", &input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Issue)
	raw, _ := os.ReadFile("testdata/issue.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestIssueClose(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Patch("/api/v1/repos/go-gitea/gitea/issues/1").
		File("testdata/close_issue.json").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.Close(context.Background(), "go-gitea/gitea", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestIssueReopen(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Patch("/api/v1/repos/go-gitea/gitea/issues/1").
		File("testdata/reopen_issue.json").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.Reopen(context.Background(), "go-gitea/gitea", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestIssueLock(t *testing.T) {
	defer gock.Off()

	mockServerVersion()
	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.Lock(context.Background(), "gogits/go-gogs-client", 1)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestIssueUnlock(t *testing.T) {
	defer gock.Off()

	mockServerVersion()
	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.Unlock(context.Background(), "gogits/go-gogs-client", 1)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

//
// issue comment sub-tests
//

func TestIssueCommentFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/issues/1/comments").
		Reply(200).
		Type("application/json").
		File("testdata/comments.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Issues.FindComment(context.Background(), "go-gitea/gitea", 1, 74)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := os.ReadFile("testdata/comment.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestIssueCommentList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/issues/1/comments").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/comments.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Issues.ListComments(context.Background(), "go-gitea/gitea", 1, &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Comment{}
	raw, _ := os.ReadFile("testdata/comments.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestIssueCommentCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Post("/api/v1/repos/go-gitea/gitea/issues/1/comments").
		Reply(201).
		Type("application/json").
		File("testdata/comment.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Issues.CreateComment(context.Background(), "go-gitea/gitea", 1, &scm.CommentInput{Body: "what?"})
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := os.ReadFile("testdata/comment.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	if gock.IsPending() {
		t.Errorf("Pending API calls")
	}
}

func TestIssueCommentDelete(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Delete("/api/v1/repos/go-gitea/gitea/issues/comments/1").
		Reply(204).
		Type("application/json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.DeleteComment(context.Background(), "go-gitea/gitea", 1, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestIssueListLabels(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/issues/1/labels").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/issue_labels.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Issues.ListLabels(context.Background(), "go-gitea/gitea", 1, &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Label{}
	raw, _ := os.ReadFile("testdata/issue_labels.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestIssueAssignIssue(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/issues/1").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	gock.New("https://demo.gitea.com").
		Patch("/api/v1/repos/go-gitea/gitea/issues/1").
		File("testdata/assign_issue.json").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.AssignIssue(context.Background(), "go-gitea/gitea", 1, []string{"a", "b"})
	if err != nil {
		t.Error(err)
	}
}

func TestIssueUnassignIssue(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/issues/1").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	gock.New("https://demo.gitea.com").
		Patch("/api/v1/repos/go-gitea/gitea/issues/1").
		File("testdata/unassign_issue.json").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.UnassignIssue(context.Background(), "go-gitea/gitea", 1, []string{"string"})
	if err != nil {
		t.Error(err)
	}
}

func TestIssueSetMilestone(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Patch("/api/v1/repos/go-gitea/gitea/issues/1").
		File("testdata/issue_set_milestone.json").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.SetMilestone(context.Background(), "go-gitea/gitea", 1, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestIssueClearMilestone(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Patch("/api/v1/repos/go-gitea/gitea/issues/1").
		File("testdata/issue_clear_milestone.json").
		Reply(200).
		Type("application/json").
		File("testdata/issue.json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Issues.ClearMilestone(context.Background(), "go-gitea/gitea", 1)
	if err != nil {
		t.Error(err)
	}
}
