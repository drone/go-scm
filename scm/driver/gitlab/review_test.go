// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/h2non/gock"
)

func TestReviewCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/merge_requests/1347").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/merge_with_diffrefs.json")

	gock.New("https://gitlab.com").
		Post("/api/v4/projects/diaspora/diaspora/merge_requests/1347/discussions").
		MatchType("json").
		JSON(map[string]interface{}{
			"body": "looks good to me",
			"position": map[string]interface{}{
				"position_type": "text",
				"base_sha":      "c380d3acebd181f13629a25d2e2acca46ffe1e00",
				"head_sha":      "2f63565e7aac07bcdadb654e253078b727143ec4",
				"start_sha":     "c380d3acebd181f13629a25d2e2acca46ffe1e00",
				"old_path":      "README.md",
				"new_path":      "README.md",
				"new_line":      float64(5),
			},
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/merge_discussion.json")

	input := &scm.ReviewInput{
		Body: "looks good to me",
		Path: "README.md",
		Line: 5,
		Side: scm.SideRight,
	}

	client := NewDefault()
	got, res, err := client.Reviews.Create(context.Background(), "diaspora/diaspora", 1347, input)
	if err != nil {
		t.Error(err)
		return
	}

	if got.ID != 1758 {
		t.Errorf("expected review ID 1758, got %d", got.ID)
	}
	if got.Body != "looks good to me" {
		t.Errorf("expected body %q, got %q", "looks good to me", got.Body)
	}
	if got.Path != "README.md" {
		t.Errorf("expected path %q, got %q", "README.md", got.Path)
	}
	if got.Line != 5 {
		t.Errorf("expected line 5, got %d", got.Line)
	}
	if got.Author.Login != "pipeline-bot" {
		t.Errorf("expected author %q, got %q", "pipeline-bot", got.Author.Login)
	}
	wantLink := "https://gitlab.com/diaspora/diaspora/merge_requests/1347#note_1758"
	if got.Link != wantLink {
		t.Errorf("expected link %q, got %q", wantLink, got.Link)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReviewCreate_LeftSide(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/merge_requests/1347").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/merge_with_diffrefs.json")

	gock.New("https://gitlab.com").
		Post("/api/v4/projects/diaspora/diaspora/merge_requests/1347/discussions").
		MatchType("json").
		JSON(map[string]interface{}{
			"body": "removed this?",
			"position": map[string]interface{}{
				"position_type": "text",
				"base_sha":      "c380d3acebd181f13629a25d2e2acca46ffe1e00",
				"head_sha":      "2f63565e7aac07bcdadb654e253078b727143ec4",
				"start_sha":     "c380d3acebd181f13629a25d2e2acca46ffe1e00",
				"old_path":      "README.md",
				"new_path":      "README.md",
				"old_line":      float64(10),
			},
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/merge_discussion.json")

	input := &scm.ReviewInput{
		Body: "removed this?",
		Path: "README.md",
		Line: 10,
		Side: scm.SideLeft,
	}

	client := NewDefault()
	_, res, err := client.Reviews.Create(context.Background(), "diaspora/diaspora", 1347, input)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
}

func TestReviewFind(t *testing.T) {
	client := NewDefault()
	_, _, err := client.Reviews.Find(context.Background(), "diaspora/diaspora", 1347, 1)
	if err != scm.ErrNotSupported {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

func TestReviewList(t *testing.T) {
	client := NewDefault()
	_, _, err := client.Reviews.List(context.Background(), "diaspora/diaspora", 1347, scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

func TestReviewDelete(t *testing.T) {
	client := NewDefault()
	_, err := client.Reviews.Delete(context.Background(), "diaspora/diaspora", 1347, 1)
	if err != scm.ErrNotSupported {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}
