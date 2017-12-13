// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github/fixtures"
)

func TestUserFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Users.Find(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testUser(result))
}

func TestUserLoginFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Users.FindLogin(context.Background(), "octocat")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Fields", testUser(result))
}

func TestUserEmailFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Users.FindEmail(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := result, "octocat@github.com"; got != want {
		t.Errorf("Want user Email %q, got %q", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func testUser(user *scm.User) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := user.Login, "octocat"; got != want {
			t.Errorf("Want user Login %v, got %v", want, got)
		}
		if got, want := user.Email, "octocat@github.com"; got != want {
			t.Errorf("Want user Email %v, got %v", want, got)
		}
		if got, want := user.Avatar, "https://github.com/images/error/octocat_happy.gif"; got != want {
			t.Errorf("Want user Avatar %v, got %v", want, got)
		}
		if got, want := user.Name, "monalisa octocat"; got != want {
			t.Errorf("Want user Name %v, got %v", want, got)
		}
	}
}
