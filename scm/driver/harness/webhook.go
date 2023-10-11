// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"crypto/sha256"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/hmac"
)

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
	switch req.Header.Get("X-Harness-Trigger") {
	// case "create":
	// 	hook, err = s.parseCreateHook(data)
	// case "delete":
	// 	hook, err = s.parseDeleteHook(data)
	// case "issues":
	// 	hook, err = s.parseIssueHook(data)
	case "branch_created":
		hook, err = s.parsePushHook(data)
	case "pullreq_created", "pullreq_reopened", "pullreq_branch_updated":
		hook, err = s.parsePullRequestHook(data)
	default:
		return nil, scm.ErrUnknownEvent
	}
	if err != nil {
		return nil, err
	}

	// get the gitea signature key to verify the payload
	// signature. If no key is provided, no validation
	// is performed.
	key, err := fn(hook)
	if err != nil {
		return hook, err
	} else if key == "" {
		return hook, nil
	}

	secret := req.FormValue("secret")
	signature := req.Header.Get("X-Harness-Signature")

	// fail if no signature passed
	if signature == "" && secret == "" {
		return hook, scm.ErrSignatureInvalid
	}

	// test signature if header not set and secret is in payload
	if signature == "" && secret != "" && secret != key {
		return hook, scm.ErrSignatureInvalid
	}

	// test signature using header
	if signature != "" && !hmac.Validate(sha256.New, data, []byte(key), signature) {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePullRequestHook(data []byte) (scm.Webhook, error) {
	dst := new(pullRequestHook)
	err := json.Unmarshal(data, dst)
	return convertPullRequestHook(dst), err
}

func (s *webhookService) parsePushHook(data []byte) (scm.Webhook, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	return convertPushHook(dst), err
}

// native data structures
type (
	repo struct {
		ID            int    `json:"id"`
		Path          string `json:"path"`
		UID           string `json:"uid"`
		DefaultBranch string `json:"default_branch"`
		GitURL        string `json:"git_url"`
	}
	principal struct {
		ID          int    `json:"id"`
		UID         string `json:"uid"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Type        string `json:"type"`
		Created     int64  `json:"created"`
		Updated     int64  `json:"updated"`
	}
	pullReq struct {
		Number        int         `json:"number"`
		State         string      `json:"state"`
		IsDraft       bool        `json:"is_draft"`
		Title         string      `json:"title"`
		SourceRepoID  int         `json:"source_repo_id"`
		SourceBranch  string      `json:"source_branch"`
		TargetRepoID  int         `json:"target_repo_id"`
		TargetBranch  string      `json:"target_branch"`
		MergeStrategy interface{} `json:"merge_strategy"`
	}
	targetRef struct {
		Name string `json:"name"`
		Repo struct {
			ID            int    `json:"id"`
			Path          string `json:"path"`
			UID           string `json:"uid"`
			DefaultBranch string `json:"default_branch"`
			GitURL        string `json:"git_url"`
		} `json:"repo"`
	}
	ref struct {
		Name string `json:"name"`
		Repo struct {
			ID            int    `json:"id"`
			Path          string `json:"path"`
			UID           string `json:"uid"`
			DefaultBranch string `json:"default_branch"`
			GitURL        string `json:"git_url"`
		} `json:"repo"`
	}
	hookCommit struct {
		Sha     string `json:"sha"`
		Message string `json:"message"`
		Author  struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When string `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"identity"`
			When string `json:"when"`
		} `json:"committer"`
	}
	// harness pull request webhook payload
	pullRequestHook struct {
		Trigger   string     `json:"trigger"`
		Repo      repo       `json:"repo"`
		Principal principal  `json:"principal"`
		PullReq   pullReq    `json:"pull_req"`
		TargetRef targetRef  `json:"target_ref"`
		Ref       ref        `json:"ref"`
		Sha       string     `json:"sha"`
		Commit    hookCommit `json:"commit"`
	}
	// harness push webhook payload
	pushHook struct {
		Trigger   string     `json:"trigger"`
		Repo      repo       `json:"repo"`
		Principal principal  `json:"principal"`
		Ref       ref        `json:"ref"`
		Commit    hookCommit `json:"commit"`
		Sha       string     `json:"sha"`
		OldSha    string     `json:"old_sha"`
		Forced    bool       `json:"forced"`
	}
)

//
// native data structure conversion
//

func convertPullRequestHook(dst *pullRequestHook) *scm.PullRequestHook {
	return &scm.PullRequestHook{
		Action: convertAction(dst.Trigger),
		PullRequest: scm.PullRequest{
			Number: dst.PullReq.Number,
			Title:  dst.PullReq.Title,
			Closed: dst.PullReq.State != "open",
			Source: dst.PullReq.SourceBranch,
			Target: dst.PullReq.TargetBranch,
			Fork:   "fork",
			Link:   dst.Ref.Repo.GitURL,
			Sha:    dst.Commit.Sha,
			Ref:    dst.Ref.Name,
			Author: scm.User{
				Name:  dst.Commit.Committer.Identity.Name,
				Email: dst.Commit.Committer.Identity.Email,
			},
		},
		Repo: scm.Repository{
			ID:     dst.Repo.UID,
			Branch: dst.Repo.DefaultBranch,
			Link:   dst.Repo.GitURL,
		},
		Sender: scm.User{
			Email: dst.Principal.Email,
		},
	}
}

func convertPushHook(dst *pushHook) *scm.PushHook {
	return &scm.PushHook{
		Ref:    dst.Sha,
		Before: dst.OldSha,
		After:  dst.Sha,
		Repo: scm.Repository{
			Name: dst.Repo.UID,
		},
		Commit: scm.Commit{
			Sha:     dst.Commit.Sha,
			Message: dst.Commit.Message,
			Author: scm.Signature{
				Name:  dst.Commit.Author.Identity.Name,
				Email: dst.Commit.Author.Identity.Email,
			},
		},
		Sender: scm.User{
			Name: dst.Principal.DisplayName,
		},
	}
}

func convertAction(src string) (action scm.Action) {
	switch src {
	case "pullreq_created":
		return scm.ActionCreate
	case "pullreq_branch_updated":
		return scm.ActionUpdate
	case "pullreq_reopened":
		return scm.ActionReopen
	default:
		return
	}
}
