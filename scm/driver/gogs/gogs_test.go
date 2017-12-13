// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gogs implements a Gogs client.
package gogs

import (
	"testing"

	"github.com/drone/go-scm/scm/driver/gogs/fixtures"
)

func TestClient(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, err := New(server.URL)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Contents", testContents(client))
	t.Run("Git", testGit(client))
	t.Run("Issues", testIssues(client))
	t.Run("Organizations", testOrgs(client))
	t.Run("PullRequests", testPullRequests(client))
	t.Run("Repositories", testRepos(client))
	t.Run("Reviews", testReviews(client))
	t.Run("Users", testUsers(client))
}
