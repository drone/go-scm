// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(ctx context.Context, repo string, index int) (*scm.PullRequest, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d", harnessURI, index)
	out := new(pr)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err

}

func (s *pullService) FindComment(context.Context, string, int, int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListComments(context.Context, string, int, scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListCommits(ctx context.Context, repo string, index int, opts scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d/commits?%s", harnessURI, index, encodeListOptions(opts))
	out := []*commit{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommits(out), res, err
}

func (s *pullService) ListChanges(context.Context, string, int, scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) CreateComment(context.Context, string, int, *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Merge(ctx context.Context, repo string, index int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Close(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

// native data structures
type (
	pr struct {
		Author struct {
			Created     int    `json:"created"`
			DisplayName string `json:"display_name"`
			Email       string `json:"email"`
			ID          int    `json:"id"`
			Type        string `json:"type"`
			UID         string `json:"uid"`
			Updated     int    `json:"updated"`
		} `json:"author"`
		Created       int    `json:"created"`
		Description   string `json:"description"`
		Edited        int    `json:"edited"`
		IsDraft       bool   `json:"is_draft"`
		MergeBaseSha  string `json:"merge_base_sha"`
		MergeHeadSha  string `json:"merge_head_sha"`
		MergeStrategy string `json:"merge_strategy"`
		Merged        int    `json:"merged"`
		Merger        struct {
			Created     int    `json:"created"`
			DisplayName string `json:"display_name"`
			Email       string `json:"email"`
			ID          int    `json:"id"`
			Type        string `json:"type"`
			UID         string `json:"uid"`
			Updated     int    `json:"updated"`
		} `json:"merger"`
		Number       int    `json:"number"`
		SourceBranch string `json:"source_branch"`
		SourceRepoID int    `json:"source_repo_id"`
		State        string `json:"state"`
		Stats        struct {
			Commits       int `json:"commits"`
			Conversations int `json:"conversations"`
			FilesChanged  int `json:"files_changed"`
		} `json:"stats"`
		TargetBranch string `json:"target_branch"`
		TargetRepoID int    `json:"target_repo_id"`
		Title        string `json:"title"`
	}

	reference struct {
		Repo repository `json:"repo"`
		Name string     `json:"ref"`
		Sha  string     `json:"sha"`
	}

	prInput struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Head  string `json:"head"`
		Base  string `json:"base"`
	}
	commit struct {
		Author struct {
			Identity struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"committer"`
		Message string `json:"message"`
		Sha     string `json:"sha"`
		Title   string `json:"title"`
	}
)

// native data structure conversion
func convertPullRequests(src []*pr) []*scm.PullRequest {
	dst := []*scm.PullRequest{}
	for _, v := range src {
		dst = append(dst, convertPullRequest(v))
	}
	return dst
}

func convertPullRequest(src *pr) *scm.PullRequest {
	return &scm.PullRequest{
		Number: src.Number,
		Title:  src.Title,
		Body:   src.Description,
		Source: src.SourceBranch,
		Target: src.TargetBranch,
		Merged: src.Merged != 0,
		Author: scm.User{
			Login: src.Author.Email,
			Name:  src.Author.DisplayName,
			ID:    src.Author.UID,
			Email: src.Author.Email,
		},
		Fork:   "fork",
		Ref:    fmt.Sprintf("refs/pull/%d/head", src.Number),
		Closed: src.State == "closed",
	}
}

func convertCommits(src []*commit) []*scm.Commit {
	dst := []*scm.Commit{}
	for _, v := range src {
		dst = append(dst, convertCommit(v))
	}
	return dst
}

func convertCommit(src *commit) *scm.Commit {
	return &scm.Commit{
		Message: src.Message,
		Sha:     src.Sha,
		Author: scm.Signature{
			Name:  src.Author.Identity.Name,
			Email: src.Author.Identity.Email,
		},
		Committer: scm.Signature{
			Name:  src.Committer.Identity.Name,
			Email: src.Committer.Identity.Email,
		},
	}
}
