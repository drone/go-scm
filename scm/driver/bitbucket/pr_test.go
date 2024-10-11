// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

func TestPullFind(t *testing.T) {
	t.Skip()
}

func TestPullList(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/octocat/hello-world/pullrequests").
		MatchParam("pagelen", "30").
		MatchParam("page", "1").
		MatchParam("state", "all").
		Reply(200).
		Type("application/json").
		File("testdata/pulls.json")

	client := NewDefault()
	got, _, err := client.PullRequests.List(context.Background(), "octocat/hello-world", &scm.PullRequestListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.PullRequest{}
	raw, _ := os.ReadFile("testdata/pulls.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Logf("got JSON: %s", data)
	}
}

func TestPullListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/pullrequests/1/diffstat").
		MatchParam("pagelen", "30").
		MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_diffstat.json")

	client, _ := New("https://api.bitbucket.org")
	got, _, err := client.PullRequests.ListChanges(context.Background(), "atlassian/atlaskit", 1, &scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := os.ReadFile("testdata/pr_diffstat.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
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

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Post("2.0/repositories/octocat/hello-world/pullrequests").
		Reply(201).
		Type("application/json").
		File("testdata/pr_create.json")

	input := &scm.PullRequestInput{
		Title: "Amazing new feature",
		Body:  "Please pull these awesome changes in!",
		Head:  "octocat:new-feature",
		Base:  "master",
	}

	client, _ := New("https://api.bitbucket.org")

	got, _, err := client.PullRequests.Create(context.Background(), "octocat/hello-world", input)
	if err != nil {
		t.Fatal(err)
	}

	want := new(scm.PullRequest)
	raw, err := os.ReadFile("testdata/pr_create.json.golden")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
		data, err := json.Marshal(got)
		if err == nil {
			t.Logf("generated %s\n", string(data))
		} else {
			t.Errorf("failed to marshal got to JSON: %s\n", err.Error())
		}
	}
}
