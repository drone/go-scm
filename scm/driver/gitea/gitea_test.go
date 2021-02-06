// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package gitea implements a Gogs client.
package gitea

import (
	"testing"

	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

var mockPageHeaders = map[string]string{
	"Link": `<https://try.gitea.io/v1/resource?page=2>; rel="next",` +
		`<https://try.gitea.io/v1/resource?page=1>; rel="prev",` +
		`<https://try.gitea.io/v1/resource?page=1>; rel="first",` +
		`<https://try.gitea.io/v1/resource?page=5>; rel="last"`,
}

func TestClient(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	client, err := New("https://try.gitea.io")
	if err != nil {
		t.Error(err)
	}
	if got, want := client.BaseURL.String(), "https://try.gitea.io/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Base(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	client, err := New("https://try.gitea.io/v1")
	if err != nil {
		t.Error(err)
	}
	if got, want := client.BaseURL.String(), "https://try.gitea.io/v1/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Error(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	_, err := New("http://a b.com/")
	if err == nil {
		t.Errorf("Expect error when invalid URL")
	}
}

func testPage(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Page.Next, 2; got != want {
			t.Errorf("Want next page %d, got %d", want, got)
		}
		if got, want := res.Page.Prev, 1; got != want {
			t.Errorf("Want prev page %d, got %d", want, got)
		}
		if got, want := res.Page.First, 1; got != want {
			t.Errorf("Want first page %d, got %d", want, got)
		}
		if got, want := res.Page.Last, 5; got != want {
			t.Errorf("Want last page %d, got %d", want, got)
		}
	}
}

func mockServerVersion() {
	gock.New("https://try.gitea.io").
		Get("/api/v1/version").
		Reply(200).
		Type("application/json").
		File("testdata/version.json")
}
