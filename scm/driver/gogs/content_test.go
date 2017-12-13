// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testContents(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testContentFind(client))
		t.Run("Create", testContentCreate(client))
		t.Run("Update", testContentUpdate(client))
		t.Run("Delete", testContentDelete(client))
	}
}

func testContentFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Contents.Find(
			context.Background(),
			"gogits/gogs",
			"README.md",
			"f05f642b892d59a0a9ef6a31f6c905a24b5db13a",
		)
		if err != nil {
			t.Error(err)
		}
		if got, want := result.Path, "README.md"; got != want {
			t.Errorf("Want file Path %q, got %q", want, got)
		}
		if got, want := string(result.Data), "Hello World!\n"; got != want {
			t.Errorf("Want file Data %q, got %q", want, got)
		}
	}
}

func testContentCreate(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Contents.Create(context.Background(), "gogits/gogs", "README.md", nil)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testContentUpdate(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Contents.Update(context.Background(), "gogits/gogs", "README.md", nil)
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testContentDelete(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Contents.Delete(context.Background(), "gogits/gogs", "README.md", "master")
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}
