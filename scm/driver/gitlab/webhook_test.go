// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

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
		// branch hooks
		// {
		// 	sig:    "c0e2b3a5e115811f8485dcb7728d50b2ce1a40d631d1cd5edf8c32ee8100b6f4",
		// 	event:  "create",
		// 	before: "samples/branch_create.json",
		// 	after:  "samples/branch_create.json.golden",
		// 	obj:    new(scm.BranchHook),
		// },
		{
			sig:    "root",
			event:  "Push Hook",
			before: "samples/branch_delete.json",
			after:  "samples/branch_delete.parsed.json",
			obj:    new(scm.BranchHook),
		},
		// // tag hooks
		// {
		// 	sig:    "e755dedc23bdf817a994c1f1705da149467f2e55b195f6e282383f80e9a53358",
		// 	event:  "create",
		// 	before: "samples/tag_create.json",
		// 	after:  "samples/tag_create.parsed.json",
		// 	obj:    new(scm.TagHook),
		// },
		// {
		// 	sig:    "99b0128bef639616331037e0b04624ad8b11dcdb4f2fc421ff22c34fca2754eb",
		// 	event:  "delete",
		// 	before: "samples/tag_delete.json",
		// 	after:  "samples/tag_delete.parsed.json",
		// 	obj:    new(scm.TagHook),
		// },
		// // push hooks
		// {
		// 	sig:    "035464962b775451f7e418c85a2302a8ab8949a1da9ced206b9eacac543cb9eb",
		// 	event:  "push",
		// 	before: "samples/push.json",
		// 	after:  "samples/push.parsed.json",
		// 	obj:    new(scm.PushHook),
		// },
		// // issue hooks
		// {
		// 	sig:    "aa45894e45f34ca8dbd38688ab6806ba7041245bf0d27ecf90fe959075c62943",
		// 	event:  "issues",
		// 	before: "samples/issues_opened.json",
		// 	after:  "samples/issues_opened.parsed.json",
		// 	obj:    new(scm.IssueHook),
		// },
		// // issue comment hooks
		// {
		// 	sig:    "2dee1c4ff5bd25899568dbd696f28cfc37a5dd7db99da042cf87c8482cbbde78",
		// 	event:  "issue_comment",
		// 	before: "samples/issue_comment_created.json",
		// 	after:  "samples/issue_comment_created.parsed.json",
		// 	obj:    new(scm.IssueCommentHook),
		// },
		// // pull request hooks
		// {
		// 	sig:    "f9522c6e7862507971ae7e250d5a93fa1629ab16bd97fb2de441a08d804aabe3",
		// 	event:  "pull_request",
		// 	before: "samples/pull_request_opened.json",
		// 	after:  "samples/pull_request_opened.parsed.json",
		// 	obj:    new(scm.PullRequestHook),
		// },
		// {
		// 	sig:    "261d7f124e251f7b2ac309e4b710f3a8cb588076bdec8aa57b3ec8586cc4e790",
		// 	event:  "pull_request",
		// 	before: "samples/pull_request_edited.json",
		// 	after:  "samples/pull_request_edited.parsed.json",
		// 	obj:    new(scm.PullRequestHook),
		// },
		// {
		// 	sig:    "86b0e5eac0561c7fc479b423ec148f743b4a8a24c29fec143c26b81741001baa",
		// 	event:  "pull_request",
		// 	before: "samples/pull_request_synchronized.json",
		// 	after:  "samples/pull_request_synchronized.parsed.json",
		// 	obj:    new(scm.PullRequestHook),
		// },
		// {
		// 	sig:    "b2cdceb63461ee5dc7d6d609f285d8498b87abaaaf7e814d784984a9a8ffce1b",
		// 	event:  "pull_request",
		// 	before: "samples/pull_request_closed.json",
		// 	after:  "samples/pull_request_closed.parsed.json",
		// 	obj:    new(scm.PullRequestHook),
		// },
		// // pull request comment hooks
		// {
		// 	sig:    "9edf3260d727e29d906bdb10c8a099666a3df4f033084e244fd56ef8828c9bea",
		// 	event:  "issue_comment",
		// 	before: "samples/pull_request_comment_created.json",
		// 	after:  "samples/pull_request_comment_created.parsed.json",
		// 	obj:    new(scm.PullRequestCommentHook),
		// },
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
		r.Header.Set("X-Gitlab-Event", test.event)
		r.Header.Set("X-Gitlab-Token", test.sig)
		r.Header.Set("X-Request-Id", "ee8d97b4-1479-43f1-9cac-fbbd1b80da55")

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
	f, _ := ioutil.ReadFile("samples/branch_delete.json")
	r, _ := http.NewRequest("GET", "/", bytes.NewBuffer(f))
	r.Header.Set("X-Gitlab-Event", "Push Hook")
	r.Header.Set("X-Gitlab-Token", "void")
	r.Header.Set("X-Request-Id", "ee8d97b4-1479-43f1-9cac-fbbd1b80da55")

	s := new(webhookService)
	_, err := s.Parse(r, secretFunc)
	if err != scm.ErrSignatureInvalid {
		t.Errorf("Expect invalid signature error, got %v", err)
	}
}

func secretFunc(interface{}) (string, error) {
	return "", nil
}
