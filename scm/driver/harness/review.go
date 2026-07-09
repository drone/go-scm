// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
)

type reviewService struct {
	client *wrapper
}

func (s *reviewService) Find(ctx context.Context, repo string, number, id int) (*scm.Review, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *reviewService) List(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Review, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *reviewService) Create(ctx context.Context, repo string, number int, input *scm.ReviewInput) (*scm.Review, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d/comments?%s", repoId, number, queryParams)

	in := &prComment{
		Text:            input.Body,
		Path:            input.Path,
		SourceCommitSha: input.Sha,
		ParentID:        input.InReplyTo,
		LineEnd:         input.Line,
		LineEndNew:      input.Side != scm.SideLeft,
	}
	if input.StartLine > 0 {
		startSide := input.StartSide
		if startSide == scm.SideUnspecified {
			startSide = input.Side
		}
		in.LineStart = input.StartLine
		in.LineStartNew = startSide != scm.SideLeft
	} else {
		in.LineStart = input.Line
		in.LineStartNew = in.LineEndNew
	}

	out := new(prCommentResponse)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertReviewComment(out, input), res, err
}

func (s *reviewService) Delete(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func convertReviewComment(from *prCommentResponse, input *scm.ReviewInput) *scm.Review {
	return &scm.Review{
		ID:   from.Id,
		Body: from.Text,
		Path: input.Path,
		Sha:  input.Sha,
		Line: input.Line,
		Author: scm.User{
			ID:      strconv.Itoa(from.Author.Id),
			Login:   from.Author.Uid,
			Name:    from.Author.DisplayName,
			Email:   from.Author.Email,
			Created: time.UnixMilli(from.Author.Created),
			Updated: time.UnixMilli(from.Author.Updated),
		},
		Created: time.UnixMilli(from.Created),
		Updated: time.UnixMilli(from.Updated),
	}
}
