// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"fmt"

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
	path := fmt.Sprintf("api/v1/repos/%s/pulls/%d/reviews", repo, number)

	comment := reviewCommentInput{
		Path: input.Path,
		Body: input.Body,
	}
	if input.Side == scm.SideLeft {
		comment.OldPosition = input.Line
	} else {
		comment.NewPosition = input.Line
	}

	in := &reviewCreateInput{
		Event:    "COMMENT",
		CommitID: input.Sha,
		Comments: []reviewCommentInput{comment},
	}

	out := new(reviewResponse)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertReview(out, input), res, err
}

func (s *reviewService) Delete(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

//
// native data structures
//

type reviewCreateInput struct {
	Event    string               `json:"event"`
	Body     string               `json:"body"`
	CommitID string               `json:"commit_id"`
	Comments []reviewCommentInput `json:"comments"`
}

type reviewCommentInput struct {
	Path        string `json:"path"`
	Body        string `json:"body"`
	OldPosition int    `json:"old_position,omitempty"`
	NewPosition int    `json:"new_position,omitempty"`
}

type reviewResponse struct {
	ID      int64  `json:"id"`
	HTMLURL string `json:"html_url"`
}

//
// native data structure conversion
//

func convertReview(src *reviewResponse, input *scm.ReviewInput) *scm.Review {
	return &scm.Review{
		ID:   int(src.ID),
		Body: input.Body,
		Path: input.Path,
		Line: input.Line,
		Link: src.HTMLURL,
	}
}
