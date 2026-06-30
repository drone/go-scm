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

func TestPullFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1347").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr.json")

	client := NewDefault()
	got, res, err := client.PullRequests.Find(context.Background(), "octocat/hello-world", 1347)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestPullList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("state", "all").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pulls.json")

	client := NewDefault()
	got, res, err := client.PullRequests.List(context.Background(), "octocat/hello-world", scm.PullRequestListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.PullRequest{}
	raw, _ := ioutil.ReadFile("testdata/pulls.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestPullListChanges(t *testing.T) {
	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1347/files").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pr_files.json")

	client := NewDefault()
	got, res, err := client.PullRequests.ListChanges(context.Background(), "octocat/hello-world", 1347, scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/pr_files.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestPullFindFileDiff(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1347/files").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_files.json")

	client := NewDefault()
	got, _, err := client.PullRequests.FindFileDiff(context.Background(), "octocat/hello-world", 1347, "asda2", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	if got == nil {
		t.Fatal("Expected a change for asda2, got nil")
	}
	if got.Path != "asda2" {
		t.Errorf("Unexpected path: %s", got.Path)
	}
	if got.Patch == "" {
		t.Error("Expected non-empty patch")
	}
}

func TestPullFindFileDiff_NotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1347/files").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_files.json")

	client := NewDefault()
	got, _, err := client.PullRequests.FindFileDiff(context.Background(), "octocat/hello-world", 1347, "does/not/exist.go", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	if got != nil {
		t.Errorf("Expected nil change for a file not in the PR, got %+v", got)
	}
}

func TestPullFindFileDiff_Paginated(t *testing.T) {
	defer gock.Off()

	// page 1 does not contain the target file and advertises a next page.
	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1347/files").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pr_files.json")

	// page 2 carries the target file.
	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1347/files").
		MatchParam("page", "2").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr_files_page2.json")

	client := NewDefault()
	got, _, err := client.PullRequests.FindFileDiff(context.Background(), "octocat/hello-world", 1347, "late_file.go", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	if got == nil {
		t.Fatal("Expected to find late_file.go on a later page, got nil")
	}
	if got.Path != "late_file.go" {
		t.Errorf("Unexpected path: %s", got.Path)
	}
	if got.Patch == "" {
		t.Error("Expected non-empty patch")
	}
	if !gock.IsDone() {
		t.Errorf("expected all pages to be requested, %d pending", len(gock.Pending()))
	}
}

func TestPullMerge(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Put("/repos/octocat/hello-world/pulls/1347/merge").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.PullRequests.Merge(context.Background(), "octocat/hello-world", 1347)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestPullClose(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Patch("/repos/octocat/hello-world/pulls/1347").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.PullRequests.Close(context.Background(), "octocat/hello-world", 1347)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/pulls").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr.json")

	input := scm.PullRequestInput{
		Title:  "new-feature",
		Body:   "Please pull these awesome changes",
		Source: "new-topic",
		Target: "master",
	}

	client := NewDefault()
	got, res, err := client.PullRequests.Create(context.Background(), "octocat/hello-world", &input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("rate", testRate(res))
}

func TestPullListCommits(t *testing.T) {
	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/pulls/1347/commits").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/commits.json")

	client := NewDefault()
	got, res, err := client.PullRequests.ListCommits(context.Background(), "octocat/hello-world", 1347, scm.ListOptions{Page: 1, Size: 30})
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
