// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

// TODO(bradrydzewski) repository html link is missing
// TODO(bradrydzewski) repository create date is missing
// TODO(bradrydzewski) repository update date is missing

func TestRepositoryFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/repo.json")

	client := NewDefault()
	got, res, err := client.Repositories.Find(context.Background(), "diaspora/diaspora")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/repo.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestSubgroupFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/gitlab-org/gitter/gitter-demo-app").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/subgroup.json")

	client := NewDefault()
	got, _, err := client.Repositories.Find(context.Background(), "gitlab-org/gitter/gitter-demo-app")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/subgroup.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Cannot parse subgroup")
		t.Log(diff)
	}
}

func TestRepositoryPerms(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/repo.json")

	client := NewDefault()
	perms, res, err := client.Repositories.FindPerms(context.Background(), "diaspora/diaspora")
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := perms.Pull, true; got != want {
		t.Errorf("Want permission Pull %v, got %v", want, got)
	}
	if got, want := perms.Push, false; got != want {
		t.Errorf("Want permission Push %v, got %v", want, got)
	}
	if got, want := perms.Admin, false; got != want {
		t.Errorf("Want permission Admin %v, got %v", want, got)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryNotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/dev/null").
		Reply(404).
		Type("application/json").
		SetHeaders(mockHeaders).
		BodyString(`{"message":"404 Project Not Found"}`)

	client := NewDefault()
	_, _, err := client.Repositories.Find(context.Background(), "dev/null")
	if err == nil {
		t.Errorf("Expect Not Found error")
		return
	}
	if got, want := err.Error(), "404 Project Not Found"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
}

func TestRepositoryList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("membership", "true").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/repos.json")

	client := NewDefault()
	got, res, err := client.Repositories.List(context.Background(), scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestStatusList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e/statuses").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/statuses.json")

	client := NewDefault()
	got, res, err := client.Repositories.ListStatus(context.Background(), "diaspora/diaspora", "6dcb09b5b57875f334f61aebed695e2e4193db5e", scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Status{}
	raw, _ := ioutil.ReadFile("testdata/statuses.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestStatusCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Post("/api/v4/projects/diaspora/diaspora/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e").
		MatchParam("name", "continuous-integration/drone").
		MatchParam("state", "success").
		MatchParam("target_url", "https://ci.example.com/diaspora/diaspora/42").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/status.json")

	in := &scm.StatusInput{
		Desc:   "Build has completed successfully",
		Label:  "continuous-integration/drone",
		State:  scm.StateSuccess,
		Target: "https://ci.example.com/diaspora/diaspora/42",
	}

	client := NewDefault()
	got, res, err := client.Repositories.CreateStatus(context.Background(), "diaspora/diaspora", "6dcb09b5b57875f334f61aebed695e2e4193db5e", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Status)
	raw, _ := ioutil.ReadFile("testdata/status.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/hooks/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	client := NewDefault()
	got, res, err := client.Repositories.FindHook(context.Background(), "diaspora/diaspora", "1")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/hooks").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/hooks.json")

	client := NewDefault()
	got, res, err := client.Repositories.ListHooks(context.Background(), "diaspora/diaspora", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Hook{}
	raw, _ := ioutil.ReadFile("testdata/hooks.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestRepositoryHookCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Post("/api/v4/projects/diaspora/diaspora/hooks").
		MatchParam("token", "topsecret").
		MatchParam("url", "https://ci.example.com/hook").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	in := &scm.HookInput{
		Name:       "drone",
		Target:     "https://ci.example.com/hook",
		Secret:     "topsecret",
		SkipVerify: false,
	}

	client := NewDefault()
	got, res, err := client.Repositories.CreateHook(context.Background(), "diaspora/diaspora", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookCreate_SkipVerification(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Post("/api/v4/projects/diaspora/diaspora/hooks").
		MatchParam("enable_ssl_verification", "false").
		MatchParam("token", "topsecret").
		MatchParam("url", "https://ci.example.com/hook").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook_skip_verification.json")

	in := &scm.HookInput{
		Name:       "drone",
		Target:     "https://ci.example.com/hook",
		Secret:     "topsecret",
		SkipVerify: true,
	}

	client := NewDefault()
	got, res, err := client.Repositories.CreateHook(context.Background(), "diaspora/diaspora", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook_skip_verification.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookUpdate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Put("/api/v4/projects/diaspora/diaspora/hooks/1").
		MatchParam("token", "topsecret").
		MatchParam("url", "https://ci.example.com/hook").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	in := &scm.HookInput{
		Name:       "drone",
		Target:     "https://ci.example.com/hook",
		Secret:     "topsecret",
		SkipVerify: false,
	}

	client := NewDefault()
	got, res, err := client.Repositories.UpdateHook(context.Background(), "diaspora/diaspora", "1", in)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/hook.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestRepositoryHookDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Delete("/api/v4/projects/diaspora/diaspora/hooks/1").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Repositories.DeleteHook(context.Background(), "diaspora/diaspora", "1")
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 204; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestConvertFromHookEvents(t *testing.T) {
	str := func(v scm.HookEvents) string {
		s := fmt.Sprintf("%+v", v)
		return regexp.MustCompile(" *[A-Za-z]+:false *").ReplaceAllString(s, "")
	}
	for _, test := range []struct {
		src scm.HookEvents
		dst string
	}{{
		src: scm.HookEvents{Push: true},
		dst: "push_events",
	}, {
		src: scm.HookEvents{Branch: true},
		dst: "push_events",
	}, {
		src: scm.HookEvents{Issue: true},
		dst: "issues_events",
	}, {
		src: scm.HookEvents{PullRequest: true},
		dst: "merge_requests_events",
	}, {
		src: scm.HookEvents{Tag: true},
		dst: "tag_push_events",
	}, {
		src: scm.HookEvents{IssueComment: true},
		dst: "note_events",
	}, {
		src: scm.HookEvents{PullRequestComment: true},
		dst: "note_events",
	}, {
		src: scm.HookEvents{ReviewComment: true},
		dst: "note_events",
	}, {
		src: scm.HookEvents{Deployment: true},
		dst: "deployment_events",
	}} {
		events := convertFromHookEvents(test.src)
		for _, k := range []string{
			"push_events",
			"issues_events",
			"merge_requests_events",
			"confidential_issues_events",
			"tag_push_events",
			"note_events",
			"confidential_note_events",
			"job_events",
			"pipeline_events",
			"wiki_page_events",
			"deployment_events",
		} {
			if v, ok := events[k]; !ok || v != (k == test.dst) {
				got, want := fmt.Sprintf("%s=%t", k, v), fmt.Sprintf("%s=%t", k, k == test.dst)
				if !ok {
					got = "none"
				}
				t.Errorf("Want %s converted to %s, got %s", str(test.src), want, got)
			}
		}
	}
}

func TestConvertState(t *testing.T) {
	tests := []struct {
		src string
		dst scm.State
	}{
		{
			src: "failed",
			dst: scm.StateFailure,
		},
		{
			src: "canceled",
			dst: scm.StateCanceled,
		},
		{
			src: "pending",
			dst: scm.StatePending,
		},
		{
			src: "running",
			dst: scm.StateRunning,
		},
		{
			src: "success",
			dst: scm.StateSuccess,
		},
		{
			src: "invalid",
			dst: scm.StateUnknown,
		},
	}
	for i, test := range tests {
		if got, want := convertState(test.src), test.dst; got != want {
			t.Errorf("Want state %s converted to %v at index %d", test.src, test.dst, i)
		}
	}
}

func TestConvertFromState(t *testing.T) {
	tests := []struct {
		src scm.State
		dst string
	}{
		{
			src: scm.StateCanceled,
			dst: "canceled",
		},
		{
			src: scm.StateError,
			dst: "failed",
		},
		{
			src: scm.StateFailure,
			dst: "failed",
		},
		{
			src: scm.StatePending,
			dst: "pending",
		},
		{
			src: scm.StateRunning,
			dst: "running",
		},
		{
			src: scm.StateSuccess,
			dst: "success",
		},
		{
			src: scm.StateUnknown,
			dst: "failed",
		},
	}
	for i, test := range tests {
		if got, want := convertFromState(test.src), test.dst; got != want {
			t.Errorf("Want state %v converted to %s at index %d", test.src, test.dst, i)
		}
	}
}

func TestConvertPrivate(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"public", false},
		{"", false},
		{"private", true},
		{"internal", true},
		{"invalid", true},
	}

	for _, test := range tests {
		if got, want := convertPrivate(test.in), test.out; got != want {
			t.Errorf("Want private %v, got %v", want, got)
		}
	}
}

func TestCanPush(t *testing.T) {
	tests := []struct {
		in  *repository
		out bool
	}{
		//
		// access granted
		//
		{
			out: true,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 30},
					GroupAccess:   access{AccessLevel: 0},
				},
			},
		},
		{
			out: true,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 31},
					GroupAccess:   access{AccessLevel: 0},
				},
			},
		},
		{
			out: true,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 0},
					GroupAccess:   access{AccessLevel: 30},
				},
			},
		},
		{
			out: true,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 0},
					GroupAccess:   access{AccessLevel: 31},
				},
			},
		},
		//
		// access denied
		//
		{
			out: false,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 29},
					GroupAccess:   access{AccessLevel: 0},
				},
			},
		},
		{
			out: false,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 0},
					GroupAccess:   access{AccessLevel: 29},
				},
			},
		},
	}

	for _, test := range tests {
		if got, want := canPush(test.in), test.out; got != want {
			t.Errorf("Want push authorization %v, got %v", want, got)
		}
	}
}

func TestCanAdmin(t *testing.T) {
	tests := []struct {
		in  *repository
		out bool
	}{
		//
		// access granted
		//
		{
			out: true,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 40},
					GroupAccess:   access{AccessLevel: 0},
				},
			},
		},
		{
			out: true,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 41},
					GroupAccess:   access{AccessLevel: 0},
				},
			},
		},
		{
			out: true,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 0},
					GroupAccess:   access{AccessLevel: 40},
				},
			},
		},
		{
			out: true,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 0},
					GroupAccess:   access{AccessLevel: 41},
				},
			},
		},
		//
		// access denied
		//
		{
			out: false,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 39},
					GroupAccess:   access{AccessLevel: 0},
				},
			},
		},
		{
			out: false,
			in: &repository{
				Permissions: permissions{
					ProjectAccess: access{AccessLevel: 0},
					GroupAccess:   access{AccessLevel: 39},
				},
			},
		},
	}

	for _, test := range tests {
		if got, want := canAdmin(test.in), test.out; got != want {
			t.Errorf("Want admin authorization %v, got %v", want, got)
		}
	}
}
