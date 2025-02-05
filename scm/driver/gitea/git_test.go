// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

//
// commit sub-tests
//

func TestCommitFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/gitea/gitea/git/commits/c43399cad8766ee521b873a32c1652407c5a4630").
		Reply(200).
		Type("application/json").
		File("testdata/commit.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Git.FindCommit(
		context.Background(),
		"gitea/gitea",
		"c43399cad8766ee521b873a32c1652407c5a4630",
	)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Commit)
	raw, _ := os.ReadFile("testdata/commit.json.golden")
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
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/commits").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Git.ListCommits(context.Background(), "go-gitea/gitea", scm.CommitListOptions{})
	if err != nil {
		t.Error(err)
	}

	var want []*scm.Commit
	raw, _ := os.ReadFile("testdata/commits.json.golden")
	err = json.Unmarshal(raw, &want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestChangeList(t *testing.T) {
	client, _ := New("https://demo.gitea.com")
	_, _, err := client.Git.ListChanges(context.Background(), "go-gitea/gitea", "f05f642b892d59a0a9ef6a31f6c905a24b5db13a", &scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestCompareCommits(t *testing.T) {
	client, _ := New("https://demo.gitea.com")
	_, _, err := client.Git.CompareCommits(context.Background(), "go-gitea/gitea", "21cf205dc770d637a9ba636644cf8bf690cc100d", "63aeb0a859499623becc1d1e7c8a2ad57439e139", &scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

//
// branch sub-tests
//

func TestBranchFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/branches/master").
		Reply(200).
		Type("application/json").
		File("testdata/branch.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Git.FindBranch(context.Background(), "go-gitea/gitea", "master")
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

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/branches").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/branches.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Git.ListBranches(context.Background(), "go-gitea/gitea", &scm.ListOptions{})
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

	t.Run("Page", testPage(res))
}

//
// tag sub-tests
//

func TestTagFind(t *testing.T) {
	client, _ := New("https://demo.gitea.com")
	_, _, err := client.Git.FindTag(context.Background(), "go-gitea/gitea", "v1.0.0")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestTagList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/tags").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/tags.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Git.ListTags(context.Background(), "go-gitea/gitea", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Reference{}
	raw, _ := os.ReadFile("testdata/tags.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestGitGetDefaultBranch(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/branches/master").
		Reply(200).
		Type("application/json").
		File("testdata/branch.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Git.GetDefaultBranch(context.Background(), "go-gitea/gitea")
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
