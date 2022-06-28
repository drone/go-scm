// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

func TestReviewFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/2/reviews/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/reviews_find.json")

	client := NewDefault()
	got, res, err := client.Reviews.Find(context.Background(), "octocat/hello-world", 2, 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := os.ReadFile("testdata/reviews_find.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

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
		Get("/repos/octocat/hello-world/pulls/1/reviews").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/reviews_list.json")

	client := NewDefault()
	got, res, err := client.Reviews.List(context.Background(), "octocat/hello-world", 1, &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Review{}
	raw, _ := os.ReadFile("testdata/reviews_list.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

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
		Post("/repos/octocat/hello-world/pulls/1/reviews").
		File("testdata/reviews_create.json").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/reviews_find.json")

	input := &scm.ReviewInput{
		Body:  "This is close to perfect! Please address the suggested inline change.",
		Sha:   "ecdd80bb57125d7ba9641ffaa4d7d2c19d3f3091",
		Event: "REQUEST_CHANGES",
		Comments: []*scm.ReviewCommentInput{
			{
				Path: "file.md",
				Line: 6,
				Body: "Please add more information here, and fix this typo.",
			},
		},
	}

	client := NewDefault()
	got, res, err := client.Reviews.Create(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := os.ReadFile("testdata/reviews_find.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReviewDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Delete("/repos/octocat/hello-world/pulls/2/reviews/1").
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

func TestReviewListComments(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1/reviews/1/comments").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/reviews_list_comments.json")

	client := NewDefault()
	got, res, err := client.Reviews.ListComments(context.Background(), "octocat/hello-world", 1, 1, &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.ReviewComment{}
	raw, _ := os.ReadFile("testdata/reviews_list_comments.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestReviewUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Put("/repos/octocat/hello-world/pulls/1/reviews/1").
		File("testdata/reviews_update.json").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/reviews_find.json")

	client := NewDefault()
	got, res, err := client.Reviews.Update(context.Background(), "octocat/hello-world", 1, 1, "Updated body")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := os.ReadFile("testdata/reviews_find.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReviewSubmit(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/pulls/1/reviews/1/events").
		File("testdata/reviews_submit.json").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/reviews_find.json")

	client := NewDefault()
	input := &scm.ReviewSubmitInput{
		Body:  "",
		Event: "APPROVE",
	}
	got, res, err := client.Reviews.Submit(context.Background(), "octocat/hello-world", 1, 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := os.ReadFile("testdata/reviews_find.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestReviewDismiss(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Put("/repos/octocat/hello-world/pulls/1/reviews/1/dismissals").
		File("testdata/reviews_dismiss.json").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/reviews_find.json")

	client := NewDefault()
	got, res, err := client.Reviews.Dismiss(context.Background(), "octocat/hello-world", 1, 1, "Dismissing")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := os.ReadFile("testdata/reviews_find.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
