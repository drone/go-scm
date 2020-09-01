// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
	"github.com/jenkins-x/go-scm/scm"
)

//
// pull request sub-tests
//

func TestPullRequestFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Get("/api/v1/repos/jcitizen/my-repo/pulls/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("https://try.gitea.io")
	got, _, err := client.PullRequests.Find(context.Background(), "jcitizen/my-repo", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullRequestList(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Get("/api/v1/repos/jcitizen/my-repo/pulls").
		Reply(200).
		Type("application/json").
		File("testdata/prs.json")

	client, _ := New("https://try.gitea.io")
	got, _, err := client.PullRequests.List(context.Background(), "jcitizen/my-repo", scm.PullRequestListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.PullRequest{}
	raw, _ := ioutil.ReadFile("testdata/prs.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullClose(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Patch("/api/v1/repos/go-gitea/gitea/pulls/1").
		File("testdata/close_pr.json").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("https://try.gitea.io")
	_, err := client.PullRequests.Close(context.Background(), "go-gitea/gitea", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestPullReopen(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Patch("/api/v1/repos/go-gitea/gitea/pulls/1").
		File("testdata/reopen_pr.json").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("https://try.gitea.io")
	_, err := client.PullRequests.Reopen(context.Background(), "go-gitea/gitea", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestPullRequestMerge(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Post("/api/v1/repos/go-gitea/gitea/pulls/1").
		Reply(204).
		Type("application/json")

	client, _ := New("https://try.gitea.io")
	_, err := client.PullRequests.Merge(context.Background(), "go-gitea/gitea", 1, nil)
	if err != nil {
		t.Error(err)
	}
}

//
// pull request change sub-tests
//

// TODO: Actually write a real test for generating changes from a patch
func TestPullRequestChanges(t *testing.T) {
	client, _ := New("https://try.gitea.io")
	_, _, err := client.PullRequests.ListChanges(context.Background(), "go-gitea/gitea", 1, scm.ListOptions{})
	if err == scm.ErrNotSupported {
		t.Errorf("Didn't expect Not Supported error")
	}
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Post("/api/v1/repos/jcitizen/my-repo/pulls").
		File("testdata/pr_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	input := &scm.PullRequestInput{
		Title: "Add License File",
		Body:  "Using a BSD License",
		Head:  "feature",
		Base:  "master",
	}

	client, _ := New("https://try.gitea.io")
	got, _, err := client.PullRequests.Create(context.Background(), "jcitizen/my-repo", input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
