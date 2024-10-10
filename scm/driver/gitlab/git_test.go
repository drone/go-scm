// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

// TODO(bradrydzewski) missing commit link
// TODO(bradrydzewski) missing commit author avatar
// TODO(bradrydzewski) missing commit committer avatar

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/commits/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/commit.json")

	client := NewDefault()
	got, res, err := client.Git.FindCommit(context.Background(), "diaspora/diaspora", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	if err != nil {
		t.Error(err)
		return
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

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitFindBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/branches/master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/branch.json")

	client := NewDefault()
	got, res, err := client.Git.FindBranch(context.Background(), "diaspora/diaspora", "master")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := os.ReadFile("testdata/branch.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitFindTag(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/tags/v1.0.0").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/tag.json")

	client := NewDefault()
	got, res, err := client.Git.FindTag(context.Background(), "diaspora/diaspora", "v1.0.0")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := os.ReadFile("testdata/tag.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("api/v4/projects/diaspora/diaspora/repository/commits").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("ref_name", "master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/commits.json")

	client := NewDefault()
	got, res, err := client.Git.ListCommits(context.Background(), "diaspora/diaspora", scm.CommitListOptions{Ref: "master", Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Commit{}
	raw, _ := os.ReadFile("testdata/commits.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/branches").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/branches.json")

	client := NewDefault()
	got, res, err := client.Git.ListBranches(context.Background(), "diaspora/diaspora", &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
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

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListTags(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/tags").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/tags.json")

	client := NewDefault()
	got, res, err := client.Git.ListTags(context.Background(), "diaspora/diaspora", &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
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

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/commits/6104942438c14ec7bd21c6cd5bd995272b3faff6/diff").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/commit_diff.json")

	client := NewDefault()
	got, res, err := client.Git.ListChanges(context.Background(), "diaspora/diaspora", "6104942438c14ec7bd21c6cd5bd995272b3faff6", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := os.ReadFile("testdata/commit_diff.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitCompareCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/gitterHQ/webapp/repository/compare").
		MatchParam("from", "6da006adb7cafe15b8495e3b7811fc318e485553").
		MatchParam("to", "c5895070235cadd2d839136dad79e01838ee2de1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/compare.json")

	client := NewDefault()
	got, res, err := client.Git.CompareCommits(context.Background(), "gitterHQ/webapp", "6da006adb7cafe15b8495e3b7811fc318e485553", "c5895070235cadd2d839136dad79e01838ee2de1", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := os.ReadFile("testdata/compare.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitCreateRef(t *testing.T) {
	baseSHA := "aa218f56b14c9653891f9e74264a383fa43fefbd"
	defer gock.Off()
	gock.New("https://gitlab.com").
		Post("/api/v4/projects/diaspora/diaspora/repository/branches").
		MatchParam("branch", "testing").
		MatchParam("ref", baseSHA).
		Reply(http.StatusCreated).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/create_branch.json")

	client := NewDefault()
	got, res, err := client.Git.CreateRef(context.Background(), "diaspora/diaspora", "testing", baseSHA)
	if err != nil {
		t.Fatal(err)
	}

	want := &scm.Reference{}
	raw, _ := os.ReadFile("testdata/create_branch.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
