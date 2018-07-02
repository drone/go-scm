// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type pullService struct {
	*issueService
}

func (s *pullService) Find(ctx context.Context, repo string, number int) (*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests/%d", repo, number)
	out := new(pullRequest)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests?%s", repo, encodePullRequestListOptions(opts))
	out := new(pullRequests)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertPullRequests(out), res, err
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests/%d/diffstat?%s", repo, number, encodeListOptions(opts))
	out := new(diffstats)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertDiffstats(out), res, err
}

func (s *pullService) Merge(ctx context.Context, repo string, number int) (*scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests/%d/merge", repo, number)
	res, err := s.client.do(ctx, "POST", path, nil, nil)
	return res, err
}

func (s *pullService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

type pullRequest struct{}

type pullRequests struct {
	pagination
	Values []*pullRequest `json:"values"`
}

func convertPullRequests(from *pullRequests) []*scm.PullRequest {
	return nil
}

func convertPullRequest(from *pullRequest) *scm.PullRequest {
	return nil
}
