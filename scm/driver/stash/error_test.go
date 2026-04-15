// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestError_WithErrorsSlice(t *testing.T) {
	e := &Error{
		Status: 400,
		Errors: []struct {
			Message         string `json:"message"`
			ExceptionName   string `json:"exceptionName"`
			CurrentVersion  int    `json:"currentVersion"`
			ExpectedVersion int    `json:"expectedVersion"`
		}{
			{Message: "field is required"},
		},
	}

	if got := e.Error(); got != "field is required" {
		t.Errorf("Error() = %q, want %q", got, "field is required")
	}
}

func TestError_WithMessage(t *testing.T) {
	e := &Error{
		Status:  404,
		Message: "Repository not found",
	}

	want := "bitbucket server: status: 404 message: Repository not found"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestError_StatusCodeOnly(t *testing.T) {
	e := &Error{
		Status: 429,
	}

	want := "bitbucket server: http status 429"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestError_UnknownError(t *testing.T) {
	e := &Error{}

	want := "bitbucket server: unknown error"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestError_StatusStampedOnEmptyBody(t *testing.T) {
	// Verify that do() stamps the HTTP status on Error.Status
	// when the response body is empty (e.g. 429 rate-limit).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		// empty body
	}))
	defer server.Close()

	client, _ := New(server.URL)
	wrapper := &wrapper{client}

	_, err := wrapper.do(context.Background(), "GET", "rest/api/1.0/repos", nil, nil)
	if err == nil {
		t.Fatal("Expected error for 429 response, got nil")
	}

	stashErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T: %v", err, err)
	}

	if stashErr.Status != 429 {
		t.Errorf("Expected Status=429, got %d", stashErr.Status)
	}

	want := "bitbucket server: http status 429"
	if got := stashErr.Error(); got != want {
		t.Errorf("Expected error %q, got %q", want, got)
	}
}

func TestError_StatusFromJSONPreserved(t *testing.T) {
	// When the JSON body includes a status-code field, it should NOT be overwritten.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status-code": 400,
			"message":     "Bad request from API",
		})
	}))
	defer server.Close()

	client, _ := New(server.URL)
	wrapper := &wrapper{client}

	_, err := wrapper.do(context.Background(), "GET", "rest/api/1.0/repos", nil, nil)
	if err == nil {
		t.Fatal("Expected error for 400 response, got nil")
	}

	stashErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T: %v", err, err)
	}

	// Status from JSON body should be preserved (not overwritten by HTTP status)
	if stashErr.Status != 400 {
		t.Errorf("Expected Status=400 from JSON, got %d", stashErr.Status)
	}

	want := "bitbucket server: status: 400 message: Bad request from API"
	if got := stashErr.Error(); got != want {
		t.Errorf("Expected error %q, got %q", want, got)
	}
}

func TestError_StatusFallbackWhenJSONHasNoStatus(t *testing.T) {
	// When JSON body has message but no status-code, HTTP status should be used as fallback.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(503)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Service unavailable",
		})
	}))
	defer server.Close()

	client, _ := New(server.URL)
	wrapper := &wrapper{client}

	_, err := wrapper.do(context.Background(), "GET", "rest/api/1.0/repos", nil, nil)
	if err == nil {
		t.Fatal("Expected error for 503 response, got nil")
	}

	stashErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T: %v", err, err)
	}

	// HTTP status should be stamped as fallback since JSON has no status-code
	if stashErr.Status != 503 {
		t.Errorf("Expected Status=503 (fallback from HTTP), got %d", stashErr.Status)
	}

	want := "bitbucket server: status: 503 message: Service unavailable"
	if got := stashErr.Error(); got != want {
		t.Errorf("Expected error %q, got %q", want, got)
	}
}
