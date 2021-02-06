// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"io/ioutil"
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

	got, res, err := client.Organizations.List(context.Background(), scm.ListOptions{Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	var want []*scm.Organization
	raw, _ := ioutil.ReadFile("testdata/orgs.json.golden")
	json.Unmarshal(raw, &want)

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

	got, _, err := client.Organizations.ListOrgMembers(context.Background(), "some-project", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	var want []*scm.TeamMember
	raw, _ := ioutil.ReadFile("testdata/org_members.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestOrganizationIsMember(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Times(2).
		Get("/rest/api/1.0/projects/some-project/permissions/users").
		Reply(200).
		Type("application/json").
		File("testdata/org_members.json")

	client, _ := New("http://example.com:7990")

	got, _, err := client.Organizations.IsMember(context.Background(), "some-project", "jcitizen")
	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(got, true); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	got, _, err = client.Organizations.IsMember(context.Background(), "some-project", "not-present")
	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(got, false); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
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
