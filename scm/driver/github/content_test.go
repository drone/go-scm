// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/contents/README").
		MatchParam("ref", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	client := NewDefault()
	got, res, err := client.Contents.Find(
		context.Background(),
		"octocat/hello-world",
		"README",
		"7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
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
