// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
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
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		json.NewEncoder(os.Stdout).Encode(got)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestUserLoginFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
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
	raw, _ := ioutil.ReadFile("testdata/user_search.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
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
	if got, want := got, "john@example.com"; got != want {
		t.Errorf("Want user Email %q, got %q", want, got)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
