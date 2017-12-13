// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitlab/fixtures"
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
	result, res, err := client.Users.FindLogin(context.Background(), "john_smith")
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
	if got, want := result, "john@example.com"; got != want {
		t.Errorf("Want user Email %q, got %q", want, got)
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func testUser(user *scm.User) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := user.Login, "john_smith"; got != want {
			t.Errorf("Want user Login %v, got %v", want, got)
		}
		if got, want := user.Email, "john@example.com"; got != want {
			t.Errorf("Want user Email %v, got %v", want, got)
		}
		if got, want := user.Avatar, "http://localhost:3000/uploads/user/avatar/1/index.jpg"; got != want {
			t.Errorf("Want user Avatar %v, got %v", want, got)
		}
		if got, want := user.Name, "John Smith"; got != want {
			t.Errorf("Want user Name %v, got %v", want, got)
		}
	}
}
