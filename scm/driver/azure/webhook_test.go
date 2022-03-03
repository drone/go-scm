// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
		// pull request events
		//
		// pull request created
		{
			sig:    "71295b197fa25f4356d2fb9965df3f2379d903d7",
			event:  "pullrequest:created",
			before: "testdata/webhooks/pr_created.json",
			after:  "testdata/webhooks/pr_created.json.golden",
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
		r, _ := http.NewRequest("GET", "/?secret=71295b197fa25f4356d2fb9965df3f2379d903d7", buf)
		r.Header.Set("x-event-key", test.event)

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
			//	_ = json.NewEncoder(os.Stdout).Encode(o)
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

func secretFunc(scm.Webhook) (string, error) {
	return "71295b197fa25f4356d2fb9965df3f2379d903d7", nil
}
