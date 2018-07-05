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

func TestRepositoryFind(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Repositories.Find(context.Background(), "PRJ/my-repo")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/repo.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryFind_NotFound(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/dev/repos/null").
		Reply(404).
		Type("application/json").
		File("testdata/error.json")

	client, _ := New("http://example.com:7990")
	_, _, err := client.Repositories.Find(context.Background(), "dev/null")
	if err == nil {
		t.Errorf("Expect not found message")
	}

	if got, want := err.Error(), "Project dev does not exist."; got != want {
		t.Errorf("Want error message %q, got %q", want, got)
	}
}

func TestRepositoryPerms(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo").
		Reply(200).
		Type("application/json").
		File("testdata/repo.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Repositories.FindPerms(context.Background(), "PRJ/my-repo")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Perm)
	raw, _ := ioutil.ReadFile("testdata/perms.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryList(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/2.0/repositories").
		MatchParam("page", "1").
		MatchParam("pagelen", "30").
		MatchParam("role", "member").
		Reply(200).
		Type("application/json").
		File("testdata/repos.json")

	client, _ := New("http://example.com:7990")
	got, res, err := client.Repositories.List(context.Background(), scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestStatusList(t *testing.T) {
	client, _ := New("http://example.com:7990")
	_, _, err := client.Repositories.ListStatus(context.Background(), "PRJ/my-repo", "a6e5e7d797edf751cbd839d6bd4aef86c941eec9", scm.ListOptions{Size: 30, Page: 1})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestStatusCreate(t *testing.T) {
	client, _ := New("http://example.com:7990")
	_, _, err := client.Repositories.CreateStatus(context.Background(), "PRJ/my-repo", "a6e5e7d797edf751cbd839d6bd4aef86c941eec9", &scm.StatusInput{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestRepositoryHookFind(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/webhooks/1").
		Reply(200).
		Type("application/json").
		File("testdata/webhook.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Repositories.FindHook(context.Background(), "PRJ/my-repo", "1")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/webhook.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryHookList(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/webhooks").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/webhooks.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Repositories.ListHooks(context.Background(), "PRJ/my-repo", scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Hook{}
	raw, _ := ioutil.ReadFile("testdata/webhooks.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepositoryHookDelete(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Delete("/rest/api/1.0/projects/PRJ/repos/my-repo/webhooks/1").
		Reply(200).
		Type("application/json")

	client, _ := New("http://example.com:7990")
	_, err := client.Repositories.DeleteHook(context.Background(), "PRJ/my-repo", "1")
	if err != nil {
		t.Error(err)
	}
}

func TestRepositoryHookCreate(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("/rest/api/1.0/projects/PRJ/repos/my-repo/webhooks").
		Reply(201).
		Type("application/json").
		File("testdata/webhook.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.Repositories.CreateHook(context.Background(), "PRJ/my-repo", &scm.HookInput{
		Name:   "example",
		Target: "http://example.com",
		Secret: "12345",
		Events: scm.HookEvents{
			Branch:             true,
			PullRequest:        true,
			PullRequestComment: true,
			Push:               true,
			Tag:                true,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Hook)
	raw, _ := ioutil.ReadFile("testdata/webhook.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestConvertPerms(t *testing.T) {
	tests := []struct {
		src *perm
		dst *scm.Perm
	}{
		{
			src: &perm{Permissions: "admin"},
			dst: &scm.Perm{Admin: true, Push: true, Pull: true},
		},
		{
			src: &perm{Permissions: "write"},
			dst: &scm.Perm{Admin: false, Push: true, Pull: true},
		},
		{
			src: &perm{Permissions: "read"},
			dst: &scm.Perm{Admin: false, Push: false, Pull: true},
		},
		{
			src: nil,
			dst: &scm.Perm{Admin: false, Push: false, Pull: false},
		},
	}
	for _, test := range tests {
		src := new(perms)
		if test.src != nil {
			src.Values = append(src.Values, test.src)
		}
		dst := convertPerms(src)
		if diff := cmp.Diff(test.dst, dst); diff != "" {
			t.Errorf("Unexpected Results")
			t.Log(diff)
		}
	}
}
