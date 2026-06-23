// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/h2non/gock"
)

func TestReviewCreate(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments").
		Reply(201).
		Type("application/json").
		File("testdata/pr_comment.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Reviews.Create(context.Background(), "PRJ/my-repo", 1, &scm.ReviewInput{
		Body: "this is a comment",
		Path: "path/to/file.go",
		Line: 42,
		Side: scm.SideRight,
	})
	if err != nil {
		t.Error(err)
		return
	}

	if got.ID != 1 {
		t.Errorf("Want review ID 1, got %d", got.ID)
	}
	if got.Body != "this is a comment" {
		t.Errorf("Want body %q, got %q", "this is a comment", got.Body)
	}
	if got.Path != "path/to/file.go" {
		t.Errorf("Want path %q, got %q", "path/to/file.go", got.Path)
	}
	if got.Line != 42 {
		t.Errorf("Want line 42, got %d", got.Line)
	}
}

func TestReviewCreate_LeftSide(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments").
		JSON(map[string]interface{}{
			"text": "old side comment",
			"anchor": map[string]interface{}{
				"diffType": "EFFECTIVE",
				"line":     7,
				"lineType": "REMOVED",
				"fileType": "FROM",
				"path":     "path/to/file.go",
			},
		}).
		Reply(201).
		Type("application/json").
		File("testdata/pr_comment.json")

	client, _ := New("http://example.com:7990")
	_, _, err := client.Reviews.Create(context.Background(), "PRJ/my-repo", 1, &scm.ReviewInput{
		Body: "old side comment",
		Path: "path/to/file.go",
		Line: 7,
		Side: scm.SideLeft,
	})
	if err != nil {
		t.Error(err)
	}
}

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
