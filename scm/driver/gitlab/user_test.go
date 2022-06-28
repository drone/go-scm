// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

func TestUserFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/user").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/user.json")

	client := NewDefault()
	got, res, err := client.Users.Find(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.User)
	raw, _ := os.ReadFile("testdata/user.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		err := json.NewEncoder(os.Stdout).Encode(got)
		if err != nil {
			t.Error(err)
		}
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestUserLoginFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/users").
		MatchParam("search", "john_smith").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/user_search.json")

	client := NewDefault()
	got, res, err := client.Users.FindLogin(context.Background(), "john_smith")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.User)
	raw, _ := os.ReadFile("testdata/user_search.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestUserLoginFind_NotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/users").
		MatchParam("search", "jcitizen").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/user_search.json")

	client := NewDefault()
	_, _, err := client.Users.FindLogin(context.Background(), "jcitizen")
	if err != scm.ErrNotFound {
		t.Errorf("Want Not Found Error, got %s", err)
	}
}

func TestUserLoginFind_NotAuthorized(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/users").
		MatchParam("search", "jcitizen").
		Reply(401).
		Type("application/json").
		SetHeaders(mockHeaders).
		BodyString(`{"message":"401 Unauthorized"}`)

	client := NewDefault()
	_, _, err := client.Users.FindLogin(context.Background(), "jcitizen")
	if err == nil {
		t.Errorf("Want 401 Unauthorized")
		return
	}
	if got, want := err.Error(), "Unauthorized"; got != want {
		t.Errorf("Want %s, got %s", want, got)
	}
}

func TestUserEmailFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/user").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/user.json")

	client := NewDefault()
	got, res, err := client.Users.FindEmail(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if want := "john@example.com"; got != want {
		t.Errorf("Want user Email %q, got %q", want, got)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
