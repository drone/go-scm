// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/livecycle/go-scm/scm"
)

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gogs.io").
		Get("/api/v1/repos/gogits/gogs/raw/f05f642b892d59a0a9ef6a31f6c905a24b5db13a/README.md").
		Reply(200).
		Type("plain/text").
		BodyString("Hello World\n")

	client, _ := New("https://try.gogs.io")
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
	if got, want := string(result.Data), "Hello World\n"; got != want {
		t.Errorf("Want file Data %q, got %q", want, got)
	}
}

func TestContentCreate(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, err := client.Contents.Create(context.Background(), "gogits/gogs", "README.md", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentUpdate(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, err := client.Contents.Update(context.Background(), "gogits/gogs", "README.md", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentDelete(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, err := client.Contents.Delete(context.Background(), "gogits/gogs", "README.md", "master")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentList(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, _, err := client.Contents.List(context.Background(), "gogits/gogs", "/", "master", scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
