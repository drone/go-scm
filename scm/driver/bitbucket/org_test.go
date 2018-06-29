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

func TestOrganizationFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Organizations.Find(context.Background(), "atlassian")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Fields", testOrganization(result))
}

func TestOrganizationList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Organizations.List(context.Background(), scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d organizations, got %d", want, got)
		return
	}
	t.Run("Page", testPage(res))
	t.Run("Fields", testOrganization(result[0]))
}

func testOrganization(organization *scm.Organization) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := organization.Name, "atlassian"; got != want {
			t.Errorf("Want organization Name %q, got %q", want, got)
		}
		if got, want := organization.Avatar, "https://bitbucket.org/account/atlassian/avatar/32/"; got != want {
			t.Errorf("Want organization Avatar %q, got %q", want, got)
		}
	}
}
