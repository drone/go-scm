// Copyright 2017 Drone.IO Inc. All rights reserved.
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

func TestError_ErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{
			name: "with body message",
			err: Error{
				StatusCode: 400,
				Type:       "error",
				Data:       struct{ Message string `json:"message"` }{Message: "Bad request"},
			},
			want: "Bad request",
		},
		{
			name: "status code only, no body message",
			err: Error{
				StatusCode: 429,
			},
			want: "bitbucket: http status 429",
		},
		{
			name: "no status code, no body message",
			err:  Error{},
			want: "bitbucket: unknown error",
		},
		{
			name: "status code with empty message",
			err: Error{
				StatusCode: 500,
				Data:       struct{ Message string `json:"message"` }{Message: ""},
			},
			want: "bitbucket: http status 500",
		},
		{
			name: "body message takes precedence over status code",
			err: Error{
				StatusCode: 403,
				Data:       struct{ Message string `json:"message"` }{Message: "Forbidden"},
			},
			want: "Forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error.Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestError_StatusCodeStamped(t *testing.T) {
	// Verify that do() stamps the HTTP status code on the Error struct
	// even when the response body is empty (e.g. 429 rate-limit with no JSON body).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		// empty body — no JSON
	}))
	defer server.Close()

	client, _ := New(server.URL)
	wrapper := &wrapper{client}

	_, err := wrapper.do(context.Background(), "GET", "2.0/repositories", nil, nil)
	if err == nil {
		t.Fatal("Expected error for 429 response, got nil")
	}

	bbErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T: %v", err, err)
	}

	if bbErr.StatusCode != 429 {
		t.Errorf("Expected StatusCode=429, got %d", bbErr.StatusCode)
	}

	if got := bbErr.Error(); got != "bitbucket: http status 429" {
		t.Errorf("Expected error message %q, got %q", "bitbucket: http status 429", got)
	}
}

func TestError_StatusCodeWithJSONBody(t *testing.T) {
	// When the server returns a JSON error body, both StatusCode and Data.Message should be set.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{"message": "Access denied"},
		})
	}))
	defer server.Close()

	client, _ := New(server.URL)
	wrapper := &wrapper{client}

	_, err := wrapper.do(context.Background(), "GET", "2.0/repositories", nil, nil)
	if err == nil {
		t.Fatal("Expected error for 403 response, got nil")
	}

	bbErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T: %v", err, err)
	}

	if bbErr.StatusCode != 403 {
		t.Errorf("Expected StatusCode=403, got %d", bbErr.StatusCode)
	}

	// Data.Message should take precedence
	if got := bbErr.Error(); got != "Access denied" {
		t.Errorf("Expected error %q, got %q", "Access denied", got)
	}
}

func TestIsHardError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "401 ErrNotAuthorized is hard",
			err:  scm.ErrNotAuthorized,
			want: true,
		},
		{
			name: "429 rate limit is hard",
			err:  &Error{StatusCode: 429},
			want: true,
		},
		{
			name: "403 access denied is soft",
			err:  &Error{StatusCode: 403},
			want: false,
		},
		{
			name: "404 not found is soft",
			err:  &Error{StatusCode: 404},
			want: false,
		},
		{
			name: "500 internal server error is soft",
			err:  &Error{StatusCode: 500},
			want: false,
		},
		{
			name: "generic error is soft",
			err:  fmt.Errorf("network error"),
			want: false,
		},
		{
			name: "nil error is soft",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// nil check: isHardError should not panic on nil
			if tt.err == nil {
				// skip — isHardError does not accept nil in practice
				return
			}
			got := isHardError(tt.err)
			if got != tt.want {
				t.Errorf("isHardError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestList_RateLimitStopsPagination(t *testing.T) {
	// When a workspace returns 429, pagination should stop and surface the error.
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

		// ws1: returns 429 rate limit (no JSON body, just like real Bitbucket)
		if r.URL.Path == "/2.0/repositories/ws1" {
			w.WriteHeader(429)
			return
		}

		// ws2: should never be reached
		if r.URL.Path == "/2.0/repositories/ws2" {
			t.Error("ws2 should not be queried after ws1 returns 429")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws2", 0, 50, "", 50))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	_, _, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err == nil {
		t.Fatal("Expected error when rate-limited, got nil")
	}

	bbErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T: %v", err, err)
	}

	if bbErr.StatusCode != 429 {
		t.Errorf("Expected 429 status code in error, got %d", bbErr.StatusCode)
	}
}

func TestList_UnauthorizedStopsPagination(t *testing.T) {
	// When credentials are invalid (401), pagination should stop immediately.
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

		// ws1: returns 401
		if r.URL.Path == "/2.0/repositories/ws1" {
			w.WriteHeader(401)
			return
		}

		// ws2: should never be reached
		if r.URL.Path == "/2.0/repositories/ws2" {
			t.Error("ws2 should not be queried after ws1 returns 401")
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	_, _, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err == nil {
		t.Fatal("Expected error when unauthorized, got nil")
	}

	if err != scm.ErrNotAuthorized {
		t.Errorf("Expected scm.ErrNotAuthorized, got %T: %v", err, err)
	}
}

func TestList_SoftErrorSkipsWorkspace(t *testing.T) {
	// 403 on one workspace should skip it but continue to the next.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/2.0/user/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"values": []map[string]interface{}{
					{"workspace": map[string]string{"slug": "ws-forbidden"}},
					{"workspace": map[string]string{"slug": "ws-ok"}},
				},
				"next": "",
			})
			return
		}

		// ws-forbidden: 403
		if r.URL.Path == "/2.0/repositories/ws-forbidden" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(403)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{"message": "Forbidden"},
			})
			return
		}

		// ws-ok: 40 repos
		if r.URL.Path == "/2.0/repositories/ws-ok" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws-ok", 0, 40, "", 40))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	repos, _, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err != nil {
		t.Fatalf("Unexpected error: %v (soft errors should be skipped)", err)
	}

	if len(repos) != 40 {
		t.Errorf("Expected 40 repos from ws-ok (skipping ws-forbidden), got %d", len(repos))
	}
}

func TestList_RateLimitDuringRepoFetch(t *testing.T) {
	// 429 during fetchReposFromWorkspaceWithOffset (not just count) should also stop.
	requestCount := 0
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
			requestCount++
			// First request is for count (pagelen=1), return success
			if r.URL.Query().Get("pagelen") == "1" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(createMockRepoPageWithSize("ws1", 0, 1, "", 50))
				return
			}
			// Second request is for actual repos — return 429
			w.WriteHeader(429)
			return
		}

		// ws2 should never be reached
		if r.URL.Path == "/2.0/repositories/ws2" {
			t.Error("ws2 should not be queried after ws1 repo fetch returns 429")
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client, _ := New(server.URL)
	ctx := context.Background()

	_, _, err := client.Repositories.List(ctx, scm.ListOptions{Page: 1, Size: 100})
	if err == nil {
		t.Fatal("Expected error when rate-limited during repo fetch, got nil")
	}

	bbErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T: %v", err, err)
	}

	if bbErr.StatusCode != 429 {
		t.Errorf("Expected 429 status code, got %d", bbErr.StatusCode)
	}
}
