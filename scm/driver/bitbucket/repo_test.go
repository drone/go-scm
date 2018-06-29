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
	result, res, err := client.Repositories.Find(context.Background(), "octocat/hello-world")
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
	result, res, err := client.Repositories.FindPerms(context.Background(), "octocat/hello-world")
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
	if got, want := err.Error(), "Not Found"; got != want {
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
	result, res, err := client.Repositories.ListStatus(context.Background(), "octocat/hello-world", "6dcb09b5b57875f334f61aebed695e2e4193db5e", scm.ListOptions{Size: 30, Page: 1})
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
	result, res, err := client.Repositories.CreateStatus(context.Background(), "octocat/hello-world", "6dcb09b5b57875f334f61aebed695e2e4193db5e", in)
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
	result, res, err := client.Repositories.FindHook(context.Background(), "octocat/hello-world", 1)
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
	result, res, err := client.Repositories.ListHooks(context.Background(), "octocat/hello-world", scm.ListOptions{Page: 1, Size: 30})
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
	_, err := client.Repositories.DeleteHook(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestRepositoryHookCreate(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, res, err := client.Repositories.CreateHook(context.Background(), "octocat/hello-world", &scm.HookInput{})
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
		if got, want := repository.ID, "1296269"; got != want {
			t.Errorf("Want repository ID %q, got %q", want, got)
		}
		if got, want := repository.Name, "Hello-World"; got != want {
			t.Errorf("Want repository Name %q, got %q", want, got)
		}
		if got, want := repository.Namespace, "octocat"; got != want {
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
		if got, want := hook.ID, 1; got != want {
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
		if got, want := status.Label, "continuous-integration/jenkins"; got != want {
			t.Errorf("Want status Label %v, got %v", want, got)
		}
		if got, want := status.Desc, "Build has completed successfully"; got != want {
			t.Errorf("Want status Desc %v, got %v", want, got)
		}
	}
}
