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
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"

	"github.com/jenkins-x/go-scm/scm"
)

func TestReviewFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/version").
		Reply(200).
		Type("application/json").
		File("testdata/version.json")

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/jcitizen/my-repo/pulls/1/reviews/1").
		Reply(200).
		Type("application/json").
		File("testdata/review.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Reviews.Find(context.Background(), "jcitizen/my-repo", 1, 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Review)
	raw, _ := os.ReadFile("testdata/review.json.golden")
	err = json.Unmarshal(raw, &want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestReviewList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/version").
		Reply(200).
		Type("application/json").
		File("testdata/version.json")

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/jcitizen/my-repo/pulls/1/reviews").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/reviews.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Reviews.List(context.Background(), "jcitizen/my-repo", 1, &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Review{}
	raw, _ := os.ReadFile("testdata/reviews.json.golden")
	err = json.Unmarshal(raw, &want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestReviewCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/version").
		Reply(200).
		Type("application/json").
		File("testdata/version.json")

	gock.New("https://demo.gitea.com").
		Post("/api/v1/repos/jcitizen/my-repo/pulls/1/reviews").
		File("testdata/review_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/review.json")

	client, _ := New("https://demo.gitea.com")

	in := &scm.ReviewInput{
		Body:  "This is a review",
		Sha:   "5c23b301e7eb47aa83de90cf08e0a75b4c0906c8",
		Event: "PENDING",
		Comments: []*scm.ReviewCommentInput{{
			Body: "some comment",
			Line: 10,
			Path: "some/file",
		}},
	}
	got, _, err := client.Reviews.Create(context.Background(), "jcitizen/my-repo", 1, in)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Review)
	raw, _ := os.ReadFile("testdata/review.json.golden")
	err = json.Unmarshal(raw, &want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestReviewDelete(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/version").
		Reply(200).
		Type("application/json").
		File("testdata/version.json")

	gock.New("https://demo.gitea.com").
		Delete("/api/v1/repos/jcitizen/my-repo/pulls/1/reviews/1").
		Reply(200)

	client, _ := New("https://demo.gitea.com")

	_, err := client.Reviews.Delete(context.Background(), "jcitizen/my-repo", 1, 1)
	if err != nil {
		t.Error(err)
	}
}
