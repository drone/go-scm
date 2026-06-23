// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

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
	// fetch the merge request to obtain the diff refs (base, head, start shas)
	// required to anchor an inline discussion to a specific line.
	mrPath := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d", encode(repo), number)
	mr := new(pr)
	if res, err := s.client.do(ctx, "GET", mrPath, nil, mr); err != nil {
		return nil, res, err
	}

	position := discussionPosition{
		PositionType: "text",
		BaseSha:      mr.DiffRefs.BaseSha,
		HeadSha:      mr.DiffRefs.HeadSha,
		StartSha:     mr.DiffRefs.StartSha,
		OldPath:      input.Path,
		NewPath:      input.Path,
	}
	if input.Side == scm.SideLeft {
		position.OldLine = input.Line
	} else {
		position.NewLine = input.Line
	}

	in := &discussionInput{
		Body:     input.Body,
		Position: position,
	}

	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/discussions", encode(repo), number)
	out := new(discussion)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertDiscussion(out, input), res, err
}

func (s *reviewService) Delete(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

type discussionInput struct {
	Body     string             `json:"body"`
	Position discussionPosition `json:"position"`
}

type discussionPosition struct {
	PositionType string `json:"position_type"`
	BaseSha      string `json:"base_sha"`
	HeadSha      string `json:"head_sha"`
	StartSha     string `json:"start_sha"`
	OldPath      string `json:"old_path,omitempty"`
	NewPath      string `json:"new_path,omitempty"`
	OldLine      int    `json:"old_line,omitempty"`
	NewLine      int    `json:"new_line,omitempty"`
}

type discussion struct {
	ID    string `json:"id"`
	Notes []struct {
		ID     int    `json:"id"`
		Body   string `json:"body"`
		Author struct {
			Username string `json:"username"`
			Avatar   string `json:"avatar_url"`
		} `json:"author"`
	} `json:"notes"`
}

func convertDiscussion(from *discussion, input *scm.ReviewInput) *scm.Review {
	to := &scm.Review{
		Body: input.Body,
		Path: input.Path,
		Line: input.Line,
		Sha:  input.Sha,
	}
	if len(from.Notes) > 0 {
		note := from.Notes[0]
		to.ID = note.ID
		if note.Body != "" {
			to.Body = note.Body
		}
		to.Author = scm.User{
			Login:  note.Author.Username,
			Avatar: note.Author.Avatar,
		}
	}
	return to
}
