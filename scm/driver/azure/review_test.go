// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/h2non/gock"
)

func TestReviewCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/pullRequests/1/threads").
		MatchType("json").
		JSON(map[string]interface{}{
			"comments": []map[string]interface{}{
				{
					"parentCommentId": float64(0),
					"content":         "looks good?",
					"commentType":     float64(1),
				},
			},
			"status": float64(1),
			"threadContext": map[string]interface{}{
				"filePath":       "/file1.txt",
				"rightFileStart": map[string]interface{}{"line": float64(5), "offset": float64(1)},
				"rightFileEnd":   map[string]interface{}{"line": float64(5), "offset": float64(1)},
			},
		}).
		Reply(200).
		Type("application/json").
		File("testdata/pr_thread.json")

	input := &scm.ReviewInput{
		Body: "looks good?",
		Line: 5,
		Path: "file1.txt",
		Sha:  "abc123",
	}

	client := NewDefault("ORG", "PROJ")
	got, res, err := client.Reviews.Create(context.Background(), "REPOID", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	if got.ID != 101 {
		t.Errorf("expected review id 101, got %d", got.ID)
	}
	if got.Body != "looks good?" {
		t.Errorf("expected body %q, got %q", "looks good?", got.Body)
	}
	if got.Path != "/file1.txt" {
		t.Errorf("expected path %q, got %q", "/file1.txt", got.Path)
	}
	if got.Line != 5 {
		t.Errorf("expected line 5, got %d", got.Line)
	}
	if res.Status != 200 {
		t.Errorf("expected status 200, got %d", res.Status)
	}
}

func TestReviewCreate_LeftSide(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/pullRequests/1/threads").
		MatchType("json").
		JSON(map[string]interface{}{
			"comments": []map[string]interface{}{
				{
					"parentCommentId": float64(0),
					"content":         "removed this?",
					"commentType":     float64(1),
				},
			},
			"status": float64(1),
			"threadContext": map[string]interface{}{
				"filePath":      "/file1.txt",
				"leftFileStart": map[string]interface{}{"line": float64(10), "offset": float64(1)},
				"leftFileEnd":   map[string]interface{}{"line": float64(10), "offset": float64(1)},
			},
		}).
		Reply(200).
		Type("application/json").
		File("testdata/pr_thread.json")

	input := &scm.ReviewInput{
		Body: "removed this?",
		Line: 10,
		Path: "file1.txt",
		Side: scm.SideLeft,
	}

	client := NewDefault("ORG", "PROJ")
	_, _, err := client.Reviews.Create(context.Background(), "REPOID", 1, input)
	if err != nil {
		t.Error(err)
		return
	}
	if !gock.IsDone() {
		t.Errorf("pending mocks: request body did not match")
	}
}

func TestReviewCreate_MultiLine(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/pullRequests/1/threads").
		MatchType("json").
		JSON(map[string]interface{}{
			"comments": []map[string]interface{}{
				{
					"parentCommentId": float64(0),
					"content":         "this whole block",
					"commentType":     float64(1),
				},
			},
			"status": float64(1),
			"threadContext": map[string]interface{}{
				"filePath":       "/file1.txt",
				"rightFileStart": map[string]interface{}{"line": float64(15), "offset": float64(1)},
				"rightFileEnd":   map[string]interface{}{"line": float64(20), "offset": float64(1)},
			},
		}).
		Reply(200).
		Type("application/json").
		File("testdata/pr_thread.json")

	input := &scm.ReviewInput{
		Body:      "this whole block",
		Line:      20,
		StartLine: 15,
		Path:      "file1.txt",
	}

	client := NewDefault("ORG", "PROJ")
	_, _, err := client.Reviews.Create(context.Background(), "REPOID", 1, input)
	if err != nil {
		t.Error(err)
		return
	}
	if !gock.IsDone() {
		t.Errorf("pending mocks: request body did not match")
	}
}
