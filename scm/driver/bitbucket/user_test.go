// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/bitbucket/fixtures"
)

func TestUserFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Users.Find(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Fields", testUser(result))
}

func TestUserLoginFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Users.FindLogin(context.Background(), "brydzewski")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Fields", testUser(result))
}

func testUser(user *scm.User) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := user.Login, "brydzewski"; got != want {
			t.Errorf("Want user Login %v, got %v", want, got)
		}
		if got, want := user.Avatar, "https://bitbucket.org/account/brydzewski/avatar/32/"; got != want {
			t.Errorf("Want user Avatar %v, got %v", want, got)
		}
		if got, want := user.Name, "Brad Rydzewski"; got != want {
			t.Errorf("Want user Name %v, got %v", want, got)
		}
	}
}
