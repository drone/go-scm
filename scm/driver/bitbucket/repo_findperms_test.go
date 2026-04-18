// Copyright 2026 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"testing"

	"github.com/h2non/gock"
)

// TestFindPerms_ErrorWhenAllWorkspacesEmpty tests the case where repo doesn't exist in any workspace
func TestFindPerms_ErrorWhenAllWorkspacesEmpty(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("page", "1").
		MatchParam("pagelen", "100").
		Reply(200).
		Type("application/json").
		BodyString(`{"values": [{"workspace": {"slug": "ws1"}}, {"workspace": {"slug": "ws2"}}], "next": ""}`)

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/ws1/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_empty.json")

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/ws2/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_empty.json")

	client, _ := New("https://api.bitbucket.org")
	_, _, err := client.Repositories.FindPerms(context.Background(), "nonexistent-repo")

	if err == nil {
		t.Fatal("Expected error when repository not found in any workspace")
	}

	if err.Error() != "repository nonexistent-repo not found in any workspace" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

// TestFindPerms_ErrorWhenWorkspaceFetchFails tests error handling when workspace fetch fails
func TestFindPerms_ErrorWhenWorkspaceFetchFails(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("page", "1").
		MatchParam("pagelen", "100").
		Reply(500).
		Type("application/json").
		BodyString(`{"error": {"message": "Internal server error"}}`)

	client, _ := New("https://api.bitbucket.org")
	_, _, err := client.Repositories.FindPerms(context.Background(), "test-repo")

	if err == nil {
		t.Fatal("Expected error when workspace fetch fails")
	}
}

// TestFindPerms_NoPermissionsInAnyWorkspace tests when user has no permissions
func TestFindPerms_NoPermissionsInAnyWorkspace(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("page", "1").
		MatchParam("pagelen", "100").
		Reply(200).
		Type("application/json").
		BodyString(`{"values": [{"workspace": {"slug": "ws1"}}], "next": ""}`)

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/ws1/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_empty.json")

	client, _ := New("https://api.bitbucket.org")
	_, _, err := client.Repositories.FindPerms(context.Background(), "test-repo")

	if err == nil {
		t.Fatal("Expected error when user has no permissions")
	}

	if err.Error() != "repository test-repo not found in any workspace" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

// TestFindPerms_NetworkErrorReturnsImmediately tests that API errors are returned immediately
func TestFindPerms_NetworkErrorReturnsImmediately(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("page", "1").
		MatchParam("pagelen", "100").
		Reply(200).
		Type("application/json").
		BodyString(`{"values": [{"workspace": {"slug": "ws1"}}, {"workspace": {"slug": "ws2"}}], "next": ""}`)

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/ws1/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(500).
		Type("application/json").
		BodyString(`{"error": {"message": "Internal server error"}}`)

	client, _ := New("https://api.bitbucket.org")
	_, res, err := client.Repositories.FindPerms(context.Background(), "test-repo")

	if err == nil {
		t.Fatal("Expected error for server failure")
	}

	if res == nil || res.Status != 500 {
		t.Error("Expected response with 500 status")
	}
}

// TestFindPerms_ContinuesToNextWorkspace tests that iteration continues across workspaces
func TestFindPerms_ContinuesToNextWorkspace(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("page", "1").
		MatchParam("pagelen", "100").
		Reply(200).
		Type("application/json").
		BodyString(`{"values": [{"workspace": {"slug": "ws1"}}, {"workspace": {"slug": "ws2"}}], "next": ""}`)

	// ws1: no access
	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/ws1/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_empty.json")

	// ws2: admin
	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/ws2/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_admin.json")

	client, _ := New("https://api.bitbucket.org")
	perm, _, err := client.Repositories.FindPerms(context.Background(), "test-repo")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !perm.Admin || !perm.Push || !perm.Pull {
		t.Error("Expected admin permissions from ws2")
	}
}

// TestFindPerms_EmptyWorkspaceList tests when user has no workspaces
func TestFindPerms_EmptyWorkspaceList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("page", "1").
		MatchParam("pagelen", "100").
		Reply(200).
		Type("application/json").
		BodyString(`{"values": [], "next": ""}`)

	client, _ := New("https://api.bitbucket.org")
	_, _, err := client.Repositories.FindPerms(context.Background(), "test-repo")

	if err == nil {
		t.Fatal("Expected error when no workspaces available")
	}

	if err.Error() != "repository test-repo not found in any workspace" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

// TestFindPerms_WorkspaceInURLExtractionWithSlash tests that URL workspace is used
func TestFindPerms_WorkspaceInURLExtractionWithSlash(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/url-workspace/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_admin.json")

	client, _ := New("https://api.bitbucket.org/repositories/url-workspace")
	perm, _, err := client.Repositories.FindPerms(context.Background(), "different-workspace/actual-repo")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !perm.Admin {
		t.Error("Expected admin permissions")
	}
}

// TestFindPerms_URLWorkspaceTakesPriority tests that URL workspace takes priority
func TestFindPerms_URLWorkspaceTakesPriority(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/url-workspace/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_admin.json")

	client, _ := New("https://api.bitbucket.org/repositories/url-workspace")
	perm, _, err := client.Repositories.FindPerms(context.Background(), "identifier-workspace/repo-slug")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !perm.Admin {
		t.Error("Expected admin permissions from url-workspace")
	}
}

// TestFindPerms_WorkspacePagination tests pagination in workspace fetching
func TestFindPerms_WorkspacePagination(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("page", "1").
		MatchParam("pagelen", "100").
		Reply(200).
		Type("application/json").
		BodyString(`{"values": [{"workspace": {"slug": "ws1"}}], "next": "https://api.bitbucket.org/2.0/user/workspaces?page=2"}`)

	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces").
		MatchParam("page", "2").
		Reply(200).
		Type("application/json").
		BodyString(`{"values": [{"workspace": {"slug": "ws2"}}], "next": ""}`)

	// ws1: no access
	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/ws1/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_empty.json")

	// ws2: write
	gock.New("https://api.bitbucket.org").
		Get("/2.0/user/workspaces/ws2/permissions/repositories").
		MatchParam("pagelen", "1").
		Reply(200).
		Type("application/json").
		File("testdata/user_repo_perm_write.json")

	client, _ := New("https://api.bitbucket.org")
	perm, _, err := client.Repositories.FindPerms(context.Background(), "test-repo")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !perm.Push || !perm.Pull || perm.Admin {
		t.Error("Expected write permissions from ws2")
	}
}
