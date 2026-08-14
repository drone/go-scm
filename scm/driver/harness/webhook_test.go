// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
)

func TestWebhooks(t *testing.T) {
	tests := []struct {
		event  string
		before string
		after  string
		obj    interface{}
	}{
		//
		// branch events
		//
		// push branch create
		{
			event:  "branch_created",
			before: "testdata/webhooks/branch_create.json",
			after:  "testdata/webhooks/branch_create.json.golden",
			obj:    new(scm.BranchHook),
		},
		// push branch update
		{
			event:  "branch_updated",
			before: "testdata/webhooks/branch_updated.json",
			after:  "testdata/webhooks/branch_updated.json.golden",
			obj:    new(scm.PushHook),
		},
		// push branch delete
		{
			event:  "branch_deleted",
			before: "testdata/webhooks/branch_delete.json",
			after:  "testdata/webhooks/branch_delete.json.golden",
			obj:    new(scm.BranchHook),
		},
		//
		// tag events
		//
		// push tag create
		{
			event:  "tag_created",
			before: "testdata/webhooks/tag_create.json",
			after:  "testdata/webhooks/tag_create.json.golden",
			obj:    new(scm.TagHook),
		},
		// push tag update
		{
			event:  "tag_updated",
			before: "testdata/webhooks/tag_update.json",
			after:  "testdata/webhooks/tag_update.json.golden",
			obj:    new(scm.PushHook),
		},
		// push tag delete
		{
			event:  "tag_deleted",
			before: "testdata/webhooks/tag_delete.json",
			after:  "testdata/webhooks/tag_delete.json.golden",
			obj:    new(scm.TagHook),
		},

		//
		// pull request events
		//
		// pull request opened
		{
			event:  "pullreq_created",
			before: "testdata/webhooks/pull_request_opened.json",
			after:  "testdata/webhooks/pull_request_opened.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request reopened
		{
			event:  "pullreq_reopened",
			before: "testdata/webhooks/pull_request_reopened.json",
			after:  "testdata/webhooks/pull_request_reopened.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request branch updated
		{
			event:  "pullreq_branch_updated",
			before: "testdata/webhooks/pull_request_branch_updated.json",
			after:  "testdata/webhooks/pull_request_branch_updated.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request comment created
		{
			event:  "pullreq_comment_created",
			before: "testdata/webhooks/pull_request_comment_created.json",
			after:  "testdata/webhooks/pull_request_comment_created.json.golden",
			obj:    new(scm.PullRequestCommentHook),
		},
		// pull request closed
		{
			event:  "pullreq_reopened",
			before: "testdata/webhooks/pull_request_closed.json",
			after:  "testdata/webhooks/pull_request_closed.json.golden",
			obj:    new(scm.PullRequestHook),
		},
		// pull request merged
		{
			event:  "pullreq_reopened",
			before: "testdata/webhooks/pull_request_merged.json",
			after:  "testdata/webhooks/pull_request_merged.json.golden",
			obj:    new(scm.PullRequestHook),
		},

		//
		// merge queue events
		//
		// merge queue checks requested
		{
			event:  "merge_queue_checks_requested",
			before: "testdata/webhooks/merge_queue_checks_requested.json",
			after:  "testdata/webhooks/merge_queue_checks_requested.json.golden",
			obj:    new(scm.MergeQueueHook),
		},
		// merge queue checks canceled
		{
			event:  "merge_queue_checks_canceled",
			before: "testdata/webhooks/merge_queue_checks_canceled.json",
			after:  "testdata/webhooks/merge_queue_checks_canceled.json.golden",
			obj:    new(scm.MergeQueueHook),
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
		r.Header.Set("X-Harness-Trigger", test.event)

		s := new(webhookService)
		o, err := s.Parse(r, secretFunc)
		if err != nil && err != scm.ErrSignatureInvalid {
			t.Error(err)
			continue
		}

		err = json.Unmarshal(after, test.obj)
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

		// switch event := o.(type) {
		// case *scm.PushHook:
		// 	if !strings.HasPrefix(event.Ref, "refs/") {
		// 		t.Errorf("Push hook reference must start with refs/")
		// 	}
		// case *scm.BranchHook:
		// 	if strings.HasPrefix(event.Ref.Name, "refs/") {
		// 		t.Errorf("Branch hook reference must not start with refs/")
		// 	}
		// case *scm.TagHook:
		// 	if strings.HasPrefix(event.Ref.Name, "refs/") {
		// 		t.Errorf("Branch hook reference must not start with refs/")
		// 	}
		// }
	}
}
func secretFunc(scm.Webhook) (string, error) {
	return "topsecret", nil
}

func TestConvertHookCommitDate(t *testing.T) {
	tests := []struct {
		name string
		when string
		want string
	}{
		{name: "utc", when: "2023-12-05T11:59:39Z", want: "2023-12-05T11:59:39Z"},
		{name: "offset", when: "2024-03-07T03:18:51-08:00", want: "2024-03-07T03:18:51-08:00"},
		{name: "fractional seconds", when: "2023-02-09T17:12:10.976Z", want: "2023-02-09T17:12:10.976Z"},
		{name: "missing", when: "", want: ""},
		{name: "malformed", when: "not-a-timestamp", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raw := []byte(`{
				"sha": "700b3dab8e7a5cebf5e1ce54e7dd5bde60099912",
				"message": "Asd",
				"author": {"identity": {"name": "Abhinav Singh", "email": "abhinav.singh@harness.io"}, "when": "` +
				test.when + `"},
				"committer": {"identity": {"name": "Harness", "email": "noreply@harness.io"}, "when": "` +
				test.when + `"}
			}`)

			in := new(hookCommit)
			if err := json.Unmarshal(raw, in); err != nil {
				t.Fatal(err)
			}

			got := convertHookCommit(*in)

			var want time.Time
			if test.want != "" {
				parsed, err := time.Parse(time.RFC3339, test.want)
				if err != nil {
					t.Fatal(err)
				}
				want = parsed
			}

			if !got.Author.Date.Equal(want) {
				t.Errorf("author date = %v, want %v", got.Author.Date, want)
			}
			if !got.Committer.Date.Equal(want) {
				t.Errorf("committer date = %v, want %v", got.Committer.Date, want)
			}
		})
	}
}
