// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/hmac"
	"github.com/drone/go-scm/scm/driver/internal/null"
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
	switch req.Header.Get("X-GitHub-Event") {
	case "push":
		hook, err = s.parsePushHook(data)
	case "create":
		hook, err = s.parseCreateHook(data)
	case "delete":
		hook, err = s.parseDeleteHook(data)
	case "pull_request":
		hook, err = s.parsePullRequestHook(data)
	case "deployment":
		hook, err = s.parseDeploymentHook(data)
	// case "pull_request_review_comment":
	// case "issues":
	// case "issue_comment":
	default:
		return nil, scm.ErrUnknownEvent
	}
	if err != nil {
		return nil, err
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
	return convertPushHook(dst), err
}

func (s *webhookService) parseCreateHook(data []byte) (scm.Webhook, error) {
	src := new(createDeleteHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	if src.RefType == "branch" {
		dst := convertBranchHook(src)
		dst.Action = scm.ActionCreate
		return dst, nil
	}
	dst := convertTagHook(src)
	dst.Action = scm.ActionCreate
	return dst, nil
}

func (s *webhookService) parseDeleteHook(data []byte) (scm.Webhook, error) {
	src := new(createDeleteHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	if src.RefType == "branch" {
		dst := convertBranchHook(src)
		dst.Action = scm.ActionDelete
		return dst, nil
	}
	dst := convertTagHook(src)
	dst.Action = scm.ActionDelete
	return dst, nil
}

func (s *webhookService) parseDeploymentHook(data []byte) (scm.Webhook, error) {
	src := new(deploymentHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertDeploymentHook(src)
	return dst, nil
}

func (s *webhookService) parsePullRequestHook(data []byte) (scm.Webhook, error) {
	src := new(pullRequestHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertPullRequestHook(src)
	switch src.Action {
	case "assigned", "unassigned", "review_requested", "review_request_removed":
		return nil, nil
	case "labeled":
		dst.Action = scm.ActionLabel
	case "unlabeled":
		dst.Action = scm.ActionUnlabel
	case "opened":
		dst.Action = scm.ActionOpen
	case "edited":
		dst.Action = scm.ActionUpdate
	case "closed":
		// if merged == true
		//    dst.Action = scm.ActionMerge
		dst.Action = scm.ActionClose
	case "reopened":
		dst.Action = scm.ActionReopen
	case "synchronize":
		dst.Action = scm.ActionSync
	}
	return dst, nil
}

//
// native data structures
//

type (
	// github create webhook payload
	createDeleteHook struct {
		Ref        string     `json:"ref"`
		RefType    string     `json:"ref_type"`
		Repository repository `json:"repository"`
		Sender     user       `json:"sender"`
	}

	// github push webhook payload
	pushHook struct {
		Ref     string `json:"ref"`
		BaseRef string `json:"base_ref"`
		Before  string `json:"before"`
		After   string `json:"after"`
		Compare string `json:"compare"`
		Head    struct {
			ID        string `json:"id"`
			TreeID    string `json:"tree_id"`
			Distinct  bool   `json:"distinct"`
			Message   string `json:"message"`
			Timestamp string `json:"timestamp"`
			URL       string `json:"url"`
			Author    struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"author"`
			Committer struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"committer"`
			Added    []interface{} `json:"added"`
			Removed  []interface{} `json:"removed"`
			Modified []string      `json:"modified"`
		} `json:"head_commit"`
		Commits []struct {
			ID        string `json:"id"`
			TreeID    string `json:"tree_id"`
			Distinct  bool   `json:"distinct"`
			Message   string `json:"message"`
			Timestamp string `json:"timestamp"`
			URL       string `json:"url"`
			Author    struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"author"`
			Committer struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"committer"`
			Added    []interface{} `json:"added"`
			Removed  []interface{} `json:"removed"`
			Modified []string      `json:"modified"`
		} `json:"commits"`
		Repository struct {
			ID    int64 `json:"id"`
			Owner struct {
				Login     string `json:"login"`
				AvatarURL string `json:"avatar_url"`
			} `json:"owner"`
			Name          string `json:"name"`
			FullName      string `json:"full_name"`
			Private       bool   `json:"private"`
			Fork          bool   `json:"fork"`
			HTMLURL       string `json:"html_url"`
			SSHURL        string `json:"ssh_url"`
			CloneURL      string `json:"clone_url"`
			DefaultBranch string `json:"default_branch"`
		} `json:"repository"`
		Pusher user `json:"pusher"`
		Sender user `json:"sender"`
	}

	pullRequestHook struct {
		Action      string     `json:"action"`
		Number      int        `json:"number"`
		PullRequest pr         `json:"pull_request"`
		Repository  repository `json:"repository"`
		Sender      user       `json:"sender"`
	}

	// github deployment webhook payload
	deploymentHook struct {
		Deployment struct {
			Creator        user        `json:"creator"`
			Description    null.String `json:"description"`
			Environment    null.String `json:"environment"`
			EnvironmentURL null.String `json:"environment_url"`
			Sha            null.String `json:"sha"`
			Ref            null.String `json:"ref"`
			Task           null.String `json:"task"`
			Payload        interface{} `json:"payload"`
		} `json:"deployment"`
		Repository repository `json:"repository"`
		Sender     user       `json:"sender"`
	}
)

//
// native data structure conversion
//

func convertPushHook(src *pushHook) *scm.PushHook {
	dst := &scm.PushHook{
		Ref:     src.Ref,
		BaseRef: src.BaseRef,
		Before:  src.Before,
		After:   src.After,
		Commit: scm.Commit{
			Sha:     src.After,
			Message: src.Head.Message,
			Link:    src.Compare,
			Author: scm.Signature{
				Login: src.Head.Author.Username,
				Email: src.Head.Author.Email,
				Name:  src.Head.Author.Name,
				// TODO (bradrydzewski) set the timestamp
			},
			Committer: scm.Signature{
				Login: src.Head.Committer.Username,
				Email: src.Head.Committer.Email,
				Name:  src.Head.Committer.Name,
				// TODO (bradrydzewski) set the timestamp
			},
		},
		Repo: scm.Repository{
			ID:        fmt.Sprint(src.Repository.ID),
			Namespace: src.Repository.Owner.Login,
			Name:      src.Repository.Name,
			Branch:    src.Repository.DefaultBranch,
			Private:   src.Repository.Private,
			Clone:     src.Repository.CloneURL,
			CloneSSH:  src.Repository.SSHURL,
			Link:      src.Repository.HTMLURL,
		},
		Sender: *convertUser(&src.Sender),
	}
	// fix https://github.com/drone/go-scm/issues/8
	if scm.IsTag(dst.Ref) && src.Head.ID != "" {
		dst.Commit.Sha = src.Head.ID
		dst.After = src.Head.ID
	}
	return dst
}

func convertBranchHook(src *createDeleteHook) *scm.BranchHook {
	return &scm.BranchHook{
		Ref: scm.Reference{
			Name: src.Ref,
		},
		Repo: scm.Repository{
			ID:        fmt.Sprint(src.Repository.ID),
			Namespace: src.Repository.Owner.Login,
			Name:      src.Repository.Name,
			Branch:    src.Repository.DefaultBranch,
			Private:   src.Repository.Private,
			Clone:     src.Repository.CloneURL,
			CloneSSH:  src.Repository.SSHURL,
			Link:      src.Repository.HTMLURL,
		},
		Sender: *convertUser(&src.Sender),
	}
}

func convertTagHook(src *createDeleteHook) *scm.TagHook {
	return &scm.TagHook{
		Ref: scm.Reference{
			Name: src.Ref,
		},
		Repo: scm.Repository{
			ID:        fmt.Sprint(src.Repository.ID),
			Namespace: src.Repository.Owner.Login,
			Name:      src.Repository.Name,
			Branch:    src.Repository.DefaultBranch,
			Private:   src.Repository.Private,
			Clone:     src.Repository.CloneURL,
			CloneSSH:  src.Repository.SSHURL,
			Link:      src.Repository.HTMLURL,
		},
		Sender: *convertUser(&src.Sender),
	}
}

func convertPullRequestHook(src *pullRequestHook) *scm.PullRequestHook {
	return &scm.PullRequestHook{
		// Action        Action
		Repo: scm.Repository{
			ID:        fmt.Sprint(src.Repository.ID),
			Namespace: src.Repository.Owner.Login,
			Name:      src.Repository.Name,
			Branch:    src.Repository.DefaultBranch,
			Private:   src.Repository.Private,
			Clone:     src.Repository.CloneURL,
			CloneSSH:  src.Repository.SSHURL,
			Link:      src.Repository.HTMLURL,
		},
		PullRequest: *convertPullRequest(&src.PullRequest),
		Sender:      *convertUser(&src.Sender),
	}
}

func convertDeploymentHook(src *deploymentHook) *scm.DeployHook {
	dst := &scm.DeployHook{
		Data: src.Deployment.Payload,
		Desc: src.Deployment.Description.String,
		Ref: scm.Reference{
			Name: src.Deployment.Ref.String,
			Path: src.Deployment.Ref.String,
			Sha:  src.Deployment.Sha.String,
		},
		Repo: scm.Repository{
			ID:        fmt.Sprint(src.Repository.ID),
			Namespace: src.Repository.Owner.Login,
			Name:      src.Repository.Name,
			Branch:    src.Repository.DefaultBranch,
			Private:   src.Repository.Private,
			Clone:     src.Repository.CloneURL,
			CloneSSH:  src.Repository.SSHURL,
			Link:      src.Repository.HTMLURL,
		},
		Sender:    *convertUser(&src.Sender),
		Task:      src.Deployment.Task.String,
		Target:    src.Deployment.Environment.String,
		TargetURL: src.Deployment.EnvironmentURL.String,
	}
	if tagRE.MatchString(dst.Ref.Name) {
		dst.Ref.Path = scm.ExpandRef(dst.Ref.Path, "refs/tags/")
	} else {
		dst.Ref.Path = scm.ExpandRef(dst.Ref.Path, "refs/heads/")
	}
	return dst
}

// regexp help determine if the named git object is a tag.
// this is not meant to be 100% accurate.
var tagRE = regexp.MustCompile("^v?(\\d+).(.+)")
