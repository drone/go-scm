// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullFind(t *testing.T) {
	t.Skip()
}

func TestPullList(t *testing.T) {
	t.Skip()
}

func TestPullListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/changes").
		// MatchParam("pagelen", "30").
		// MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_change.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.PullRequests.ListChanges(context.Background(), "PRJ/my-repo", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/pr-change.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullMerge(t *testing.T) {
	t.Skip()
}

func TestPullClose(t *testing.T) {
	client, _ := New("https://api.bitbucket.org")
	_, err := client.PullRequests.Close(context.Background(), "atlassian/atlaskit", 1)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
