// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestOrgMembers(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/orgs/myorg/members").
		MatchParam("per_page", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/org_members.json")

	client := NewDefault()
	got, res, err := client.Organizations.ListOrgMembers(context.Background(), "myorg", scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.TeamMember{}
	raw, _ := ioutil.ReadFile("testdata/org_members.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestIsAdmin(t *testing.T) {
	defer gock.Off()

	testOrg := "testOrg"
	testUser := "testUser"

	gock.New("https://api.github.com").
		Get(fmt.Sprintf("/orgs/%s/memberships/%s", testOrg, testUser)).
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/membership_admin.json")

	client := NewDefault()
	isAdmin, _, err := client.Organizations.IsAdmin(context.Background(), testOrg, testUser)
	if err != nil {
		t.Error(err)
		return
	}

	assert.True(t, isAdmin)
}

func TestIsAdminFalse(t *testing.T) {
	defer gock.Off()

	testOrg := "testOrg"
	testUser := "testUser"

	gock.New("https://api.github.com").
		Get(fmt.Sprintf("/orgs/%s/memberships/%s", testOrg, testUser)).
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/membership_member.json")

	client := NewDefault()
	isAdmin, _, err := client.Organizations.IsAdmin(context.Background(), testOrg, testUser)
	if err != nil {
		t.Error(err)
		return
	}

	assert.False(t, isAdmin)
}

func TestOrganizationService_IsMember_False(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/orgs/octocat/members/someuser").
		Reply(404).
		SetHeaders(mockHeaders)

	client := NewDefault()
	got, res, err := client.Organizations.IsMember(context.Background(), "octocat", "someuser")
	if err != nil {
		t.Error(err)
		return
	}

	if got {
		t.Errorf("Expected user to not be a member")
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestOrganizationService_IsMember_True(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/orgs/octocat/members/someuser").
		Reply(204).
		SetHeaders(mockHeaders)

	client := NewDefault()
	got, res, err := client.Organizations.IsMember(context.Background(), "octocat", "someuser")
	if err != nil {
		t.Error(err)
		return
	}

	if !got {
		t.Errorf("Expected user to be a member")
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestListPendingInvitations(t *testing.T) {
	defer gock.Off()

	testOrg := "testOrg"
	testUser := "monalisa"

	gock.New("https://api.github.com").
		Get(fmt.Sprintf("/orgs/%s/invitations", testOrg)).
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/list_pending_invitations.json")

	client := NewDefault()
	invites, res, err := client.Organizations.ListPendingInvitations(context.Background(), testOrg, scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 1, len(invites), "should have found one pending invite")
	assert.Equal(t, testUser, invites[0].Login, fmt.Sprintf("should have found a pending invite for user %s", testUser))

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestAcceptOrganizationInvitation(t *testing.T) {
	defer gock.Off()

	testOrg := "octocat"

	gock.New("https://api.github.com").
		Patch(fmt.Sprintf("/user/memberships/orgs/%s", testOrg)).
		Reply(200).
		File("testdata/org_accept_invitation.json").
		SetHeader("X-GitHub-Request-Id", "DD0E:6011:12F21A8:1926790:5A2064E2").
		SetHeader("X-RateLimit-Limit", "60").
		SetHeader("X-RateLimit-Remaining", "59").
		SetHeader("X-RateLimit-Reset", "1512076018")

	client := NewDefault()
	res, err := client.Organizations.AcceptOrganizationInvitation(context.Background(), testOrg)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestListMemberships(t *testing.T) {
	defer gock.Off()

	testOrg1 := "github"
	testState1 := "active"
	testRole1 := "admin"
	testOrg2 := "jenkins-x"
	testState2 := "pending"
	testRole2 := "admin"

	gock.New("https://api.github.com").
		Get("/user/memberships/orgs").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/list_memberships.json")

	client := NewDefault()
	memberships, res, err := client.Organizations.ListMemberships(context.Background(), scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 2, len(memberships), "should have found two memberships")
	for _, m := range memberships {
		switch m.OrganizationName {
		case testOrg1:
			assert.Equal(t, testState1, m.State, fmt.Sprintf("should have found a state %s", testState1))
			assert.Equal(t, testRole1, m.Role, fmt.Sprintf("should have found a role %s", testRole1))

		case testOrg2:
			assert.Equal(t, testState2, m.State, fmt.Sprintf("should have found a state %s", testState2))
			assert.Equal(t, testRole2, m.Role, fmt.Sprintf("should have found a role %s", testRole2))
		default:
			t.Error(fmt.Errorf("unrecognised organisation name %s", m.OrganizationName))
			return
		}
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
