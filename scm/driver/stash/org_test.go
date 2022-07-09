// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

func TestOrganizationFind(t *testing.T) {
	client, _ := New("https://api.bitbucket.org")
	_, _, err := client.Organizations.Find(context.Background(), "atlassian")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestOrganizationList(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/orgs.json")

	client, _ := New("http://example.com:7990")

	got, res, err := client.Organizations.List(context.Background(), &scm.ListOptions{Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	var want []*scm.Organization
	raw, _ := os.ReadFile("testdata/orgs.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
}

func TestOrganizationListOrgMembers(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/some-project/permissions/users").
		Reply(200).
		Type("application/json").
		File("testdata/org_members.json")

	client, _ := New("http://example.com:7990")

	got, _, err := client.Organizations.ListOrgMembers(context.Background(), "some-project", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	var want []*scm.TeamMember
	raw, _ := os.ReadFile("testdata/org_members.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestOrganizationIsMember(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Times(3).
		Get("/rest/api/1.0/projects/some-project/permissions/users").
		Reply(200).
		Type("application/json").
		File("testdata/org_members.json")

	gock.New("http://example.com:7990").
		Times(3).
		Get("rest/api/1.0/projects/some-project/permissions/groups").
		Reply(200).
		Type("application/json").
		File("testdata/project_groups.json")

	gock.New("http://example.com:7990").
		Times(3).
		Get("rest/api/1.0/admin/groups/more-members").
		ParamPresent("context").
		Reply(200).
		Type("application/json").
		File("testdata/group_members.json")

	testCases := []struct {
		description string
		user        string
		isMember    bool
	}{
		{
			description: "user assigned directly to project",
			user:        "jcitizen",
			isMember:    true,
		},
		{
			description: "user part of group assigned directly to project",
			user:        "jx-user",
			isMember:    true,
		},
		{
			description: "user not assigned to project",
			user:        "not-present",
			isMember:    false,
		},
	}

	for k, v := range testCases {
		t.Logf("Running test %q: %s", k, v.description)
		client, _ := New("http://example.com:7990")

		got, _, err := client.Organizations.IsMember(context.Background(), "some-project", v.user)
		if err != nil {
			t.Error(err)
			return
		}

		if diff := cmp.Diff(got, v.isMember); diff != "" {
			t.Errorf("Unexpected Results")
			t.Log(diff)
		}
	}
}

func TestOrganizationIsAdmin(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Times(2).
		Get("/rest/api/1.0/projects/some-project/permissions/users").
		Reply(200).
		Type("application/json").
		File("testdata/org_members.json")

	client, _ := New("http://example.com:7990")

	got, _, err := client.Organizations.IsAdmin(context.Background(), "some-project", "jcitizen")
	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(got, true); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	got, _, err = client.Organizations.IsAdmin(context.Background(), "some-project", "bob")
	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(got, false); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
