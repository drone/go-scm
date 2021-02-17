// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/livecycle/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/src/425863f9dbe56d70c8dcdbf2e4e0805e85591fcc/README").
		Reply(200).
		Type("text/plain").
		File("testdata/content.txt")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Contents.Find(context.Background(), "atlassian/atlaskit", "README", "425863f9dbe56d70c8dcdbf2e4e0805e85591fcc")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestContentCreate(t *testing.T) {
	content := new(contentService)
	_, err := content.Create(context.Background(), "atlassian/atlaskit", "README", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentUpdate(t *testing.T) {
	content := new(contentService)
	_, err := content.Update(context.Background(), "atlassian/atlaskit", "README", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentDelete(t *testing.T) {
	content := new(contentService)
	_, err := content.Delete(context.Background(), "atlassian/atlaskit", "README", "master")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/src/master/packages/activity").
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Contents.List(context.Background(), "atlassian/atlaskit", "packages/activity", "master", scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.ContentInfo{}
	raw, _ := ioutil.ReadFile("testdata/content_list.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
