// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"fmt"
	"time"

	"github.com/livecycle/go-scm/scm"
)

type issueService struct {
	client *wrapper
}

func (s *issueService) Find(ctx context.Context, repo string, number int) (*scm.Issue, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues/%d", repo, number)
	out := new(issue)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertIssue(out), res, err
}

func (s *issueService) FindComment(ctx context.Context, repo string, index, id int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) List(ctx context.Context, repo string, opts scm.IssueListOptions) ([]*scm.Issue, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues?%s", repo, encodeIssueListOptions(opts))
	out := []*issue{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertIssueList(out), res, err
}

func (s *issueService) ListComments(ctx context.Context, repo string, index int, opts scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues/%d/comments?%s", repo, index, encodeListOptions(opts))
	out := []*issueComment{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertIssueCommentList(out), res, err
}

func (s *issueService) Create(ctx context.Context, repo string, input *scm.IssueInput) (*scm.Issue, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues", repo)
	in := &issueInput{
		Title: input.Title,
		Body:  input.Body,
	}
	out := new(issue)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertIssue(out), res, err
}

func (s *issueService) CreateComment(ctx context.Context, repo string, index int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues/%d/comments", repo, index)
	in := &issueCommentInput{
		Body: input.Body,
	}
	out := new(issueComment)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertIssueComment(out), res, err
}

func (s *issueService) DeleteComment(ctx context.Context, repo string, index, id int) (*scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/issues/%d/comments/%d", repo, index, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *issueService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *issueService) Lock(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *issueService) Unlock(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

//
// native data structures
//

type (
	// gitea issue response object.
	issue struct {
		ID          int       `json:"id"`
		Number      int       `json:"number"`
		User        user      `json:"user"`
		Title       string    `json:"title"`
		Body        string    `json:"body"`
		State       string    `json:"state"`
		Labels      []string  `json:"labels"`
		Comments    int       `json:"comments"`
		Created     time.Time `json:"created_at"`
		Updated     time.Time `json:"updated_at"`
		PullRequest *struct {
			Merged   bool        `json:"merged"`
			MergedAt interface{} `json:"merged_at"`
		} `json:"pull_request"`
	}

	// gitea issue request object.
	issueInput struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	// gitea issue comment response object.
	issueComment struct {
		ID        int       `json:"id"`
		HTMLURL   string    `json:"html_url"`
		User      user      `json:"user"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// gitea issue comment request object.
	issueCommentInput struct {
		Body string `json:"body"`
	}
)

//
// native data structure conversion
//

func convertIssueList(from []*issue) []*scm.Issue {
	to := []*scm.Issue{}
	for _, v := range from {
		to = append(to, convertIssue(v))
	}
	return to
}

func convertIssue(from *issue) *scm.Issue {
	return &scm.Issue{
		Number:  from.Number,
		Title:   from.Title,
		Body:    from.Body,
		Link:    "", // TODO construct the link to the issue.
		Closed:  from.State == "closed",
		Author:  *convertUser(&from.User),
		Created: from.Created,
		Updated: from.Updated,
	}
}

func convertIssueCommentList(from []*issueComment) []*scm.Comment {
	to := []*scm.Comment{}
	for _, v := range from {
		to = append(to, convertIssueComment(v))
	}
	return to
}

func convertIssueComment(from *issueComment) *scm.Comment {
	return &scm.Comment{
		ID:      from.ID,
		Body:    from.Body,
		Author:  *convertUser(&from.User),
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
	}
}
