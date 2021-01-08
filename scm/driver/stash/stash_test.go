// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"github.com/jenkins-x/go-scm/scm"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := New("https://api.bitbucket.org")
	if err != nil {
		t.Error(err)
	}
	if got, want := client.BaseURL.String(), "https://api.bitbucket.org/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Base(t *testing.T) {
	client, err := New("https://api.bitbucket.org/v1")
	if err != nil {
		t.Error(err)
	}
	if got, want := client.BaseURL.String(), "https://api.bitbucket.org/v1/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Error(t *testing.T) {
	_, err := New("http://a b.com/")
	if err == nil {
		t.Errorf("Expect error when invalid URL")
	}
}

var mockHeaders = map[string]string{
	"X-AREQUESTID": "@1XQQXPVx888x374364x0",
}

func testRequest(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Header.Get("X-AREQUESTID"), "@1XQQXPVx888x374364x0"; got != want {
			t.Errorf("Want X-AREQUESTID %q, got %q", want, got)
		}
	}
}
