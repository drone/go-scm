// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

//
// commit sub-tests
//

func TestCommitFind(t *testing.T) {
	gock.New("https://try.gogs.io").
		Get("/api/v1/repos/gogs/gogs/commits/2c3e2b701e012294d457937e6bfbffd63dd8ae4f").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")
	client, _ := New("https://try.gogs.io")
	got, _, err := client.Git.FindCommit(
		context.Background(),
		"gogs/gogs",
		"2c3e2b701e012294d457937e6bfbffd63dd8ae4f",
	)
	if err != nil {
		t.Error(err)
	}
	want := new(scm.Commit)
	raw, _ := os.ReadFile("testdata/commits.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestCommitList(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, _, err := client.Git.ListCommits(context.Background(), "gogits/gogs", scm.CommitListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestChangeList(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, _, err := client.Git.ListChanges(context.Background(), "gogits/gogs", "f05f642b892d59a0a9ef6a31f6c905a24b5db13a", &scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestCompareCommits(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, _, err := client.Git.CompareCommits(context.Background(), "gogits/gogs", "f05f642b892d59a0a9ef6a31f6c905a24b5db13a", "atotallydifferentshaofexactlysamelengths", &scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

//
// branch sub-tests
//

func TestBranchFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gogs.io").
		Get("/api/v1/repos/gogits/gogs/branches/master").
		Reply(200).
		Type("application/json").
		File("testdata/branch.json")

	client, _ := New("https://try.gogs.io")
	got, _, err := client.Git.FindBranch(context.Background(), "gogits/gogs", "master")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Reference)
	raw, _ := os.ReadFile("testdata/branch.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestBranchList(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gogs.io").
		Get("/api/v1/repos/gogits/gogs/branches").
		Reply(200).
		Type("application/json").
		File("testdata/branches.json")

	client, _ := New("https://try.gogs.io")
	got, _, err := client.Git.ListBranches(context.Background(), "gogits/gogs", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Reference{}
	raw, _ := os.ReadFile("testdata/branches.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

//
// tag sub-tests
//

func TestTagFind(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, _, err := client.Git.FindTag(context.Background(), "gogits/gogs", "v1.0.0")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestTagList(t *testing.T) {
	client, _ := New("https://try.gogs.io")
	_, _, err := client.Git.ListTags(context.Background(), "gogits/gogs", &scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
