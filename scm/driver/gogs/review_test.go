// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testReviews(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testReviewFind(client))
		t.Run("List", testReviewList(client))
		t.Run("Create", testReviewCreate(client))
		t.Run("Delete", testReviewDelete(client))
	}
}

func testReviewFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Reviews.Find(context.Background(), "gogits/gogs", 1, 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testReviewList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Reviews.List(context.Background(), "gogits/gogs", 1, scm.ListOptions{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testReviewCreate(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Reviews.Create(context.Background(), "gogits/gogs", 1, nil)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testReviewDelete(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Reviews.Delete(context.Background(), "gogits/gogs", 1, 1)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}
