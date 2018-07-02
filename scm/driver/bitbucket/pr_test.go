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

func TestPullFind(t *testing.T) {
	t.Skip()
}

func TestPullList(t *testing.T) {
	t.Skip()
}

func TestPullListChanges(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.PullRequests.ListChanges(context.Background(), "atlassian/atlaskit", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if len(result) == 0 {
		t.Errorf("Want non-empty diff")
		return
	}
	if got, want := result[0].Path, "CONTRIBUTING.md"; got != want {
		t.Errorf("Want file path %q, got %q", want, got)
	}
}

func TestPullMerge(t *testing.T) {
	t.Skip()
}

func TestPullClose(t *testing.T) {
	t.Skip()
}
