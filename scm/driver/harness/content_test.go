// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

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

	gock.New(gockOrigin).
		Get("/api/v1/repos/1/content/README.md").
		Reply(200).
		Type("plain/text").
		File("testdata/content.json")

	client, _ := New(gockOrigin)
	result, _, err := client.Contents.Find(
		context.Background(),
		"1",
		"README.md",
		"",
	)
	if err != nil {
		t.Error(err)
	}

	if got, want := result.Path, "README.md"; got != want {
		t.Errorf("Want file Path %q, got %q", want, got)
	}
	if got, want := string(result.Data), "# string\nstrinasdasdsag"; got != want {
		t.Errorf("Want file Data %q, got %q", want, got)
	}
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New(gockOrigin).
		Get("/api/v1/repos/1/").
		Reply(200).
		Type("application/json").
		File("testdata/content_list.json")

	client, _ := New(gockOrigin)
	got, _, err := client.Contents.List(
		context.Background(),
		"1",
		"",
		"",
		scm.ListOptions{},
	)
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
