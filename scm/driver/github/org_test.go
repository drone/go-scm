// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestOrganizationFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/orgs/github").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/org.json")

	client := NewDefault()
	got, res, err := client.Organizations.Find(context.Background(), "github")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Organization)
	raw, _ := ioutil.ReadFile("testdata/org.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestOrganizationList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/user/orgs").
		MatchParam("per_page", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/orgs.json")

	client := NewDefault()
	got, res, err := client.Organizations.List(context.Background(), scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Organization{}
	raw, _ := ioutil.ReadFile("testdata/orgs.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestTeamList(t *testing.T) {
	defer gock.Off()

	org := "myorg"

	gock.New("https://api.github.com").
		Get("/orgs/myorg/teams").
		MatchParam("per_page", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/teams.json")

	client := NewDefault()
	got, res, err := client.Organizations.ListTeams(context.Background(), org, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Team{}
	raw, _ := ioutil.ReadFile("testdata/teams.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestTeamMembers(t *testing.T) {
	defer gock.Off()

	teamID := 1
	role := "all"

	gock.New("https://api.github.com").
		Get("/teams/1/members").
		MatchParam("role", role).
		MatchParam("per_page", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/team_members.json")

	client := NewDefault()
	got, res, err := client.Organizations.ListTeamMembers(context.Background(), teamID, role, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.TeamMember{}
	raw, _ := ioutil.ReadFile("testdata/team_members.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}
