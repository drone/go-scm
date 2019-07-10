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
	"time"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/internal/hmac"
	"github.com/jenkins-x/go-scm/scm/driver/internal/null"
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

	guid := req.Header.Get("X-GitHub-Delivery")
	if guid == "" {
		return nil, scm.MissingHeader{"X-GitHub-Delivery"}
	}

	var hook scm.Webhook
	event := req.Header.Get("X-GitHub-Event")
	switch event {
	case "push":
		hook, err = s.parsePushHook(data, guid)
	case "create":
		hook, err = s.parseCreateHook(data)
	case "delete":
		hook, err = s.parseDeleteHook(data)
	case "pull_request":
		hook, err = s.parsePullRequestHook(data, guid)
	case "pull_request_review_comment":
		hook, err = s.parsePullRequestReviewCommentHook(data)
	case "deployment":
		hook, err = s.parseDeploymentHook(data)
	// case "issues":
	case "issue_comment":
		hook, err = s.parseIssueCommentHook(data)
	default:
		return nil, scm.UnknownWebhook{event}
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

func (s *webhookService) parsePushHook(data []byte, guid string) (*scm.PushHook, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	to := convertPushHook(dst)
	if to != nil {
		to.GUID = guid
	}
	return to, err
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

func (s *webhookService) parsePullRequestHook(data []byte, guid string) (scm.Webhook, error) {
	src := new(pullRequestHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertPullRequestHook(src)
	dst.GUID = guid
	switch src.Action {
	case "assigned":
		dst.Action = scm.ActionAssigned
	case "unassigned":
		dst.Action = scm.ActionUnassigned
	case "review_requested":
		dst.Action = scm.ActionReviewRequested
	case "review_request_removed":
		dst.Action = scm.ActionReviewRequestRemoved
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

func (s *webhookService) parsePullRequestReviewCommentHook(data []byte) (scm.Webhook, error) {
	src := new(pullRequestReviewCommentHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertPullRequestReviewCommentHook(src)
	return dst, nil
}

func (s *webhookService) parseIssueCommentHook(data []byte) (*scm.IssueCommentHook, error) {
	src := new(issueCommentHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertIssueCommentHook(src)
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

	pullRequestHookChanges struct {
		Base struct {
			Ref struct {
				From string `json:"from"`
			} `json:"ref"`
			Sha struct {
				From string `json:"from"`
			} `json:"sha"`
		} `json:"base"`
	}

	pullRequestHook struct {
		Action      string                 `json:"action"`
		Number      int                    `json:"number"`
		PullRequest pr                     `json:"pull_request"`
		Repository  repository             `json:"repository"`
		Label       label                  `json:"label"`
		Sender      user                   `json:"sender"`
		Changes     pullRequestHookChanges `json:"changes"`
	}

	label struct {
		URL         string `json:"url"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}

	pullRequestReviewCommentHook struct {
		// Action see https://developer.github.com/v3/activity/events/types/#pullrequestreviewcommentevent
		Action      string        `json:"action"`
		PullRequest pr            `json:"pull_request"`
		Repository  repository    `json:"repository"`
		Comment     reviewComment `json:"comment"`
	}

	issueCommentHook struct {
		Action     string       `json:"action"`
		Issue      issue        `json:"issue"`
		Repository repository   `json:"repository"`
		Comment    issueComment `json:"comment"`
		Sender     user         `json:"sender"`
	}

	// reviewComment describes a Pull Request review comment
	reviewComment struct {
		ID        int       `json:"id"`
		ReviewID  int       `json:"pull_request_review_id"`
		User      user      `json:"user"`
		Body      string    `json:"body"`
		Path      string    `json:"path"`
		HTMLURL   string    `json:"html_url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		// Position will be nil if the code has changed such that the comment is no
		// longer relevant.
		Position *int `json:"position"`
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
	// fix https://github.com/jenkins-x/go-scm/issues/8
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
		Label:       convertLabel(src.Label),
		PullRequest: *convertPullRequest(&src.PullRequest),
		Sender:      *convertUser(&src.Sender),
		Changes:     *convertPullRequestChanges(&src.Changes),
	}
}

func convertPullRequestChanges(src *pullRequestHookChanges) *scm.PullRequestHookChanges {
	to := &scm.PullRequestHookChanges{}
	to.Base.Sha.From = src.Base.Sha.From
	to.Base.Ref.From = src.Base.Ref.From
	return to
}

func convertLabel(src label) scm.Label {
	return scm.Label{
		Color:       src.Color,
		Description: src.Description,
		Name:        src.Name,
		URL:         src.URL,
	}
}

func convertPullRequestReviewCommentHook(src *pullRequestReviewCommentHook) *scm.PullRequestCommentHook {
	return &scm.PullRequestCommentHook{
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
		Comment:     *convertPullRequestComment(&src.Comment),
		Sender:      *convertUser(&src.Comment.User),
	}
}

/*
func convertIssueHook(dst *issueHook) *scm.IssueHook {
	return &scm.IssueHook{
		Action: convertAction(dst.Action),
		Issue:  *convertIssue(&dst.Issue),
		Repo:   *convertRepository(&dst.Repository),
		Sender: *convertUser(&dst.Sender),
	}
}
*/

func convertIssueCommentHook(dst *issueCommentHook) *scm.IssueCommentHook {
	return &scm.IssueCommentHook{
		Action:  convertAction(dst.Action),
		Issue:   *convertIssue(&dst.Issue),
		Comment: *convertIssueComment(&dst.Comment),
		Repo:    *convertRepository(&dst.Repository),
		Sender:  *convertUser(&dst.Sender),
	}
}

func convertPullRequestComment(comment *reviewComment) *scm.Comment {
	return &scm.Comment{
		ID:      comment.ID,
		Body:    comment.Body,
		Author:  *convertUser(&comment.User),
		Created: comment.CreatedAt,
		Updated: comment.UpdatedAt,
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

func convertAction(src string) (action scm.Action) {
	switch src {
	case "create", "created":
		return scm.ActionCreate
	case "delete", "deleted":
		return scm.ActionDelete
	case "update", "updated", "edit", "edited":
		return scm.ActionUpdate
	case "open", "opened":
		return scm.ActionOpen
	case "reopen", "reopened":
		return scm.ActionReopen
	case "close", "closed":
		return scm.ActionClose
	case "label", "labeled":
		return scm.ActionLabel
	case "unlabel", "unlabeled":
		return scm.ActionUnlabel
	case "merge", "merged":
		return scm.ActionMerge
	case "synchronize", "synchronized":
		return scm.ActionSync
	default:
		return
	}
}
