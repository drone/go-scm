// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

//
// repository sub-tests
//

func TestRepoFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Repositories.Find(context.Background(), "go-gitea/gitea")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Repository)
	raw, _ := os.ReadFile("testdata/repo.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepoFindPerm(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Repositories.FindPerms(context.Background(), "go-gitea/gitea")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Repository)
	raw, _ := os.ReadFile("testdata/repo.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want.Perm); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepoList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/user/repos").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/repos.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Repositories.List(context.Background(), &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Repository{}
	raw, _ := os.ReadFile("testdata/repos.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestRepoNotFound(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/gogits/go-gogs-client").
		Reply(404).
		Type("text/plain")

	client, _ := New("https://demo.gitea.com")
	_, _, err := client.Repositories.FindPerms(context.Background(), "gogits/go-gogs-client")
	if err == nil {
		t.Errorf("Expect Not Found error")
	} else if got, want := err.Error(), "404 Not Found"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
}

//
// hook sub-tests
//

func TestHookFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/hooks/20").
		Reply(200).
		Type("application/json").
		File("testdata/hook.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Repositories.FindHook(context.Background(), "go-gitea/gitea", "20")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Hook)
	raw, _ := os.ReadFile("testdata/hook.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestHookList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/go-gitea/gitea/hooks").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/hooks.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Repositories.ListHooks(context.Background(), "go-gitea/gitea", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Hook{}
	raw, _ := os.ReadFile("testdata/hooks.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestHookCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Post("/api/v1/repos/go-gitea/gitea/hooks").
		Reply(201).
		Type("application/json").
		File("testdata/hook.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Repositories.CreateHook(context.Background(), "go-gitea/gitea", &scm.HookInput{})
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Hook)
	raw, _ := os.ReadFile("testdata/hook.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestHookDelete(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Delete("/api/v1/repos/go-gitea/gitea/hooks/20").
		Reply(204).
		Type("application/json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Repositories.DeleteHook(context.Background(), "go-gitea/gitea", "20")
	if err != nil {
		t.Error(err)
	}
}

func TestHookEvents(t *testing.T) {
	tests := []struct {
		in  scm.HookEvents
		out []string
	}{
		{
			in:  scm.HookEvents{Push: true},
			out: []string{"push"},
		},
		{
			in:  scm.HookEvents{Branch: true},
			out: []string{"create", "delete"},
		},
		{
			in:  scm.HookEvents{IssueComment: true},
			out: []string{"issue_comment"},
		},
		{
			in:  scm.HookEvents{PullRequestComment: true},
			out: []string{"issue_comment"},
		},
		{
			in:  scm.HookEvents{Issue: true},
			out: []string{"issues"},
		},
		{
			in:  scm.HookEvents{PullRequest: true},
			out: []string{"pull_request"},
		},
		{
			in:  scm.HookEvents{Review: true},
			out: []string{"pull_request_review"},
		},
		{
			in:  scm.HookEvents{Release: true},
			out: []string{"release"},
		},
		{
			in: scm.HookEvents{
				Branch:             true,
				Issue:              true,
				IssueComment:       true,
				PullRequest:        true,
				PullRequestComment: true,
				Push:               true,
				Release:            true,
				Review:             true,
				ReviewComment:      true,
				Tag:                true,
			},
			out: []string{"pull_request", "pull_request_review", "pull_request_review_comment", "issues", "issue_comment", "create", "delete", "push", "release"},
		},
	}
	for _, test := range tests {
		got, want := convertHookEvent(test.in), test.out
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Unexpected Results")
			t.Log(diff)
		}
	}
}

//
// status sub-tests
//

func TestStatusList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/jcitizen/my-repo/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e/statuses").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/statuses.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Repositories.ListStatus(context.Background(), "jcitizen/my-repo", "6dcb09b5b57875f334f61aebed695e2e4193db5e", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Status{}
	raw, _ := os.ReadFile("testdata/statuses.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestStatusCreate(t *testing.T) {
	in := &scm.StatusInput{
		Desc:   "Build has completed successfully",
		Label:  "continuous-integration/drone",
		State:  scm.StateSuccess,
		Target: "https://example.com/jcitizen/my-repo/1000",
	}

	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Post("/api/v1/repos/jcitizen/my-repo/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e").
		Reply(201).
		Type("application/json").
		File("testdata/status.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Repositories.CreateStatus(context.Background(), "jcitizen/my-repo", "6dcb09b5b57875f334f61aebed695e2e4193db5e", in)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Status)
	raw, _ := os.ReadFile("testdata/status.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
