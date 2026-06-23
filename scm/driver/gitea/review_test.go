// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/h2non/gock"
)

func TestReviewFind(t *testing.T) {
	client, _ := New("https://try.gitea.io")
	_, _, err := client.Reviews.Find(context.Background(), "go-gitea/gitea", 1, 1)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewList(t *testing.T) {
	client, _ := New("https://try.gitea.io")
	_, _, err := client.Reviews.List(context.Background(), "go-gitea/gitea", 1, scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Post("/api/v1/repos/jcitizen/my-repo/pulls/1/reviews").
		MatchType("json").
		JSON(map[string]interface{}{
			"event":     "COMMENT",
			"body":      "",
			"commit_id": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			"comments": []map[string]interface{}{
				{
					"path":         "file1.txt",
					"body":         "what?",
					"new_position": float64(1),
				},
			},
		}).
		Reply(200).
		Type("application/json").
		File("testdata/review.json")

	input := &scm.ReviewInput{
		Body: "what?",
		Line: 1,
		Path: "file1.txt",
		Sha:  "6dcb09b5b57875f334f61aebed695e2e4193db5e",
	}

	client, _ := New("https://try.gitea.io")
	got, res, err := client.Reviews.Create(context.Background(), "jcitizen/my-repo", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 200 {
		t.Errorf("expected status 200, got %d", res.Status)
	}
	if got.ID != 42 {
		t.Errorf("expected ID 42, got %d", got.ID)
	}
	if got.Body != "what?" {
		t.Errorf("expected body %q, got %q", "what?", got.Body)
	}
	if got.Path != "file1.txt" {
		t.Errorf("expected path %q, got %q", "file1.txt", got.Path)
	}
	if got.Line != 1 {
		t.Errorf("expected line 1, got %d", got.Line)
	}
	if got.Link != "https://try.gitea.io/jcitizen/my-repo/pulls/1#issuecomment-42" {
		t.Errorf("unexpected link %q", got.Link)
	}
}

func TestReviewCreate_LeftSide(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Post("/api/v1/repos/jcitizen/my-repo/pulls/1/reviews").
		MatchType("json").
		JSON(map[string]interface{}{
			"event":     "COMMENT",
			"body":      "",
			"commit_id": "abc123",
			"comments": []map[string]interface{}{
				{
					"path":         "file1.txt",
					"body":         "removed this?",
					"old_position": float64(10),
				},
			},
		}).
		Reply(200).
		Type("application/json").
		File("testdata/review.json")

	input := &scm.ReviewInput{
		Body: "removed this?",
		Line: 10,
		Path: "file1.txt",
		Sha:  "abc123",
		Side: scm.SideLeft,
	}

	client, _ := New("https://try.gitea.io")
	_, res, err := client.Reviews.Create(context.Background(), "jcitizen/my-repo", 1, input)
	if err != nil {
		t.Error(err)
		return
	}
	if res.Status != 200 {
		t.Errorf("expected status 200, got %d", res.Status)
	}
}

func TestReviewDelete(t *testing.T) {
	client, _ := New("https://try.gitea.io")
	_, err := client.Reviews.Delete(context.Background(), "go-gitea/gitea", 1, 1)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
