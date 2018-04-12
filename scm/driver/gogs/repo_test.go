// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gogs

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
)

func testRepos(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testRepoFind(client))
		t.Run("FindPerm", testRepoFindPerm(client))
		t.Run("NotFound", testRepoNotFound(client))
		t.Run("List", testRepoList(client))
		t.Run("Hooks", testHooks(client))
		t.Run("Statuses", testStatuses(client))
	}
}

//
// repository sub-tests
//

func testRepoFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Repositories.Find(context.Background(), "gogits/gogs")
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Repository", testRepository(result))
			t.Run("Permissions", testPermissions(result.Perm))
		}
	}
}

func testRepoFindPerm(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Repositories.FindPerms(context.Background(), "gogits/gogs")
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Permissions", testPermissions(result))
		}
	}
}

func testRepoList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Repositories.List(context.Background(), scm.ListOptions{})
		if err != nil {
			t.Error(err)
		} else if got, want := len(result), 1; got != want {
			t.Errorf("Want %d repositories, got %d", want, got)
		} else {
			t.Run("Repository", testRepository(result[0]))
			t.Run("Permissions", testPermissions(result[0].Perm))
		}
	}
}

func testRepoNotFound(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Repositories.FindPerms(context.Background(), "gogits/go-gogs-client")
		if err == nil {
			t.Errorf("Expect Not Found error")
		} else if got, want := err.Error(), "Not Found"; got != want {
			t.Errorf("Want error %q, got %q", want, got)
		}
	}
}

//
// hook sub-tests
//

func testHooks(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Find", testHookFind(client))
		t.Run("List", testHookList(client))
		t.Run("Create", testHookCreate(client))
		t.Run("Delete", testHookDelete(client))
	}
}

func testHookFind(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Repositories.FindHook(context.Background(), "gogits/gogs", 20)
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Hook", testHook(result))
		}
	}
}

func testHookList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Repositories.ListHooks(context.Background(), "gogits/gogs", scm.ListOptions{})
		if err != nil {
			t.Error(err)
		} else if got, want := len(result), 1; got != want {
			t.Errorf("Want %d hooks, got %d", want, got)
		} else {
			t.Run("Hook", testHook(result[0]))
		}
	}
}

func testHookCreate(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		result, _, err := client.Repositories.CreateHook(context.Background(), "gogits/gogs", &scm.HookInput{})
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Hook", testHook(result))
		}
	}
}

func testHookDelete(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := client.Repositories.DeleteHook(context.Background(), "gogits/gogs", 20)
		if err != nil {
			t.Error(err)
		}
	}
}

//
// status sub-tests
//

func testStatuses(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("List", testStatusList(client))
		t.Run("Create", testStatusCreate(client))
	}
}

func testStatusList(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Repositories.ListStatus(context.Background(), "gogits/gogs", "master", scm.ListOptions{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

func testStatusCreate(client *scm.Client) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := client.Repositories.CreateStatus(context.Background(), "gogits/gogs", "master", &scm.StatusInput{})
		if err != scm.ErrNotSupported {
			t.Errorf("Expect Not Supported error")
		}
	}
}

//
// struct value sub-tests
//

func testRepository(repository *scm.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := repository.ID, "1"; got != want {
			t.Errorf("Want repository ID %v, got %v", want, got)
		}
		if got, want := repository.Name, "gogs"; got != want {
			t.Errorf("Want repository Name %q, got %q", want, got)
		}
		if got, want := repository.Namespace, "gogits"; got != want {
			t.Errorf("Want repository Owner %q, got %q", want, got)
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
		if got, want := hook.ID, 20; got != want {
			t.Errorf("Want hook ID %v, got %v", want, got)
		}
		if got, want := hook.Active, true; got != want {
			t.Errorf("Want hook Active %v, got %v", want, got)
		}
		if got, want := hook.Target, "http://gogs.io"; got != want {
			t.Errorf("Want hook Target %v, got %v", want, got)
		}
	}
}
