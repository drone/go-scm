// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"context"
	"regexp"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
)

//
// user sub-tests
//

func testUsers(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testUserFind(client))
	}
}

func testUserFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Users.FindLogin(context.Background(), "octocat")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("User", testUser(result))
	}
}

//
// struct sub-tests
//

func testUser(user *scm.User) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		if got, want := user.Login, "octocat"; got != want {
			t.Errorf("Want user Login %q, got %q", want, got)
		}
		if got, want := user.Name, "The Octocat"; got != want {
			t.Errorf("Want user Name %q, got %q", want, got)
		}
		userAvatar := `https://avatars[0-9]?\.githubusercontent\.com/u/583231\?v=4`
		if matched, err := regexp.MatchString(userAvatar, user.Avatar); !matched {
			t.Errorf("Want user Avatar %q, got %q", userAvatar, user.Avatar)
		} else if err != nil {
			t.Errorf("invalid regexp %q: %v", userAvatar, err)
		}
	}
}
