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

func TestRepositoryFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Repositories.Find(context.Background(), "atlassian/stash-example-plugin")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Repository", testRepository(result))
}

func TestRepositoryPerms(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Repositories.FindPerms(context.Background(), "atlassian/stash-example-plugin")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Permissions", testPermissions(result))
}

func TestRepositoryList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Repositories.List(context.Background(), scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d repositories, got %d", want, got)
		return
	}
	t.Run("Page", testPage(res))
	t.Run("Repository", testRepository(result[0]))
}

func TestStatusList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Repositories.ListStatus(context.Background(), "atlassian/stash-example-plugin", "a6e5e7d797edf751cbd839d6bd4aef86c941eec9", scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d statuses, got %d", want, got)
		return
	}
	t.Run("Page", testPage(res))
	t.Run("Status", testStatus(result[0]))
}

func TestStatusCreate(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	in := &scm.StatusInput{
		Desc:   "Build has completed successfully",
		Label:  "continuous-integration/drone",
		State:  scm.StateSuccess,
		Target: "https://ci.example.com/1000/output",
	}

	client, _ := New(server.URL)
	result, _, err := client.Repositories.CreateStatus(context.Background(), "atlassian/stash-example-plugin", "a6e5e7d797edf751cbd839d6bd4aef86c941eec9", in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Status", testStatus(result))
}

func TestRepositoryHookFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Repositories.FindHook(context.Background(), "atlassian/stash-example-plugin", "{d53603cc-3f67-45ea-b310-aaa5ef6ec061}")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Hook", testHook(result))
}

func TestRepositoryHookList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Repositories.ListHooks(context.Background(), "atlassian/stash-example-plugin", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d hooks, got %d", want, got)
		return
	}
	t.Run("Hook", testHook(result[0]))
}

func TestRepositoryHookDelete(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	_, err := client.Repositories.DeleteHook(context.Background(), "atlassian/stash-example-plugin", "{d53603cc-3f67-45ea-b310-aaa5ef6ec061}")
	if err != nil {
		t.Error(err)
	}
}

func TestRepositoryHookCreate(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Repositories.CreateHook(context.Background(), "atlassian/stash-example-plugin", &scm.HookInput{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Hook", testHook(result))
}

func testRepository(repository *scm.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := repository.ID, "{7dd600e6-0d9c-4801-b967-cb4cc17359ff}"; got != want {
			t.Errorf("Want repository ID %q, got %q", want, got)
		}
		if got, want := repository.Name, "stash-example-plugin"; got != want {
			t.Errorf("Want repository Name %q, got %q", want, got)
		}
		if got, want := repository.Namespace, "atlassian"; got != want {
			t.Errorf("Want repository Namespace %q, got %q", want, got)
		}
		if got, want := repository.Branch, "master"; got != want {
			t.Errorf("Want repository Branch %q, got %q", want, got)
		}
		if got, want := repository.Private, true; got != want {
			t.Errorf("Want repository Private %v, got %v", want, got)
		}
	}
}

func testPermissions(perms *scm.Perm) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := perms.Pull, true; got != want {
			t.Errorf("Want permission Pull %v, got %v", want, got)
		}
		if got, want := perms.Push, true; got != want {
			t.Errorf("Want permission Push %v, got %v", want, got)
		}
		if got, want := perms.Admin, true; got != want {
			t.Errorf("Want permission Admin %v, got %v", want, got)
		}
	}
}

func testHook(hook *scm.Hook) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := hook.ID, "{d53603cc-3f67-45ea-b310-aaa5ef6ec061}"; got != want {
			t.Errorf("Want hook ID %v, got %v", want, got)
		}
		if got, want := hook.Active, true; got != want {
			t.Errorf("Want hook Active %v, got %v", want, got)
		}
		if got, want := hook.Target, "http://example.com/webhook"; got != want {
			t.Errorf("Want hook Target %v, got %v", want, got)
		}
	}
}

func testStatus(status *scm.Status) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := status.Target, "https://ci.example.com/1000/output"; got != want {
			t.Errorf("Want status Target %v, got %v", want, got)
		}
		if got, want := status.State, scm.StateSuccess; got != want {
			t.Errorf("Want status State %v, got %v", want, got)
		}
		if got, want := status.Label, "drone"; got != want {
			t.Errorf("Want status Label %v, got %v", want, got)
		}
		if got, want := status.Desc, "Build has completed successfully"; got != want {
			t.Errorf("Want status Desc %v, got %v", want, got)
		}
	}
}
