// Copyright 2026 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/drone/go-scm/scm"
)

// TestList_AllWorkspacesEmpty tests when all workspaces have zero repositories
func TestList_AllWorkspacesEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws1"}},
					{"workspace": map[string]string{"slug": "ws2"}},
				},
				"next": "",
			})
			return
		}

		// Both workspaces have 0 repos
		if r.URL.Path == "/2.0/repositories/ws1" || r.URL.Path == "/2.0/repositories/ws2" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{},
				"size":   0,
				"next":   "",
			})
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	repos, res, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(repos) != 0 {
		t.Errorf("Expected 0 repos, got %d", len(repos))
	}

	if res.Page.Next != 0 {
		t.Errorf("Expected no next page, got %d", res.Page.Next)
	}
}

// TestList_NoWorkspaces tests when user has no workspaces
func TestList_NoWorkspaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{},
				"next":   "",
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	repos, _, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(repos) != 0 {
		t.Errorf("Expected 0 repos when no workspaces, got %d", len(repos))
	}
}

// TestList_WorkspaceFetchError tests error handling when workspace fetch fails
func TestList_WorkspaceFetchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.WriteHeader(500)
			w.Write([]byte(`{"error": {"message": "Internal server error"}}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	_, _, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err == nil {
		t.Fatal("Expected error when workspace fetch fails")
	}
}

// TestList_SizeZeroDefaultsTo100 tests that size=0 defaults to 100
func TestList_SizeZeroDefaultsTo100(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws1"}},
				},
				"next": "",
			})
			return
		}

		if r.URL.Path == "/2.0/repositories/ws1" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", 0, 150, "", 150))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Size=0 should default to 100
	repos, res, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 0})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(repos) != 100 {
		t.Errorf("Expected 100 repos (default size), got %d", len(repos))
	}

	if res.Page.Next != 2 {
		t.Errorf("Expected next page 2, got %d", res.Page.Next)
	}
}

// TestList_VeryLargePageSize tests large page size
func TestList_LargePageSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws1"}},
				},
				"next": "",
			})
			return
		}

		if r.URL.Path == "/2.0/repositories/ws1" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", 0, 50, "", 50))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Request 1000 repos but only 50 exist
	repos, res, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 1000})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(repos) != 50 {
		t.Errorf("Expected 50 repos (all available), got %d", len(repos))
	}

	if res.Page.Next != 0 {
		t.Errorf("Expected no next page, got %d", res.Page.Next)
	}
}

// TestList_ExactPageBoundary tests when workspace repos end exactly at page boundary
func TestList_ExactPageBoundary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws1"}},
					{"workspace": map[string]string{"slug": "ws2"}},
				},
				"next": "",
			})
			return
		}

		// ws1 has exactly 100 repos
		if r.URL.Path == "/2.0/repositories/ws1" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", 0, 100, "", 100))
			return
		}

		// ws2 has 50 repos
		if r.URL.Path == "/2.0/repositories/ws2" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws2", 0, 50, "", 50))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Page 1: exactly 100 from ws1
	repos1, res1, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err != nil {
		t.Fatalf("Page 1 error: %v", err)
	}

	if len(repos1) != 100 {
		t.Errorf("Page 1: expected 100 repos, got %d", len(repos1))
	}

	if res1.Page.Next != 2 {
		t.Errorf("Page 1: expected next page 2, got %d", res1.Page.Next)
	}

	// Page 2: 50 from ws2
	repos2, res2, err := client.Repositories.List(ctx, scm.ListOptions{Page: 2, Size: 100})
	if err != nil {
		t.Fatalf("Page 2 error: %v", err)
	}

	if len(repos2) != 50 {
		t.Errorf("Page 2: expected 50 repos, got %d", len(repos2))
	}

	if res2.Page.Next != 0 {
		t.Errorf("Page 2: expected no next page, got %d", res2.Page.Next)
	}
}

// TestList_PartialWorkspaceAccessErrors tests resilience when some workspaces are inaccessible
func TestList_SkipsInaccessibleWorkspaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws1"}},
					{"workspace": map[string]string{"slug": "ws2"}},
					{"workspace": map[string]string{"slug": "ws3"}},
				},
				"next": "",
			})
			return
		}

		// ws1: accessible, 50 repos
		if r.URL.Path == "/2.0/repositories/ws1" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", 0, 50, "", 50))
			return
		}

		// ws2: inaccessible (403)
		if r.URL.Path == "/2.0/repositories/ws2" {
			w.WriteHeader(403)
			w.Write([]byte(`{"error": {"message": "Access denied"}}`))
			return
		}

		// ws3: accessible, 30 repos
		if r.URL.Path == "/2.0/repositories/ws3" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws3", 0, 30, "", 30))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	repos, _, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should get 50 from ws1 + 30 from ws3 = 80 repos (skipping ws2)
	if len(repos) != 80 {
		t.Errorf("Expected 80 repos (skipping inaccessible ws2), got %d", len(repos))
	}

	// Verify repos are from correct workspaces
	ws1Count := 0
	ws3Count := 0
	for _, repo := range repos {
		if repo.Namespace == "ws1" {
			ws1Count++
		} else if repo.Namespace == "ws3" {
			ws3Count++
		}
	}

	if ws1Count != 50 {
		t.Errorf("Expected 50 repos from ws1, got %d", ws1Count)
	}

	if ws3Count != 30 {
		t.Errorf("Expected 30 repos from ws3, got %d", ws3Count)
	}
}

// TestList_SingleRepoMultipleWorkspaces tests correct distribution
func TestList_SmallPageAcrossWorkspaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws1"}},
					{"workspace": map[string]string{"slug": "ws2"}},
					{"workspace": map[string]string{"slug": "ws3"}},
				},
				"next": "",
			})
			return
		}

		if r.URL.Path == "/2.0/repositories/ws1" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", 0, 3, "", 3))
			return
		}

		if r.URL.Path == "/2.0/repositories/ws2" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws2", 0, 3, "", 3))
			return
		}

		if r.URL.Path == "/2.0/repositories/ws3" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws3", 0, 3, "", 3))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Request page size of 5 (should get 3 from ws1, 2 from ws2)
	repos, res, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 5})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(repos) != 5 {
		t.Errorf("Expected 5 repos, got %d", len(repos))
	}

	// Should have next page since 9 total > 5
	if res.Page.Next != 2 {
		t.Errorf("Expected next page 2, got %d", res.Page.Next)
	}

	// Verify distribution
	ws1Count := 0
	ws2Count := 0
	for _, repo := range repos {
		if repo.Namespace == "ws1" {
			ws1Count++
		} else if repo.Namespace == "ws2" {
			ws2Count++
		}
	}

	if ws1Count != 3 {
		t.Errorf("Expected 3 repos from ws1, got %d", ws1Count)
	}

	if ws2Count != 2 {
		t.Errorf("Expected 2 repos from ws2, got %d", ws2Count)
	}
}

// TestList_OffsetCalculation tests various offset scenarios
func TestList_OffsetCalculationAcrossPages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws1"}},
				},
				"next": "",
			})
			return
		}

		if r.URL.Path == "/2.0/repositories/ws1" {
			page := r.URL.Query().Get("page")
			pagelenStr := r.URL.Query().Get("pagelen")

			// Check for count query (pagelen=1)
			if pagelenStr == "1" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", 0, 1, "", 500))
				return
			}

			// Parse pagelen, default to 100 if not specified
			pagelen := 100
			if pagelenStr != "" {
				if n, err := strconv.Atoi(pagelenStr); err == nil && n > 0 {
					pagelen = n
				}
			}

			// Parse page, default to 1
			pageNum := 1
			if page != "" {
				if n, err := strconv.Atoi(page); err == nil && n > 0 {
					pageNum = n
				}
			}

			// Calculate start index based on page and pagelen
			startIdx := (pageNum - 1) * pagelen
			count := pagelen
			if startIdx+count > 500 {
				count = 500 - startIdx
			}

			// Calculate next page
			nextPage := ""
			if startIdx+count < 500 {
				nextPage = strconv.Itoa(pageNum + 1)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", startIdx, count, nextPage, 500))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Page 3 with size 50 (offset = (3-1)*50 = 100, should get repos 100-149)
	repos, res, err := client.Repositories.List(ctx, scm.ListOptions{Page: 3, Size: 50})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(repos) != 50 {
		t.Errorf("Expected 50 repos, got %d", len(repos))
	}

	// Debug: print what we got
	if len(repos) > 0 {
		t.Logf("First repo: %s, Last repo: %s", repos[0].Name, repos[len(repos)-1].Name)
	}

	// Should start from repo-100 (offset 100)
	if len(repos) > 0 && repos[0].Name != "repo-100" {
		t.Errorf("Expected first repo to be repo-100, got %s", repos[0].Name)
	}

	// Should end at repo-149
	if len(repos) >= 50 && repos[49].Name != "repo-149" {
		t.Errorf("Expected last repo to be repo-149, got %s", repos[49].Name)
	}

	if res.Page.Next != 4 {
		t.Errorf("Expected next page 4, got %d", res.Page.Next)
	}
}

// TestListV2_EmptySearchResults tests search with no results
func TestListV2_NoSearchResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws1"}},
				},
				"next": "",
			})
			return
		}

		if r.URL.Path == "/2.0/repositories/ws1" {
			// Verify search query
			query := r.URL.Query().Get("q")
			if query != `name~"nonexistent"` {
				t.Errorf("Expected search query name~\"nonexistent\", got: %s", query)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{},
				"size":   0,
				"next":   "",
			})
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	repos, res, err := client.Repositories.ListV2(ctx, scm.RepoListOptions{
		RepoSearchTerm: scm.RepoSearchTerm{RepoName: "nonexistent"},
		ListOptions:    scm.ListOptions{Page: 1, Size: 100},
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(repos) != 0 {
		t.Errorf("Expected 0 repos for non-matching search, got %d", len(repos))
	}

	if res.Page.Next != 0 {
		t.Errorf("Expected no next page, got %d", res.Page.Next)
	}
}

// Note: createMockRepoPageWithSize is defined in repo_workspace_pagination_test.go
