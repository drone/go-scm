// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
)

func TestWebhooks(t *testing.T) {
	tests := []struct {
		sig    string
		event  string
		before string
		after  string
		obj    interface{}
	}{
		//
		// push events
		//

		// push hooks
		{
			sig:    "sha1=8256c70004120d7241f9833be7378f40060c0763",
			event:  "push",
			before: "samples/push.json",
			after:  "samples/push.json.golden",
			obj:    new(scm.PushHook),
		},
		// push tag create hooks
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "push",
			before: "samples/push_tag.json",
			after:  "samples/push_tag.json.golden",
			obj:    new(scm.PushHook),
		},
		// push tag delete hooks
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "push",
			before: "samples/push_tag_delete.json",
			after:  "samples/push_tag_delete.json.golden",
			obj:    new(scm.PushHook),
		},
		// push branch create
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "push",
			before: "samples/push_branch_create.json",
			after:  "samples/push_branch_create.json.golden",
			obj:    new(scm.PushHook),
		},
		// push branch delete
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "push",
			before: "samples/push_branch_delete.json",
			after:  "samples/push_branch_delete.json.golden",
			obj:    new(scm.PushHook),
		},

		//
		// branch events
		//

		// push branch create
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "create",
			before: "samples/branch_create.json",
			after:  "samples/branch_create.json.golden",
			obj:    new(scm.BranchHook),
		},
		// push branch delete
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "delete",
			before: "samples/branch_delete.json",
			after:  "samples/branch_delete.json.golden",
			obj:    new(scm.BranchHook),
		},

		//
		// tag events
		//

		// push tag create
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "create",
			before: "samples/tag_create.json",
			after:  "samples/tag_create.json.golden",
			obj:    new(scm.TagHook),
		},
		// push tag delete
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "delete",
			before: "samples/tag_delete.json",
			after:  "samples/tag_delete.json.golden",
			obj:    new(scm.TagHook),
		},

		//
		// pull request events
		//

		// pull request synced
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pull_request",
			before: "samples/pr_sync.json",
			after:  "samples/pr_sync.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request opened
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pull_request",
			before: "samples/pr_opened.json",
			after:  "samples/pr_opened.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request closed
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pull_request",
			before: "samples/pr_closed.json",
			after:  "samples/pr_closed.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request reopened
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pull_request",
			before: "samples/pr_reopened.json",
			after:  "samples/pr_reopened.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request edited
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pull_request",
			before: "samples/pr_edited.json",
			after:  "samples/pr_edited.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request labeled
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pull_request",
			before: "samples/pr_labeled.json",
			after:  "samples/pr_labeled.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request unlabeled
		{
			sig:    "sha1=71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pull_request",
			before: "samples/pr_unlabeled.json",
			after:  "samples/pr_unlabeled.json.golden",
			obj:    new(scm.PullRequestHook),
		},
	}

	for _, test := range tests {
		before, err := ioutil.ReadFile(test.before)
		if err != nil {
			t.Error(err)
			continue
		}
		after, err := ioutil.ReadFile(test.after)
		if err != nil {
			t.Error(err)
			continue
		}

		buf := bytes.NewBuffer(before)
		r, _ := http.NewRequest("GET", "/", buf)
		r.Header.Set("X-GitHub-Event", test.event)
		r.Header.Set("X-Hub-Signature", test.sig)
		r.Header.Set("X-GitHub-Delivery", "f2467dea-70d6-11e8-8955-3c83993e0aef")

		s := new(webhookService)
		o, err := s.Parse(r, secretFunc)
		if err != nil {
			t.Error(err)
			continue
		}

		err = json.Unmarshal(after, &test.obj)
		if err != nil {
			t.Error(err)
			continue
		}

		if diff := cmp.Diff(test.obj, o); diff != "" {
			t.Errorf("Error unmarshaling %s", test.before)
			t.Log(diff)

			// debug only. remove once implemented
			json.NewEncoder(os.Stdout).Encode(o)
		}
	}
}

func TestWebhookInvalid(t *testing.T) {
	f, _ := ioutil.ReadFile("samples/push.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
	r.Header.Set("X-GitHub-Event", "push")
	r.Header.Set("X-GitHub-Delivery", "ee8d97b4-1479-43f1-9cac-fbbd1b80da55")
	r.Header.Set("X-Hub-Signature", "sha1=380f462cd2e160b84765144beabdad2e930a7ec5")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != scm.ErrSignatureInvalid {
		t.Errorf("Expect invalid signature error, got %v", err)
	}
}

func secretFunc(interface{}) (string, error) {
	return "", nil
	// return "root", nil
}
