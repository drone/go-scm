// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestOrganizationFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/workspaces/atlassian").
		Reply(200).
		Type("application/json").
		File("testdata/team.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Organizations.Find(context.Background(), "atlassian")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Organization)
	raw, _ := ioutil.ReadFile("testdata/team.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestOrganizationList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("pagelen", "30").
		MatchParam("page", "1").
		MatchParam("role", "member").
		Reply(200).
		Type("application/json").
		File("testdata/user_workspaces.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.Organizations.List(context.Background(), scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Organization{}
	raw, _ := ioutil.ReadFile("testdata/user_workspaces.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestConvertWorkspace(t *testing.T) {
	tests := []struct {
		name      string
		workspace *workspace
		want      *scm.Organization
	}{
		{
			name: "workspace with avatar link",
			workspace: &workspace{
				Slug: "my-workspace",
				Links: struct {
					Avatar link `json:"avatar"`
				}{
					Avatar: link{Href: "https://bitbucket.org/account/my-workspace/avatar/32/"},
				},
			},
			want: &scm.Organization{
				Name:   "my-workspace",
				Avatar: "https://bitbucket.org/account/my-workspace/avatar/32/",
			},
		},
		{
			name: "workspace without avatar link",
			workspace: &workspace{
				Slug: "test-workspace",
				Links: struct {
					Avatar link `json:"avatar"`
				}{
					Avatar: link{Href: ""},
				},
			},
			want: &scm.Organization{
				Name:   "test-workspace",
				Avatar: "https://bitbucket.org/account/test-workspace/avatar/32",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertWorkspace(tt.workspace)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("convertWorkspace() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertWorkspaceAccessList(t *testing.T) {
	tests := []struct {
		name string
		from *workspaceAccessList
		want []*scm.Organization
	}{
		{
			name: "valid workspaces",
			from: &workspaceAccessList{
				Values: []*workspaceAccess{
					{
						Workspace: &workspace{
							Slug: "workspace1",
							Links: struct {
								Avatar link `json:"avatar"`
							}{
								Avatar: link{Href: "https://bitbucket.org/account/workspace1/avatar/32/"},
							},
						},
					},
					{
						Workspace: &workspace{
							Slug: "workspace2",
							Links: struct {
								Avatar link `json:"avatar"`
							}{
								Avatar: link{Href: ""},
							},
						},
					},
				},
			},
			want: []*scm.Organization{
				{
					Name:   "workspace1",
					Avatar: "https://bitbucket.org/account/workspace1/avatar/32/",
				},
				{
					Name:   "workspace2",
					Avatar: "https://bitbucket.org/account/workspace2/avatar/32",
				},
			},
		},
		{
			name: "nil workspace",
			from: &workspaceAccessList{
				Values: []*workspaceAccess{
					{
						Workspace: nil,
					},
				},
			},
			want: []*scm.Organization{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertWorkspaceAccessList(tt.from)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("convertWorkspaceAccessList() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
