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

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/commit.json")

	client := NewDefault()

	got, _, err := client.Git.FindCommit(context.Background(), "ORG/PROJ/REPOID", "14897f4465d2d63508242b5cbf68aa2865f693e7")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Commit)
	raw, _ := os.ReadFile("testdata/commit.json.golden")

	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitCreateRef(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(201).
		Type("application/json").
		File("testdata/branch_create.json")

	params := &scm.ReferenceInput{
		Name: "test_branch",
		Sha:  "312797ba52425353dec56871a255e2a36fc96344",
	}

	client := NewDefault()
	got, res, err := client.Git.CreateRef(context.Background(), "ORG/PROJ/REPOID", params.Name, params.Sha)

	if err != nil {
		t.Error(err)
		return
	}

	if res.Status != 201 {
		t.Errorf("Unexpected Results")
	}

	want := &scm.Reference{}
	raw, _ := os.ReadFile("testdata/branch_create.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestGitListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client := NewDefault()
	got, _, err := client.Git.ListCommits(context.Background(), "ORG/PROJ/REPOID", scm.CommitListOptions{})
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

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/branches.json")

	client := NewDefault()
	got, _, err := client.Git.ListBranches(context.Background(), "ORG/PROJ/REPOID", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := os.ReadFile("testdata/branches.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitCompareCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/compare.json")

	client := NewDefault()
	got, _, err := client.Git.CompareCommits(context.Background(), "ORG/PROJ/REPOID", "9788e5ddf8b387cb79228628f34d8dc18582d606", "66df312dad61e84dd896d1e8d14ee3dce53b62f0", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := os.ReadFile("testdata/compare.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}
