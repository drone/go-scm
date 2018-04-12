// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitlab/fixtures"
)

func TestRepositoryFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Repositories.Find(context.Background(), "diaspora/diaspora")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Repository", testRepository(result))
	t.Run("Permissions", testPermissions(result.Perm))
}

func TestRepositoryPerms(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Repositories.FindPerms(context.Background(), "diaspora/diaspora")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Permissions", testPermissions(result))
}

func TestRepositoryNotFound(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	_, _, err := client.Repositories.FindPerms(context.Background(), "not/found")
	if err == nil {
		t.Errorf("Expect Not Found error")
		return
	}
	if got, want := err.Error(), "404 Project Not Found"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
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
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Repository", testRepository(result[0]))
	t.Run("Permissions", testPermissions(result[0].Perm))
}

func TestStatusList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Repositories.ListStatus(context.Background(), "diaspora/diaspora", "18f3e63d05582537db6d183d9d557be09e1f90c8", scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d statuses, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Status", testStatus(result[0]))
}

func TestStatusCreate(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	in := &scm.StatusInput{
		Desc:   "Build has completed successfully",
		Label:  "continuous-integration/jenkins",
		State:  scm.StateSuccess,
		Target: "https://ci.example.com/1000/output",
	}

	client, _ := New(server.URL)
	result, res, err := client.Repositories.CreateStatus(context.Background(), "diaspora/diaspora", "18f3e63d05582537db6d183d9d557be09e1f90c8", in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Status", testStatus(result))
}

func TestRepositoryHookFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Repositories.FindHook(context.Background(), "diaspora/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Hook", testHook(result))
}

func TestRepositoryHookList(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Repositories.ListHooks(context.Background(), "diaspora/diaspora", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(result), 1; got != want {
		t.Errorf("Want %d hooks, got %d", want, got)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
	t.Run("Hook", testHook(result[0]))
}

func TestRepositoryHookDelete(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	_, err := client.Repositories.DeleteHook(context.Background(), "diaspora/diaspora", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestRepositoryHookCreate(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	in := &scm.HookInput{
		Target: "http://example.com/hook",
		Events: scm.HookEvents{
			Push: true,
		},
	}

	client, _ := New(server.URL)
	result, res, err := client.Repositories.CreateHook(context.Background(), "diaspora/diaspora", in)
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Hook", testHook(result))
}

func testRepository(repository *scm.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := repository.ID, "178504"; got != want {
			t.Errorf("Want repository ID %q, got %q", want, got)
		}
		if got, want := repository.Name, "diaspora"; got != want {
			t.Errorf("Want repository Name %q, got %q", want, got)
		}
		if got, want := repository.Namespace, "diaspora"; got != want {
			t.Errorf("Want repository Namespace %q, got %q", want, got)
		}
		if got, want := repository.Branch, "master"; got != want {
			t.Errorf("Want repository Branch %q, got %q", want, got)
		}
		if got, want := repository.Private, false; got != want {
			t.Errorf("Want repository Private %v, got %v", want, got)
		}
	}
}

func testPermissions(perms *scm.Perm) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := perms.Pull, true; got != want {
			t.Errorf("Want permission Pull %v, got %v", want, got)
		}
		if got, want := perms.Push, false; got != want {
			t.Errorf("Want permission Push %v, got %v", want, got)
		}
		if got, want := perms.Admin, false; got != want {
			t.Errorf("Want permission Admin %v, got %v", want, got)
		}
	}
}

func testHook(hook *scm.Hook) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := hook.ID, 1; got != want {
			t.Errorf("Want hook ID %v, got %v", want, got)
		}
		if got, want := hook.Active, true; got != want {
			t.Errorf("Want hook Active %v, got %v", want, got)
		}
		if got, want := hook.Target, "http://example.com/hook"; got != want {
			t.Errorf("Want hook Target %v, got %v", want, got)
		}
	}
}

func testStatus(status *scm.Status) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := status.Target, "https://gitlab.example.com/thedude/gitlab-ce/builds/91"; got != want {
			t.Errorf("Want status Target %v, got %v", want, got)
		}
		if got, want := status.State, scm.StatePending; got != want {
			t.Errorf("Want status State %v, got %v", want, got)
		}
		if got, want := status.Label, "default"; got != want {
			t.Errorf("Want status Label %v, got %v", want, got)
		}
		if got, want := status.Desc, "the dude abides"; got != want {
			t.Errorf("Want status Desc %v, got %v", want, got)
		}
	}
}
