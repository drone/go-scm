// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/drone/go-scm/scm"
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
	// switch req.Header.Get("SOMETHING") {
	// case "pull_request":
	hook, err = s.parsePullRequest(data)
	// default:
	// 	return nil, scm.ErrUnknownEvent
	// }
	return hook, err
}

func (s *webhookService) parsePullRequest(data []byte) (scm.Webhook, error) {
	src := new(pullRequestHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}
	dst := convertPullRequestHook(src)
	switch src.EventType {
	case "git.pullrequest.created":
		dst.Action = scm.ActionOpen
	case "git.pullrequest.updated":
		dst.Action = scm.ActionUpdate
	default:
		dst.Action = scm.ActionUnknown
	}

	return dst, nil
}

func convertPullRequestHook(src *pullRequestHook) (returnVal *scm.PullRequestHook) {
	returnVal = &scm.PullRequestHook{
		PullRequest: scm.PullRequest{
			Number: src.Resource.PullRequestID,
			Title:  src.Resource.Title,
			Body:   src.Resource.Description,
			Sha:    src.Resource.MergeID,
			Ref:    src.Resource.SourceRefName,
			Source: src.Resource.SourceRefName,
			Target: src.Resource.TargetRefName,
			Link:   src.Resource.URL,
			Author: scm.User{
				Login:  src.Resource.CreatedBy.ID,
				Name:   src.Resource.CreatedBy.DisplayName,
				Email:  src.Resource.CreatedBy.UniqueName,
				Avatar: src.Resource.CreatedBy.ImageURL,
			},
			Created: src.Resource.CreationDate,
		},
		Repo: scm.Repository{
			ID:        src.Resource.Repository.ID,
			Name:      src.Resource.Repository.ID,
			Namespace: src.Resource.Repository.Name,
			Link:      src.Resource.Repository.URL,
		},
		Sender: scm.User{
			Login:  src.Resource.CreatedBy.ID,
			Name:   src.Resource.CreatedBy.DisplayName,
			Email:  src.Resource.CreatedBy.UniqueName,
			Avatar: src.Resource.CreatedBy.ImageURL,
		},
	}
	return returnVal
}

type pullRequestHook struct {
	ID          string `json:"id"`
	EventType   string `json:"eventType"`
	PublisherID string `json:"publisherId"`
	Scope       string `json:"scope"`
	Message     struct {
		Text     string `json:"text"`
		HTML     string `json:"html"`
		Markdown string `json:"markdown"`
	} `json:"message"`
	DetailedMessage struct {
		Text     string `json:"text"`
		HTML     string `json:"html"`
		Markdown string `json:"markdown"`
	} `json:"detailedMessage"`
	Resource struct {
		Repository struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			URL     string `json:"url"`
			Project struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				URL   string `json:"url"`
				State string `json:"state"`
			} `json:"project"`
			DefaultBranch string `json:"defaultBranch"`
			RemoteURL     string `json:"remoteUrl"`
		} `json:"repository"`
		PullRequestID int    `json:"pullRequestId"`
		Status        string `json:"status"`
		CreatedBy     struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
			UniqueName  string `json:"uniqueName"`
			URL         string `json:"url"`
			ImageURL    string `json:"imageUrl"`
		} `json:"createdBy"`
		CreationDate          time.Time `json:"creationDate"`
		Title                 string    `json:"title"`
		Description           string    `json:"description"`
		SourceRefName         string    `json:"sourceRefName"`
		TargetRefName         string    `json:"targetRefName"`
		MergeStatus           string    `json:"mergeStatus"`
		MergeID               string    `json:"mergeId"`
		LastMergeSourceCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeSourceCommit"`
		LastMergeTargetCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeTargetCommit"`
		LastMergeCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeCommit"`
		Reviewers []struct {
			ReviewerURL interface{} `json:"reviewerUrl"`
			Vote        int         `json:"vote"`
			ID          string      `json:"id"`
			DisplayName string      `json:"displayName"`
			UniqueName  string      `json:"uniqueName"`
			URL         string      `json:"url"`
			ImageURL    string      `json:"imageUrl"`
			IsContainer bool        `json:"isContainer"`
		} `json:"reviewers"`
		URL string `json:"url"`
	} `json:"resource"`
	ResourceVersion    string `json:"resourceVersion"`
	ResourceContainers struct {
		Collection struct {
			ID string `json:"id"`
		} `json:"collection"`
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
	} `json:"resourceContainers"`
	CreatedDate time.Time `json:"createdDate"`
}
