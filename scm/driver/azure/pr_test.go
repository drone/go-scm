// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(201).
		Type("application/json").
		File("testdata/pr_active.json")

	input := scm.PullRequestInput{
		Title: "test_pr",
		Body:  "test_pr_body",
		Head:  "pr_branch",
		Base:  "main",
	}

	client := NewDefault()
	got, _, err := client.PullRequests.Create(context.Background(), "ORG/PROJ/REPOID", &input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := os.ReadFile("testdata/pr_active.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/pullrequests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_active.json")

	client := NewDefault()
	got, _, err := client.PullRequests.Find(context.Background(), "ORG/PROJ/REPOID", 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := os.ReadFile("testdata/pr_active.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestMerge(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/pullrequests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_active.json")

	gock.New("https://dev.azure.com/").
		Patch("/ORG/PROJ/_apis/git/repositories/REPOID/pullrequests/1").
		File("testdata/pr_merge.json").
		Reply(200).
		Type("application/json").
		File("testdata/pr_completed.json")

	client := NewDefault()
	_, err := client.PullRequests.Merge(context.Background(), "ORG/PROJ/REPOID", 1, &scm.PullRequestMergeOptions{})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestFailedMerge(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/pullrequests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_active.json")

	// Special Azure Devops behavior. Sometimes (when you've just created a PR)
	// the patch is accepted (200) but the result is still an `active` pr.
	gock.New("https://dev.azure.com/").
		Patch("/ORG/PROJ/_apis/git/repositories/REPOID/pullrequests/1").
		File("testdata/pr_merge.json").
		Reply(200).
		Type("application/json").
		File("testdata/pr_active.json")

	client := NewDefault()
	_, err := client.PullRequests.Merge(context.Background(), "ORG/PROJ/REPOID", 1, &scm.PullRequestMergeOptions{})

	if err.Error() != "patch accepted, but status still active" {
		if err != nil {
			t.Errorf("expected custom error text, wanted 'patch accepted, but status still active' but got %s", err.Error())
		} else {
			t.Errorf("expected custom error text, wanted 'patch accepted, but status still active' but got no error")
		}
	}
}

func TestClose(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Patch("/ORG/PROJ/_apis/git/repositories/REPOID/pullrequests/1").
		File("testdata/pr_abandon.json").
		Reply(200)

	client := NewDefault()
	_, err := client.PullRequests.Close(context.Background(), "ORG/PROJ/REPOID", 1)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestPullListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/pullRequests/1/commits").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client := NewDefault()
	got, _, err := client.PullRequests.ListCommits(context.Background(), "ORG/PROJ/REPOID", 1, &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Commit{}
	raw, _ := os.ReadFile("testdata/commits.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
