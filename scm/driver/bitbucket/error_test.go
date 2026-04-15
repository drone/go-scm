// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

