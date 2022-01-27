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

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/commits/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/commit.json")

	client := NewDefault()
	got, res, err := client.Git.FindCommit(context.Background(), "octocat/hello-world", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Commit)
	raw, _ := ioutil.ReadFile("testdata/commit.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitFindBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/branches/master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/branch.json")

	client := NewDefault()
	got, res, err := client.Git.FindBranch(context.Background(), "octocat/hello-world", "master")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/branch.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitFindTag(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/git/ref/tags/v0.1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/tag.json")

	client := NewDefault()
	got, res, err := client.Git.FindTag(context.Background(), "octocat/hello-world", "v0.1")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/tag.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitCreateBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/git/refs").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/branch_create.json")

	params := &scm.CreateBranch{
		Name: "Hello",
		Sha:  "312797ba52425353dec56871a255e2a36fc96344",
	}

	client := NewDefault()
	res, err := client.Git.CreateBranch(
		context.Background(),
		"octocat/hello-world",
		params,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 201 {
		t.Errorf("Unexpected Results")
	}
}

func TestGitListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/commits").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("sha", "master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/commits.json")

	client := NewDefault()
	got, res, err := client.Git.ListCommits(context.Background(), "octocat/hello-world", scm.CommitListOptions{Ref: "master", Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/branches").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/branches.json")

	client := NewDefault()
	got, res, err := client.Git.ListBranches(context.Background(), "octocat/hello-world", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/branches.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListTags(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/tags").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/tags.json")

	client := NewDefault()
	got, res, err := client.Git.ListTags(context.Background(), "octocat/hello-world", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/tags.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/commits/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/changes.json")

	client := NewDefault()
	got, res, err := client.Git.ListChanges(context.Background(), "octocat/hello-world", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/changes.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitCompareChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/compare/553c2077f0edc3d5dc5d17262f6aa498e69d6f8e...7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/compare.json")

	client := NewDefault()
	got, res, err := client.Git.CompareChanges(context.Background(), "octocat/hello-world", "553c2077f0edc3d5dc5d17262f6aa498e69d6f8e", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/compare.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
