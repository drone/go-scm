// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestGitCreateBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(201).
		Type("application/json").
		File("testdata/branch_create.json")

	params := &scm.CreateBranch{
		Name: "test_branch",
		Sha:  "312797ba52425353dec56871a255e2a36fc96344",
	}

	client := NewDefault("ORG", "PROJ")
	res, err := client.Git.CreateBranch(context.Background(), "REPOID", params)

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

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Git.ListCommits(context.Background(), "REPOID", scm.CommitListOptions{})
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
}

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/REPOID/").
		Reply(200).
		Type("application/json").
		File("testdata/branches.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Git.ListBranches(context.Background(), "REPOID", scm.ListOptions{})
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
}
