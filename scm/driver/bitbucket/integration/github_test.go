// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"net/http"
	"os"
	"testing"

	"github.com/drone/go-scm/scm/driver/bitbucket"
	"github.com/drone/go-scm/scm/token"
)

func TestGitLab(t *testing.T) {
	if os.Getenv("BITBUCKET_TOKEN") == "" {
		t.Skipf("missing BITBUCKET_TOKEN environment variable")
		return
	}

	client := bitbucket.NewDefault()
	client.Client = &http.Client{
		Transport: &token.Transport{
			SetToken: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+os.Getenv("BITBUCKET_TOKEN"))
			},
		},
	}

	t.Run("Contents", testContents(client))
	t.Run("Git", testGit(client))
	t.Run("Issues", testIssues(client))
	t.Run("Organizations", testOrgs(client))
	t.Run("PullRequests", testPullRequests(client))
	t.Run("Repositories", testRepos(client))
	t.Run("Users", testUsers(client))
}
