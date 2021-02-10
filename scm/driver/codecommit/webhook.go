// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codecommit

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"
	sns "github.com/robbiet480/go.sns"
)

const (
	eventSrcCodeCommit = "aws:codecommit"
)

type webhookService struct {
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (scm.Webhook, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}

	var notification sns.Payload
	err = json.Unmarshal([]byte(data), &notification)
	if err != nil {
		return nil, err
	}

	err = notification.VerifyPayload()
	if err != nil {
		return nil, scm.ErrSignatureInvalid
	}

	src := new(message)
	err = json.Unmarshal([]byte(notification.Message), src)
	if err != nil {
		return nil, err
	}

	t, _ := time.Parse(time.RFC3339, notification.Timestamp)
	for _, record := range src.Records {
		if record.EventSource != eventSrcCodeCommit {
			continue
		}

		for _, reference := range record.Codecommit.References {
			hook, err := s.parseReference(reference, record, t)
			if err != nil {
				return nil, err
			}

			if hook != nil {
				return hook, nil
			}
		}
	}

	return nil, scm.ErrUnknownEvent
}

// AWS codecommit triggers handles
// 1. push to existing branch
// 2. branch create/delete
// 3. tag create/delete
func (s *webhookService) parseReference(ref reference, r record, t time.Time) (scm.Webhook, error) {
	action := getAction(ref)
	var dst scm.Webhook
	if strings.HasPrefix(ref.Ref, "refs/tags/") {
		dst = convertTagHook(ref, r, t, action)
	} else if strings.HasPrefix(ref.Ref, "refs/heads/") {
		if action == scm.ActionUnknown {
			dst = convertPushHook(ref, r, t)
		} else {
			dst = convertBranchHook(ref, r, t, action)
		}
	} else {
		return nil, scm.ErrUnknownEvent
	}
	return dst, nil
}

func convertTagHook(ref reference, r record, t time.Time, action scm.Action) *scm.TagHook {
	tag := ref.Ref[len("refs/tags/"):]
	repo := getRepo(r.EventSourceARN)
	return &scm.TagHook{
		Ref: scm.Reference{
			Name: tag,
			Path: ref.Ref,
			Sha:  ref.Commit,
		},
		Repo: scm.Repository{
			ID:      r.EventSourceARN,
			Name:    repo,
			Private: true,
		},
		Sender: scm.User{
			Login:   r.UserIdentityARN,
			Created: t,
			Updated: t,
		},
		Action: action,
	}
}

func convertPushHook(ref reference, r record, t time.Time) *scm.PushHook {
	branch := ref.Ref[len("refs/heads/"):]
	repo := getRepo(r.EventSourceARN)
	c := scm.Commit{
		Sha: ref.Commit,
		Author: scm.Signature{
			Login: r.UserIdentityARN,
			Date:  t,
		},
		Committer: scm.Signature{
			Login: r.UserIdentityARN,
			Date:  t,
		},
	}

	return &scm.PushHook{
		Ref:    ref.Ref,
		Commit: c,
		Repo: scm.Repository{
			ID:      r.EventSourceARN,
			Name:    repo,
			Branch:  branch,
			Private: true,
		},
		Sender: scm.User{
			Login:   r.UserIdentityARN,
			Created: t,
			Updated: t,
		},
		Commits: []scm.Commit{c},
	}
}

func convertBranchHook(ref reference, r record, t time.Time, action scm.Action) *scm.BranchHook {
	branch := ref.Ref[len("refs/heads/"):]
	repo := getRepo(r.EventSourceARN)
	return &scm.BranchHook{
		Ref: scm.Reference{
			Name: branch,
			Path: ref.Ref,
			Sha:  ref.Commit,
		},
		Repo: scm.Repository{
			ID:      r.EventSourceARN,
			Name:    repo,
			Branch:  branch,
			Private: true,
			Created: t,
			Updated: t,
		},
		Sender: scm.User{
			Login:   r.UserIdentityARN,
			Created: t,
			Updated: t,
		},
		Action: action,
	}
}

func getRepo(eventSourceARN string) string {
	parts := strings.Split(eventSourceARN, ":")
	if len(parts) != 6 {
		return ""
	}

	return parts[5]
}

func getAction(r reference) scm.Action {
	if r.Created {
		return scm.ActionCreate
	}
	if r.Deleted {
		return scm.ActionDelete
	}
	return scm.ActionUnknown
}

type (
	reference struct {
		Commit  string `json:"commit"`
		Ref     string `json:"ref"`
		Created bool   `json:"created"`
		Deleted bool   `json:"deleted"`
	}

	record struct {
		AwsRegion  string `json:"awsRegion"`
		Codecommit struct {
			References []reference `json:"references"`
		} `json:"codecommit"`
		EventID          string `json:"eventId"`
		EventName        string `json:"eventName"`
		EventPartNumber  int    `json:"eventPartNumber"`
		EventSource      string `json:"eventSource"`
		EventSourceARN   string `json:"eventSourceARN"`
		EventTime        string `json:"eventTime"`
		EventTotalParts  int    `json:"eventTotalParts"`
		EventTriggerName string `json:"eventTriggerName"`
		EventVersion     string `json:"eventVersion"`
		UserIdentityARN  string `json:"userIdentityARN"`
	}

	message struct {
		Records []record `json:"records"`
	}
)
