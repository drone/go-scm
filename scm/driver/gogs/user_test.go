// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testUsers(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testUserFind(client))
		t.Run("FindLogin", testUserFindLogin(client))
		t.Run("FindEmail", testUserFindEmail(client))
	}
}

func testUserFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Users.Find(context.Background())
		if err != nil {
			t.Error(err)
		}
		if got, want := result.Login, "janedoe"; got != want {
			t.Errorf("Want user Login %q, got %q", want, got)
		}
	}
}

func testUserFindLogin(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Users.FindLogin(context.Background(), "janedoe")
		if err != nil {
			t.Error(err)
		}
		if got, want := result.Login, "janedoe"; got != want {
			t.Errorf("Want user Login %q, got %q", want, got)
		}
	}
}

func testUserFindEmail(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Users.FindEmail(context.Background())
		if err != nil {
			t.Error(err)
		}
		if got, want := result, "janedoe@gmail.com"; got != want {
			t.Errorf("Want user Email %q, got %q", want, got)
		}
	}
}
