// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"

	"github.com/jenkins-x/go-scm/scm"
)

type issueService struct {
	client *wrapper
}

func (s *issueService) Search(context.Context, scm.SearchOptions) ([]*scm.SearchIssue, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) AssignIssue(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	// TODO implement
	return nil, nil
}

func (s *issueService) UnassignIssue(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	// TODO implement
	return nil, nil
}

func (s *issueService) ListEvents(context.Context, string, int, scm.ListOptions) ([]*scm.ListedIssueEvent, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) ListLabels(context.Context, string, int, scm.ListOptions) ([]*scm.Label, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) AddLabel(ctx context.Context, repo string, number int, label string) (*scm.Response, error) {
	// TODO implement
	return nil, nil
}

func (s *issueService) DeleteLabel(ctx context.Context, repo string, number int, label string) (*scm.Response, error) {
	// TODO implement
	return nil, nil
}

func (s *issueService) Find(ctx context.Context, repo string, number int) (*scm.Issue, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) FindComment(ctx context.Context, repo string, index, id int) (*scm.Comment, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) List(ctx context.Context, repo string, opts scm.IssueListOptions) ([]*scm.Issue, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) ListComments(ctx context.Context, repo string, index int, opts scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) Create(ctx context.Context, repo string, input *scm.IssueInput) (*scm.Issue, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) CreateComment(ctx context.Context, repo string, number int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	// TODO implement
	return nil, nil, nil
}

func (s *issueService) DeleteComment(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	// TODO implement
	return nil, nil
}

func (s *issueService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	// TODO implement
	return nil, nil
}

func (s *issueService) Lock(ctx context.Context, repo string, number int) (*scm.Response, error) {
	// TODO implement
	return nil, nil
}

func (s *issueService) Unlock(ctx context.Context, repo string, number int) (*scm.Response, error) {
	// TODO implement
	return nil, nil
}
