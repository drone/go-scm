// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "repo:refs_changed",
			before: "testdata/webhooks/push.json",
			after:  "testdata/webhooks/push.json.golden",
			obj:    new(scm.PushHook),
		},
		// v5 test
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "repo:refs_changed",
			before: "testdata/webhooks/push_v5.json",
			after:  "testdata/webhooks/push_v5.json.golden",
			obj:    new(scm.PushHook),
		},
		//
		// tag events
		//

		// create
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "repo:refs_changed",
			before: "testdata/webhooks/push_tag_create.json",
			after:  "testdata/webhooks/push_tag_create.json.golden",
			obj:    new(scm.TagHook),
		},
		// delete
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "repo:refs_changed",
			before: "testdata/webhooks/push_tag_delete.json",
			after:  "testdata/webhooks/push_tag_delete.json.golden",
			obj:    new(scm.TagHook),
		},

		//
		// branch events
		//

		// create
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "repo:refs_changed",
			before: "testdata/webhooks/push_branch_create.json",
			after:  "testdata/webhooks/push_branch_create.json.golden",
			obj:    new(scm.BranchHook),
		},
		// delete
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "repo:refs_changed",
			before: "testdata/webhooks/push_branch_delete.json",
			after:  "testdata/webhooks/push_branch_delete.json.golden",
			obj:    new(scm.BranchHook),
		},

		//
		// pull request events
		//

		// pull request opened
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pr:opened",
			before: "testdata/webhooks/pr_open.json",
			after:  "testdata/webhooks/pr_open.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request source branch updated
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pr:from_ref_updated",
			before: "testdata/webhooks/pr_from_ref_updated.json",
			after:  "testdata/webhooks/pr_from_ref_updated.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request modified, target branch updated
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pr:modified",
			before: "testdata/webhooks/pr_modified_toref.json",
			after:  "testdata/webhooks/pr_modified_toref.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request modified, no change to toref
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pr:modified",
			before: "testdata/webhooks/pr_modified_meta.json",
			after:  "testdata/webhooks/pr_modified_meta.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request fulfilled (merged)
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pr:merged",
			before: "testdata/webhooks/pr_merged.json",
			after:  "testdata/webhooks/pr_merged.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request rejected (closed, declined)
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pr:declined",
			before: "testdata/webhooks/pr_declined.json",
			after:  "testdata/webhooks/pr_declined.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request deleted
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pr:deleted",
			before: "testdata/webhooks/pr_deleted.json",
			after:  "testdata/webhooks/pr_deleted.json.golden",
			obj:    new(scm.PullRequestHook),
		},

		//
		// real-world payloads
		//

		// FNB (Bitbucket DC 10.3) branch create push webhook
		{
			sig:    "",
			event:  "repo:refs_changed",
			before: "testdata/fnb_push_branch_create.json",
			after:  "testdata/fnb_push_branch_create.json.golden",
			obj:    new(scm.BranchHook),
		},
		// FNB (Bitbucket DC 10.3) push update webhook — ISO timestamp strings in commit objects
		{
			sig:    "",
			event:  "repo:refs_changed",
			before: "testdata/fnb_push_update.json",
			after:  "testdata/fnb_push_update.json.golden",
			obj:    new(scm.PushHook),
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
		r.Header.Set("X-Event-Key", test.event)

		s := new(webhookService)
		o, err := s.Parse(r, secretFunc)
		if err != nil && err != scm.ErrSignatureInvalid {
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
			_ = json.NewEncoder(os.Stdout).Encode(o)
		}

		switch event := o.(type) {
		case *scm.PushHook:
			if !strings.HasPrefix(event.Ref, "refs/") {
				t.Errorf("Push hook reference must start with refs/")
			}
		case *scm.BranchHook:
			if strings.HasPrefix(event.Ref.Name, "refs/") {
				t.Errorf("Branch hook reference must not start with refs/")
			}
		case *scm.TagHook:
			if strings.HasPrefix(event.Ref.Name, "refs/") {
				t.Errorf("Branch hook reference must not start with refs/")
			}
		}
	}
}

func TestWebhookInvalid(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/webhooks/push.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
	r.Header.Set("X-Event-Key", "repo:refs_changed")
	r.Header.Set("X-Hub-Signature", "sha256=380f462cd2e160b84765144beabdad2e930a7ec5")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != scm.ErrSignatureInvalid {
		t.Errorf("Expect invalid signature error, got %v", err)
	}
}

func TestWebhookVerified(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/webhooks/push.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
	r.Header.Set("X-Event-Key", "repo:refs_changed")
	r.Header.Set("X-Hub-Signature", "sha256=c90565fa018f3039414a7929c9187a147f1ac463076961c4cf411e3c67c541f8")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != nil {
		t.Errorf("Expect valid signature error, got %v", err)
	}
}

// TestWebhooksIsoDates verifies that PR webhooks with ISO 8601 date strings
// (Bitbucket Data Center 10.3+) produce the same output as epoch-ms payloads.
func TestWebhooksIsoDates(t *testing.T) {
	tests := []struct {
		name       string
		event      string
		isoFile    string
		goldenFile string
		obj        interface{}
	}{
		{
			name:       "pr_opened ISO dates",
			event:      "pr:opened",
			isoFile:    "testdata/webhooks/pr_open_iso_dates.json",
			goldenFile: "testdata/webhooks/pr_open_iso_dates.json.golden",
			obj:        new(scm.PullRequestHook),
		},
		{
			name:       "pr_merged ISO dates",
			event:      "pr:merged",
			isoFile:    "testdata/webhooks/pr_merged_iso_dates.json",
			goldenFile: "testdata/webhooks/pr_merged_iso_dates.json.golden",
			obj:        new(scm.PullRequestHook),
		},
		{
			name:       "pr_declined ISO dates",
			event:      "pr:declined",
			isoFile:    "testdata/webhooks/pr_declined_iso_dates.json",
			goldenFile: "testdata/webhooks/pr_declined_iso_dates.json.golden",
			obj:        new(scm.PullRequestHook),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before, err := ioutil.ReadFile(test.isoFile)
			if err != nil {
				t.Fatal(err)
			}
			after, err := ioutil.ReadFile(test.goldenFile)
			if err != nil {
				t.Fatal(err)
			}

			r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(before))
			r.Header.Set("X-Event-Key", test.event)

			s := new(webhookService)
			got, err := s.Parse(r, secretFunc)
			if err != nil && err != scm.ErrSignatureInvalid {
				t.Fatalf("Parse error: %v", err)
			}

			if err := json.Unmarshal(after, &test.obj); err != nil {
				t.Fatalf("golden unmarshal error: %v", err)
			}

			if diff := cmp.Diff(test.obj, got); diff != "" {
				t.Errorf("ISO dates webhook diff from epoch-ms golden:\n%s", diff)
			}
		})
	}
}

// TestWebhookPRIsoDatesNoError verifies that a PR webhook with ISO 8601 dates
// does not return a parse error (regression test for PIPE-34150).
func TestWebhookPRIsoDatesNoError(t *testing.T) {
	events := []struct {
		file  string
		event string
	}{
		{"testdata/webhooks/pr_open_iso_dates.json", "pr:opened"},
		{"testdata/webhooks/pr_merged_iso_dates.json", "pr:merged"},
		{"testdata/webhooks/pr_declined_iso_dates.json", "pr:declined"},
	}

	for _, ev := range events {
		t.Run(ev.event, func(t *testing.T) {
			f, err := ioutil.ReadFile(ev.file)
			if err != nil {
				t.Fatal(err)
			}
			r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
			r.Header.Set("X-Event-Key", ev.event)

			s := new(webhookService)
			hook, err := s.Parse(r, secretFunc)
			if err != nil && err != scm.ErrSignatureInvalid {
				t.Errorf("unexpected error for ISO date PR webhook: %v", err)
			}
			if hook == nil {
				t.Errorf("expected non-nil hook, got nil")
			}
		})
	}
}

func secretFunc(scm.Webhook) (string, error) {
	return "71295b197fa25f4356d2fb9965df3f2379d903d7", nil
}
