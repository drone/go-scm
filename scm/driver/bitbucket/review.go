// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type prInlineCommentInput struct {
	Content struct {
		Raw string `json:"raw"`
	} `json:"content"`
	Inline *prInline `json:"inline,omitempty"`
}

type prInline struct {
	Path      string `json:"path"`
	From      int    `json:"from,omitempty"`
	To        int    `json:"to,omitempty"`
	StartFrom int    `json:"start_from,omitempty"`
	StartTo   int    `json:"start_to,omitempty"`
}

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
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests/%d/comments", repo, number)
	in := &prInlineCommentInput{}
	in.Content.Raw = input.Body
	in.Inline = &prInline{
		Path: input.Path,
	}
	if input.Side == scm.SideLeft {
		in.Inline.From = input.Line
	} else {
		in.Inline.To = input.Line
	}
	if input.StartLine > 0 {
		startSide := input.StartSide
		if startSide == scm.SideUnspecified {
			startSide = input.Side
		}
		if startSide == scm.SideLeft {
			in.Inline.StartFrom = input.StartLine
		} else {
			in.Inline.StartTo = input.StartLine
		}
	}
	out := new(prComment)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertReviewComment(out), res, err
}

func convertReviewComment(from *prComment) *scm.Review {
	line := from.Inline.To
	if line == 0 {
		line = from.Inline.From
	}
	return &scm.Review{
		ID:   from.ID,
		Body: from.Content.Raw,
		Path: from.Inline.Path,
		Line: line,
		Link: from.Links.HTML.Href,
		Author: scm.User{
			ID:     from.User.UUID,
			Login:  from.User.Nickname,
			Name:   from.User.DisplayName,
			Avatar: from.User.Links.Avatar.Href,
		},
		Created: from.CreatedOn,
		Updated: from.UpdatedOn,
	}
}

func (s *reviewService) Delete(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}
