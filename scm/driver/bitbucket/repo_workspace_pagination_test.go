// Copyright 2026 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drone/go-scm/scm"
)

// TestList_WorkspacePagination_SingleWorkspace tests pagination within a single workspace
func TestList_WorkspacePagination_SingleWorkspace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock workspace list endpoint (sorted)
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "workspace1"}},
				},
				"next": "",
			})
			return
		}

		// Mock repository list for workspace1 (250 total repos, paginated)
		if r.URL.Path == "/2.0/repositories/workspace1" {
			page := r.URL.Query().Get("page")
			pagelen := r.URL.Query().Get("pagelen")

			if pagelen == "" {
				pagelen = "10" // default
			}

			switch page {
			case "", "1":
				// Page 1: repos 0-99 (include total size)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(createMockRepoPageWithSize("workspace1", 0, 100, "2", 250))
			case "2":
				// Page 2: repos 100-199
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(createMockRepoPageWithSize("workspace1", 100, 100, "3", 250))
			case "3":
				// Page 3: repos 200-249 (last page, partial)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(createMockRepoPageWithSize("workspace1", 200, 50, "", 250))
			}
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Page 1: Should get 100 repos
	repos1, res1, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err != nil {
		t.Fatalf("Page 1 error: %v", err)
	}
	if len(repos1) != 100 {
		t.Errorf("Page 1: expected 100 repos, got %d", len(repos1))
	}
	if res1.Page.Next != 2 {
		t.Errorf("Page 1: expected Next=2, got %d", res1.Page.Next)
	}

	// Page 2: Should get 100 repos
	repos2, res2, err := client.Repositories.List(ctx, scm.ListOptions{Page: 2, Size: 100})
	if err != nil {
		t.Fatalf("Page 2 error: %v", err)
	}
	if len(repos2) != 100 {
		t.Errorf("Page 2: expected 100 repos, got %d", len(repos2))
	}
	if res2.Page.Next != 3 {
		t.Errorf("Page 2: expected Next=3, got %d", res2.Page.Next)
	}

	// Page 3: Should get 50 repos (partial page)
	repos3, res3, err := client.Repositories.List(ctx, scm.ListOptions{Page: 3, Size: 100})
	if err != nil {
		t.Fatalf("Page 3 error: %v", err)
	}
	if len(repos3) != 50 {
		t.Errorf("Page 3: expected 50 repos, got %d", len(repos3))
	}
	if res3.Page.Next != 0 {
		t.Errorf("Page 3: expected Next=0 (no more pages), got %d", res3.Page.Next)
	}
}

// TestList_WorkspacePagination_CrossWorkspace tests pagination across multiple workspaces
func TestList_WorkspacePagination_CrossWorkspace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock workspace list endpoint (sorted)
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws-alpha"}},
					{"workspace": map[string]string{"slug": "ws-beta"}},
					{"workspace": map[string]string{"slug": "ws-gamma"}},
				},
				"next": "",
			})
			return
		}

		// ws-alpha: 150 repos (positions 0-149)
		if r.URL.Path == "/2.0/repositories/ws-alpha" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws-alpha", 0, 150, "", 150))
			return
		}

		// ws-beta: 120 repos (positions 150-269)
		if r.URL.Path == "/2.0/repositories/ws-beta" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws-beta", 0, 120, "", 120))
			return
		}

		// ws-gamma: 80 repos (positions 270-349)
		if r.URL.Path == "/2.0/repositories/ws-gamma" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws-gamma", 0, 80, "", 80))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Page 1 (offset 0-99): All from ws-alpha
	repos1, res1, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err != nil {
		t.Fatalf("Page 1 error: %v", err)
	}
	if len(repos1) != 100 {
		t.Errorf("Page 1: expected 100 repos, got %d", len(repos1))
	}
	if repos1[0].Namespace != "ws-alpha" {
		t.Errorf("Page 1: expected first repo from ws-alpha, got %s", repos1[0].Namespace)
	}
	if res1.Page.Next != 2 {
		t.Errorf("Page 1: expected Next=2, got %d", res1.Page.Next)
	}

	// Page 2 (offset 100-199): 50 from ws-alpha + 50 from ws-beta (boundary crossing!)
	repos2, res2, err := client.Repositories.List(ctx, scm.ListOptions{Page: 2, Size: 100})
	if err != nil {
		t.Fatalf("Page 2 error: %v", err)
	}
	if len(repos2) != 100 {
		t.Errorf("Page 2: expected 100 repos, got %d", len(repos2))
	}

	// Verify workspace boundary crossing
	alphaCount := 0
	betaCount := 0
	for _, repo := range repos2 {
		if repo.Namespace == "ws-alpha" {
			alphaCount++
		} else if repo.Namespace == "ws-beta" {
			betaCount++
		}
	}
	if alphaCount != 50 {
		t.Errorf("Page 2: expected 50 from ws-alpha, got %d", alphaCount)
	}
	if betaCount != 50 {
		t.Errorf("Page 2: expected 50 from ws-beta, got %d", betaCount)
	}
	if res2.Page.Next != 3 {
		t.Errorf("Page 2: expected Next=3, got %d", res2.Page.Next)
	}

	// Page 3 (offset 200-299): 70 from ws-beta + 30 from ws-gamma
	repos3, res3, err := client.Repositories.List(ctx, scm.ListOptions{Page: 3, Size: 100})
	if err != nil {
		t.Fatalf("Page 3 error: %v", err)
	}
	if len(repos3) != 100 {
		t.Errorf("Page 3: expected 100 repos, got %d", len(repos3))
	}

	// Verify workspace boundary crossing
	betaCount = 0
	gammaCount := 0
	for _, repo := range repos3 {
		if repo.Namespace == "ws-beta" {
			betaCount++
		} else if repo.Namespace == "ws-gamma" {
			gammaCount++
		}
	}
	if betaCount != 70 {
		t.Errorf("Page 3: expected 70 from ws-beta, got %d", betaCount)
	}
	if gammaCount != 30 {
		t.Errorf("Page 3: expected 30 from ws-gamma, got %d", gammaCount)
	}
	if res3.Page.Next != 4 {
		t.Errorf("Page 3: expected Next=4, got %d", res3.Page.Next)
	}

	// Page 4 (offset 300-349): 50 from ws-gamma (partial page, last page)
	repos4, res4, err := client.Repositories.List(ctx, scm.ListOptions{Page: 4, Size: 100})
	if err != nil {
		t.Fatalf("Page 4 error: %v", err)
	}
	if len(repos4) != 50 {
		t.Errorf("Page 4: expected 50 repos (partial), got %d", len(repos4))
	}
	if repos4[0].Namespace != "ws-gamma" {
		t.Errorf("Page 4: expected repos from ws-gamma, got %s", repos4[0].Namespace)
	}
	if res4.Page.Next != 0 {
		t.Errorf("Page 4: expected Next=0 (no more pages), got %d", res4.Page.Next)
	}
}

// TestList_DefaultPage tests that Page=0 defaults to page 1
func TestList_DefaultPage(t *testing.T) {
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

		if r.URL.Path == "/2.0/repositories/ws1" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", 0, 50, "", 50))
			return
		}

		if r.URL.Path == "/2.0/repositories/ws2" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws2", 0, 30, "", 30))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Page=0 should default to page 1 and return first page worth of repos
	repos, res, err := client.Repositories.List(ctx, scm.ListOptions{Page: 0, Size: 100})
	if err != nil {
		t.Fatalf("Default page error: %v", err)
	}

	// Should get 80 repos (50 from ws1 + 30 from ws2, total < page size)
	if len(repos) != 80 {
		t.Errorf("Default page: expected 80 repos, got %d", len(repos))
	}

	// No next page since 80 < 100
	if res.Page.Next != 0 {
		t.Errorf("Default page: expected Next=0, got %d", res.Page.Next)
	}
}

// TestList_WorkspaceOrder tests that repos are returned in workspace order
func TestList_WorkspaceOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "aaa-first"}},
					{"workspace": map[string]string{"slug": "zzz-last"}},
				},
				"next": "",
			})
			return
		}

		if r.URL.Path == "/2.0/repositories/aaa-first" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("aaa-first", 0, 60, "", 60))
			return
		}

		if r.URL.Path == "/2.0/repositories/zzz-last" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("zzz-last", 0, 40, "", 40))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Page 1: Should get all 60 from aaa-first + 40 from zzz-last (fills to 100)
	repos1, _, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(repos1) != 100 {
		t.Errorf("Expected 100 repos, got %d", len(repos1))
	}

	// Verify order: first 60 from aaa-first, then 40 from zzz-last
	if repos1[0].Namespace != "aaa-first" {
		t.Errorf("Expected first repo from aaa-first, got %s", repos1[0].Namespace)
	}
	if repos1[59].Namespace != "aaa-first" {
		t.Errorf("Expected repo 59 from aaa-first, got %s", repos1[59].Namespace)
	}
	if repos1[60].Namespace != "zzz-last" {
		t.Errorf("Expected repo 60 from zzz-last, got %s", repos1[60].Namespace)
	}
	if repos1[99].Namespace != "zzz-last" {
		t.Errorf("Expected last repo from zzz-last, got %s", repos1[99].Namespace)
	}
}

// TestListV2_WithSearchTerm tests pagination with search filter
func TestListV2_WithSearchTerm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "workspace1"}},
				},
				"next": "",
			})
			return
		}

		if r.URL.Path == "/2.0/repositories/workspace1" {
			// Verify search query is present
			query := r.URL.Query().Get("q")
			if query != `name~"test"` {
				t.Errorf("Expected search query name~\"test\", got: %s", query)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("workspace1", 0, 25, "", 25))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	// Test with search term
	repos, _, err := client.Repositories.ListV2(ctx, scm.RepoListOptions{
		RepoSearchTerm: scm.RepoSearchTerm{RepoName: "test"},
		ListOptions:    scm.ListOptions{Page: 1, Size: 100},
	})

	if err != nil {
		t.Fatalf("ListV2 error: %v", err)
	}

	if len(repos) != 25 {
		t.Errorf("Expected 25 repos, got %d", len(repos))
	}
}

// Helper function to create mock repository page response
// workspace: workspace slug
// startIdx: starting index for repo numbering
// count: number of repos in this page
// nextPage: next page number (empty if last page)
// Optional: totalSize - total repos in workspace (for "size" field)
func createMockRepoPage(workspace string, startIdx, count int, nextPage string) map[string]interface{} {
	return createMockRepoPageWithSize(workspace, startIdx, count, nextPage, 0)
}

func createMockRepoPageWithSize(workspace string, startIdx, count int, nextPage string, totalSize int) map[string]interface{} {
	values := []map[string]interface{}{}

	for i := 0; i < count; i++ {
		repoNum := startIdx + i
		values = append(values, map[string]interface{}{
			"uuid":       fmt.Sprintf("{repo-%d-uuid}", repoNum),
			"full_name":  fmt.Sprintf("%s/repo-%d", workspace, repoNum),
			"is_private": false,
			"scm":        "git",
			"links": map[string]interface{}{
				"html":  map[string]string{"href": fmt.Sprintf("https://bitbucket.org/%s/repo-%d", workspace, repoNum)},
				"clone": []map[string]string{},
			},
			"mainbranch": map[string]string{"name": "main"},
		})
	}

	result := map[string]interface{}{
		"values": values,
		"next":   "",
	}

	// Add size field (total count) if provided
	if totalSize > 0 {
		result["size"] = totalSize
		result["pagelen"] = 100
		result["page"] = 1
	}

	if nextPage != "" {
		result["next"] = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s?page=%s&role=member", workspace, nextPage)
	}

	return result
}
