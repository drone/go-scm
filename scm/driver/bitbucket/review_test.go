// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

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
	_, _, err := NewDefault().Reviews.Find(context.Background(), "", 0, 0)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewList(t *testing.T) {
	_, _, err := NewDefault().Reviews.List(context.Background(), "", 0, scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewDelete(t *testing.T) {
	_, err := NewDefault().Reviews.Delete(context.Background(), "", 0, 0)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/brianharness/test/pullrequests/3/comments").
		MatchType("json").
		JSON(map[string]interface{}{
			"content": map[string]interface{}{
				"raw": "Lovely comment",
			},
			"inline": map[string]interface{}{
				"path": "README.md",
				"to":   float64(5),
			},
		}).
		Reply(201).
		Type("application/json").
		File("testdata/review.json")

	input := &scm.ReviewInput{
		Body: "Lovely comment",
		Path: "README.md",
		Line: 5,
	}

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Reviews.Create(context.Background(), "brianharness/test", 3, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := ioutil.ReadFile("testdata/review.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestReviewCreate_LeftSide(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/brianharness/test/pullrequests/3/comments").
		MatchType("json").
		JSON(map[string]interface{}{
			"content": map[string]interface{}{
				"raw": "removed this?",
			},
			"inline": map[string]interface{}{
				"path": "README.md",
				"from": float64(10),
			},
		}).
		Reply(201).
		Type("application/json").
		File("testdata/review.json")

	input := &scm.ReviewInput{
		Body: "removed this?",
		Path: "README.md",
		Line: 10,
		Side: scm.SideLeft,
	}

	client, _ := New("https://api.bitbucket.org")
	_, _, err := client.Reviews.Create(context.Background(), "brianharness/test", 3, input)
	if err != nil {
		t.Error(err)
		return
	}
	if !gock.IsDone() {
		t.Errorf("expected request to match mock, %d pending", len(gock.Pending()))
	}
}

func TestReviewCreate_MultiLine(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("/2.0/repositories/brianharness/test/pullrequests/3/comments").
		MatchType("json").
		JSON(map[string]interface{}{
			"content": map[string]interface{}{
				"raw": "this whole block",
			},
			"inline": map[string]interface{}{
				"path":     "README.md",
				"to":       float64(20),
				"start_to": float64(15),
			},
		}).
		Reply(201).
		Type("application/json").
		File("testdata/review.json")

	input := &scm.ReviewInput{
		Body:      "this whole block",
		Path:      "README.md",
		Line:      20,
		StartLine: 15,
	}

	client, _ := New("https://api.bitbucket.org")
	_, _, err := client.Reviews.Create(context.Background(), "brianharness/test", 3, input)
	if err != nil {
		t.Error(err)
		return
	}
	if !gock.IsDone() {
		t.Errorf("expected request to match mock, %d pending", len(gock.Pending()))
	}
}
