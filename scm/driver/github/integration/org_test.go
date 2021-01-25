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
// organization sub-tests
//

func testOrgs(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Find", testOrgFind(client))
	}
}

func testOrgFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		result, _, err := client.Organizations.Find(context.Background(), "github")
		if err != nil {
			t.Error(err)
			return
		}
		t.Run("Organization", testOrg(result))
	}
}

//
// struct sub-tests
//

func testOrg(organization *scm.Organization) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		if got, want := organization.Name, "github"; got != want {
			t.Errorf("Want organization Name %q, got %q", want, got)
		}
		orgAvatar := `https://avatars[0-9]?\.githubusercontent\.com/u/9919\?v=4`
		if matched, err := regexp.MatchString(orgAvatar, organization.Avatar); !matched {
			t.Errorf("Want organization Avatar %q, got %q", orgAvatar, organization.Avatar)
		} else if err != nil {
			t.Errorf("invalid regexp %q: %v", orgAvatar, err)
		}
	}
}
