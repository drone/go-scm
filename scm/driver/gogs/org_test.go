// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testOrgs(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testOrgFind(client))
		t.Run("List", testOrgList(client))
	}
}

func testOrgFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Organizations.Find(context.Background(), "gogits")
		if err != nil {
			t.Error(err)
		}
		t.Run("Organization", testOrganization(result))
	}
}

func testOrgList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Organizations.List(
			context.Background(),
			scm.ListOptions{},
		)
		if err != nil {
			t.Error(err)
		}
		if got, want := len(result), 1; got != want {
			t.Errorf("Want %d organizations, got %d", want, got)
		} else {
			t.Run("Organization", testOrganization(result[0]))
		}
	}
}

func testOrganization(organization *scm.Organization) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := organization.Name, "gogits"; got != want {
			t.Errorf("Want organization Name %q, got %q", want, got)
		}
		if got, want := organization.Avatar, "http://gogits.io/avatars/1"; got != want {
			t.Errorf("Want organization Avatar %q, got %q", want, got)
		}
	}
}
