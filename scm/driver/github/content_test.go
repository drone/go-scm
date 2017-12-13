// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github/fixtures"
)

func TestContentFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Contents.Find(context.Background(), "octocat/hello-world", "README", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := result.Path, "README"; got != want {
		t.Errorf("Want content Path %q, got %q", want, got)
	}
	if got, want := string(result.Data), "Hello World!\n"; got != want {
		t.Errorf("Want content Body %q, got %q", want, got)
	}
}

func TestContentCreate(t *testing.T) {
	content := new(contentService)
	_, err := content.Create(context.Background(), "octocat/hello-world", "README", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentUpdate(t *testing.T) {
	content := new(contentService)
	_, err := content.Update(context.Background(), "octocat/hello-world", "README", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentDelete(t *testing.T) {
	content := new(contentService)
	_, err := content.Delete(context.Background(), "octocat/hello-world", "README", "master")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
