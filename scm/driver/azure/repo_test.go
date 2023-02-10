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

func TestRepositoryListOrgAndProject(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories").
		Reply(200).
		Type("application/json").
		File("testdata/repos.json")

	client := NewDefault()
	got, _, err := client.Repositories.ListOrganisation(context.Background(), "ORG/PROJ", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := os.ReadFile("testdata/repos.json.golden")
	jsonErr := json.Unmarshal(raw, &want)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryListOrgOnly(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/_apis/git/repositories").
		Reply(200).
		Type("application/json").
		File("testdata/repos.json")

	client := NewDefault()
	got, _, err := client.Repositories.ListOrganisation(context.Background(), "ORG", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := os.ReadFile("testdata/repos.json.golden")
	jsonErr := json.Unmarshal(raw, &want)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/test_project").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	client := NewDefault()
	got, _, err := client.Repositories.Find(context.Background(), "ORG/PROJ/test_project")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := os.ReadFile("testdata/repo.json.golden")
	jsonErr := json.Unmarshal(raw, &want)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/_apis/projects/PROJ").
		Reply(200).
		Type("application/json").
		File("testdata/project.json")

	gock.New("https://dev.azure.com/").
		Post("/ORG/PROJ/_apis/git/repositories").
		File("testdata/repo_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	client := NewDefault()

	input := &scm.RepositoryInput{Name: "test_project", Namespace: "ORG/PROJ"}

	got, _, err := client.Repositories.Create(context.Background(), input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := os.ReadFile("testdata/repo.json.golden")
	jsonErr := json.Unmarshal(raw, &want)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories/test_project").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	gock.New("https://dev.azure.com/").
		Delete("/ORG/PROJ/_apis/git/repositories/91f0d4cb-4c36-49a5-b28d-2d72da089c4d").
		Reply(204)

	client := NewDefault()

	_, err := client.Repositories.Delete(context.Background(), "ORG/PROJ/test_project")
	if err != nil {
		t.Error(err)
		return
	}
}
