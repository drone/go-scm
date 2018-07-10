// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/hmac"
)

// TODO(bradrydzewski) push hook does not include commit message
// TODO(bradrydzewski) push hook does not include commit link
// TODO(bradrydzewski) push hook does not include repository git+http link
// TODO(bradrydzewski) push hook does not include repository git+ssh link
// TODO(bradrydzewski) push hook does not include repository html link
// TODO(bradrydzewski) missing pull request synchrnoized webhook. See https://jira.atlassian.com/browse/BSERV-10279
// TODO(bradrydzewski) pr hook does not include repository git+http link
// TODO(bradrydzewski) pr hook does not include repository git+ssh link
// TODO(bradrydzewski) pr hook does not include repository html link

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (scm.Webhook, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}

	var hook scm.Webhook
	switch req.Header.Get("X-Event-Key") {
	case "repo:refs_changed":
		hook, err = s.parsePushHook(data)
	case "pr:opened", "pr:declined", "pr:merged":
		hook, err = s.parsePullRequest(data)
	}
	if err != nil {
		return nil, err
	}
	if hook == nil {
		return nil, nil
	}

	// get the gogs signature key to verify the payload
	// signature. If no key is provided, no validation
	// is performed.
	key, err := fn(hook)
	if err != nil {
		return hook, err
	} else if key == "" {
		return hook, nil
	}

	sig := req.Header.Get("X-Hub-Signature")
	if !hmac.ValidatePrefix(data, []byte(key), sig) {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePushHook(data []byte) (scm.Webhook, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	if err != nil {
		return nil, err
	}
	if len(dst.Changes) == 0 {
		return nil, errors.New("Push hook has empty changeset")
	}
	change := dst.Changes[0]
	switch {
	case change.Ref.Type == "BRANCH" && change.Type != "UPDATE":
		return convertBranchHook(dst), nil
	case change.Ref.Type == "TAG":
		return convertTagHook(dst), nil
	default:
		return convertPushHook(dst), err
	}
}

func (s *webhookService) parsePullRequest(data []byte) (scm.Webhook, error) {
	src := new(pullRequestHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertPullRequestHook(src)
	switch src.EventKey {
	case "pr:opened":
		dst.Action = scm.ActionOpen
	case "pr:declined":
		dst.Action = scm.ActionClose
	case "pr:merged":
		dst.Action = scm.ActionMerge
	default:
		return nil, nil
	}
	return dst, nil
}

//
// native data structures
//

type pushHook struct {
	EventKey   string      `json:"eventKey"`
	Date       string      `json:"date"`
	Actor      *user       `json:"actor"`
	Repository *repository `json:"repository"`
	Changes    []*change   `json:"changes"`
}

type pullRequestHook struct {
	EventKey    string       `json:"eventKey"`
	Date        string       `json:"date"`
	Actor       *user        `json:"actor"`
	PullRequest *pullRequest `json:"pullRequest"`
}

type change struct {
	Ref struct {
		ID        string `json:"id"`
		DisplayID string `json:"displayId"`
		Type      string `json:"type"`
	} `json:"ref"`
	RefID    string `json:"refId"`
	FromHash string `json:"fromHash"`
	ToHash   string `json:"toHash"`
	Type     string `json:"type"`
}

//
// push hooks
//

func convertPushHook(src *pushHook) *scm.PushHook {
	change := src.Changes[0]
	repo := convertRepository(src.Repository)
	sender := convertUser(src.Actor)
	signer := convertSignature(src.Actor)
	signer.Date, _ = time.Parse("2006-01-02T15:04:05+0000", src.Date)
	return &scm.PushHook{
		Ref: change.RefID,
		Commit: scm.Commit{
			Sha:       change.ToHash,
			Message:   "",
			Link:      "",
			Author:    signer,
			Committer: signer,
		},
		Repo:   *repo,
		Sender: *sender,
	}
}

func convertTagHook(src *pushHook) *scm.TagHook {
	change := src.Changes[0]
	sender := convertUser(src.Actor)
	repo := convertRepository(src.Repository)

	dst := &scm.TagHook{
		Ref: scm.Reference{
			Name: change.Ref.DisplayID,
			Sha:  change.ToHash,
		},
		Action: scm.ActionCreate,
		Repo:   *repo,
		Sender: *sender,
	}
	if change.Type == "DELETE" {
		dst.Action = scm.ActionDelete
		dst.Ref.Sha = change.FromHash
	}
	return dst
}

func convertBranchHook(src *pushHook) *scm.BranchHook {
	change := src.Changes[0]
	sender := convertUser(src.Actor)
	repo := convertRepository(src.Repository)

	dst := &scm.BranchHook{
		Ref: scm.Reference{
			Name: change.Ref.DisplayID,
			Sha:  change.ToHash,
		},
		Action: scm.ActionCreate,
		Repo:   *repo,
		Sender: *sender,
	}
	if change.Type == "DELETE" {
		dst.Action = scm.ActionDelete
		dst.Ref.Sha = change.FromHash
	}
	return dst
}

func convertSignature(actor *user) scm.Signature {
	return scm.Signature{
		Name:   actor.DisplayName,
		Email:  actor.EmailAddress,
		Login:  actor.Slug,
		Avatar: avatarLink(actor.EmailAddress),
	}
}

func convertPullRequestHook(src *pullRequestHook) *scm.PullRequestHook {
	repo := convertRepository(&src.PullRequest.ToRef.Repository)
	pr := convertPullRequest(src.PullRequest)
	sender := convertUser(src.Actor)

	return &scm.PullRequestHook{
		Action:      scm.ActionOpen,
		Repo:        *repo,
		PullRequest: *pr,
		Sender:      *sender,
	}
}
