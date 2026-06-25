// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type reviewService struct {
	client *wrapper
}

type reviewAnchor struct {
	DiffType string `json:"diffType"`
	Line     int    `json:"line"`
	LineType string `json:"lineType"`
	FileType string `json:"fileType"`
	Path     string `json:"path"`
	SrcPath  string `json:"srcPath,omitempty"`
}

type pullRequestReviewCommentInput struct {
	Text   string       `json:"text"`
	Anchor reviewAnchor `json:"anchor"`
}

func (s *reviewService) Find(ctx context.Context, repo string, number, id int) (*scm.Review, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *reviewService) List(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Review, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *reviewService) Create(ctx context.Context, repo string, number int, input *scm.ReviewInput) (*scm.Review, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/comments", namespace, name, number)
	in := &pullRequestReviewCommentInput{
		Text: input.Body,
		Anchor: reviewAnchor{
			DiffType: "EFFECTIVE",
			Line:     input.Line,
			LineType: stashLineType(input.Side),
			FileType: stashFileType(input.Side),
			Path:     input.Path,
		},
	}
	out := new(pullRequestComment)
	res, err := s.client.do(ctx, "POST", path, in, out)
	if err != nil {
		return nil, res, err
	}
	review := convertReviewComment(out, input)
	review.Link = s.commentLink(namespace, name, number, out.ID)
	return review, res, err
}

func (s *reviewService) Delete(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *reviewService) commentLink(namespace, name string, number, commentID int) string {
	if s.client.BaseURL == nil {
		return ""
	}
	return fmt.Sprintf("%s://%s/projects/%s/repos/%s/pull-requests/%d/overview?commentId=%d",
		s.client.BaseURL.Scheme, s.client.BaseURL.Host, namespace, name, number, commentID)
}

func stashLineType(side scm.Side) string {
	if side == scm.SideLeft {
		return "REMOVED"
	}
	return "ADDED"
}

func stashFileType(side scm.Side) string {
	if side == scm.SideLeft {
		return "FROM"
	}
	return "TO"
}

func convertReviewComment(from *pullRequestComment, input *scm.ReviewInput) *scm.Review {
	comment := convertPullRequestComment(from)
	return &scm.Review{
		ID:      comment.ID,
		Body:    comment.Body,
		Path:    input.Path,
		Line:    input.Line,
		Sha:     input.Sha,
		Author:  comment.Author,
		Created: comment.Created,
		Updated: comment.Updated,
	}
}
