// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/raw/README").
		MatchParam("at", "5c64a07cd6c0f21b753bf261ef059c7e7633c50a").
		Reply(200).
		Type("text/plain").
		File("testdata/content.txt")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Contents.Find(context.Background(), "PRJ/my-repo", "README", "5c64a07cd6c0f21b753bf261ef059c7e7633c50a")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Content)
	raw, _ := os.ReadFile("testdata/content.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/files/pkg").
		MatchParam("at", "5c64a07cd6c0f21b753bf261ef059c7e7633c50a").
		Reply(200).
		Type("text/plain").
		File("testdata/content_list.json")

	client, _ := New("http://example.com:7990")
	_, resp, err := client.Contents.List(context.Background(), "PRJ/my-repo", "pkg", "5c64a07cd6c0f21b753bf261ef059c7e7633c50a", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	if resp.Status != 200 {
		t.Errorf("got %d", resp.Status)
	}
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Put("/rest/api/1.0/projects/octocat/repos/hello-world/browse/README").
		Reply(200).
		Type("application/json").
		File("testdata/content_create.json")

	params := &scm.ContentParams{
		Message: "my commit message",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
		Signature: scm.Signature{
			Name:  "Zaki",
			Email: "zaki@example.com",
		},
	}

	client, _ := New("http://example.com:7990")
	resp, err := client.Contents.Create(context.Background(), "octocat/hello-world", "README", params)
	if err != nil {
		t.Error(err)
	}

	if resp.Status != 200 {
		t.Errorf("got %d", resp.Status)
	}
}

func TestContentUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Put("/rest/api/1.0/projects/octocat/repos/hello-world/browse/README").
		Reply(200).
		Type("application/json").
		File("testdata/content_update.json")

	params := &scm.ContentParams{
		Message: "my commit message",
		Data:    []byte("bXkgbmV3IGZpbGUgY29udGVudHM="),
		Sha:     "95b966ae1c166bd92f8ae7d1c313e738c731dfc3",
		Signature: scm.Signature{
			Name:  "Zaki",
			Email: "zaki@example.com",
		},
	}

	client, _ := New("http://example.com:7990")
	resp, err := client.Contents.Update(context.Background(), "octocat/hello-world", "README", params)
	if err != nil {
		t.Error(err)
	}

	if resp.Status != 200 {
		t.Errorf("got %d", resp.Status)
	}
}

func TestContentDelete(t *testing.T) {
	content := new(contentService)
	_, err := content.Delete(context.Background(), "atlassian/atlaskit", "README", &scm.ContentParams{Ref: "master"})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
