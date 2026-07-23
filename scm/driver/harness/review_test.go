// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/h2non/gock"
)

func TestReviewFind(t *testing.T) {
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	_, _, err := client.Reviews.Find(context.Background(), harnessRepo, 1, 1)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewList(t *testing.T) {
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	_, _, err := client.Reviews.List(context.Background(), harnessRepo, 1, scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewDelete(t *testing.T) {
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	_, err := client.Reviews.Delete(context.Background(), harnessRepo, 1, 1)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewCreate(t *testing.T) {
	defer gock.Off()
	gock.New(gockOrigin).
		Post(fmt.Sprintf("/gateway/code/api/v1/repos/%s/pullreq/1/comments", harnessRepo)).
		MatchParam("accountIdentifier", harnessAccount).
		MatchParam("orgIdentifier", harnessOrg).
		MatchParam("projectIdentifier", harnessProject).
		MatchType("json").
		JSON(map[string]interface{}{
			"text":              "Needs a guard here",
			"path":              "file1.txt",
			"source_commit_sha": "abc123",
			"target_commit_sha": "",
			"parent_id":         float64(0),
			"line_start":        float64(5),
			"line_start_new":    true,
			"line_end":          float64(5),
			"line_end_new":      true,
		}).
		Reply(201).
		Type("application/json").
		File("testdata/review.json")

	client, _ := New(gockOrigin, harnessAccount, harnessOrg, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}

	input := &scm.ReviewInput{
		Body: "Needs a guard here",
		Sha:  "abc123",
		Path: "file1.txt",
		Line: 5,
	}

	got, _, err := client.Reviews.Create(context.Background(), harnessRepo, 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Review)
	raw, _ := ioutil.ReadFile("testdata/review.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want,
		cmpopts.IgnoreFields(scm.Review{}, "Created", "Updated"),
		cmpopts.IgnoreFields(scm.User{}, "Created", "Updated")); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
	if !gock.IsDone() {
		t.Errorf("expected request to match mock, %d pending", len(gock.Pending()))
	}
}

func TestReviewCreate_LeftSide(t *testing.T) {
	defer gock.Off()
	gock.New(gockOrigin).
		Post(fmt.Sprintf("/gateway/code/api/v1/repos/%s/pullreq/1/comments", harnessRepo)).
		MatchType("json").
		JSON(map[string]interface{}{
			"text":              "removed this?",
			"path":              "file1.txt",
			"source_commit_sha": "",
			"target_commit_sha": "",
			"parent_id":         float64(0),
			"line_start":        float64(10),
			"line_start_new":    false,
			"line_end":          float64(10),
			"line_end_new":      false,
		}).
		Reply(201).
		Type("application/json").
		File("testdata/review.json")

	client, _ := New(gockOrigin, harnessAccount, harnessOrg, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}

	input := &scm.ReviewInput{
		Body: "removed this?",
		Path: "file1.txt",
		Line: 10,
		Side: scm.SideLeft,
	}

	_, _, err := client.Reviews.Create(context.Background(), harnessRepo, 1, input)
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
	gock.New(gockOrigin).
		Post(fmt.Sprintf("/gateway/code/api/v1/repos/%s/pullreq/1/comments", harnessRepo)).
		MatchType("json").
		JSON(map[string]interface{}{
			"text":              "this whole block",
			"path":              "file1.txt",
			"source_commit_sha": "",
			"target_commit_sha": "",
			"parent_id":         float64(0),
			"line_start":        float64(15),
			"line_start_new":    true,
			"line_end":          float64(20),
			"line_end_new":      true,
		}).
		Reply(201).
		Type("application/json").
		File("testdata/review.json")

	client, _ := New(gockOrigin, harnessAccount, harnessOrg, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}

	input := &scm.ReviewInput{
		Body:      "this whole block",
		Path:      "file1.txt",
		Line:      20,
		StartLine: 15,
	}

	_, _, err := client.Reviews.Create(context.Background(), harnessRepo, 1, input)
	if err != nil {
		t.Error(err)
		return
	}
	if !gock.IsDone() {
		t.Errorf("expected request to match mock, %d pending", len(gock.Pending()))
	}
}
