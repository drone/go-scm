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
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/sha1"
)

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (interface{}, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}

	var hook interface{}
	switch req.Header.Get("X-GitHub-Event") {
	case "push":
		hook, err = s.parsePushHook(data)
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
	sig = strings.TrimPrefix(sig, "sha1=")

	if !sha1.ValidateEncoded(data, []byte(key), sig) {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePushHook(data []byte) (interface{}, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	return convertPushHook(dst), err
}

//
// native data structures
//

type (
	// gogs push webhook payload
	pushHook struct {
		Ref     string `json:"ref"`
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
)

//
// native data structure conversion
//

func convertPushHook(src *pushHook) *scm.PushHook {
	return &scm.PushHook{
		Ref: src.Ref,
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
}
