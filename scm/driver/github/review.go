// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type reviewService struct {
	client *wrapper
}

func (s *reviewService) Find(ctx context.Context, repo string, number, id int) (*scm.Review, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/comments/%d", repo, id)
	out := new(review)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertReview(out), res, err
}

func (s *reviewService) List(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Review, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/comments?%s", repo, number, encodeListOptions(opts))
	out := []*review{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertReviewList(out), res, err
}

func (s *reviewService) Create(ctx context.Context, repo string, number int, input *scm.ReviewInput) (*scm.Review, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/comments", repo, number)
	in := buildReviewBody(input)
	jsonBytes, _ := json.Marshal(in)
	fmt.Printf("[go-scm] CreateReview: repo=%s number=%d endpoint=%s\n", repo, number, path)
	fmt.Printf("[go-scm] ReviewInput: Line=%d Side=%s StartLine=%d StartSide=%s SubjectType=%s InReplyTo=%d Sha=%s Path=%s\n",
		input.Line, input.Side.String(), input.StartLine, input.StartSide.String(), input.SubjectType.String(), input.InReplyTo, input.Sha, input.Path)
	fmt.Printf("[go-scm] JSON payload to GitHub: %s\n", string(jsonBytes))
	out := new(review)
	res, err := s.client.do(ctx, "POST", path, in, out)
	if err != nil {
		statusCode := 0
		if res != nil {
			statusCode = res.Status
		}
		fmt.Printf("[go-scm] CreateReview ERROR: status=%d err=%v\n", statusCode, err)
	} else {
		fmt.Printf("[go-scm] CreateReview SUCCESS: status=%d comment_id=%d url=%s\n", res.Status, out.ID, out.HTMLURL)
	}
	return convertReview(out), res, err
}

func buildReviewBody(in *scm.ReviewInput) *reviewInput {
	body := &reviewInput{
		Body:     in.Body,
		Path:     in.Path,
		CommitID: in.Sha,
	}
	switch {
	case in.InReplyTo > 0:
		body.InReplyTo = in.InReplyTo
	case in.SubjectType == scm.SubjectTypeFile:
		body.SubjectType = in.SubjectType.String()
	default:
		body.Line = in.Line
		body.Side = in.Side.String()
		if in.StartLine > 0 {
			body.StartLine = in.StartLine
			startSide := in.StartSide
			if startSide == scm.SideUnspecified {
				startSide = in.Side
			}
			body.StartSide = startSide.String()
		}
	}
	return body
}

func (s *reviewService) Delete(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/pulls/comments/%d", repo, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

type review struct {
	ID       int    `json:"id"`
	CommitID string `json:"commit_id"`
	Line     int    `json:"line"`
	Path     string `json:"path"`
	User     struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"user"`
	Body      string    `json:"body"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type reviewInput struct {
	Body        string `json:"body"`
	Path        string `json:"path,omitempty"`
	CommitID    string `json:"commit_id,omitempty"`
	Line        int    `json:"line,omitempty"`
	Side        string `json:"side,omitempty"`
	StartLine   int    `json:"start_line,omitempty"`
	StartSide   string `json:"start_side,omitempty"`
	SubjectType string `json:"subject_type,omitempty"`
	InReplyTo   int    `json:"in_reply_to,omitempty"`
}

func convertReviewList(from []*review) []*scm.Review {
	to := []*scm.Review{}
	for _, v := range from {
		to = append(to, convertReview(v))
	}
	return to
}

func convertReview(from *review) *scm.Review {
	return &scm.Review{
		ID:   from.ID,
		Body: from.Body,
		Path: from.Path,
		Line: from.Line,
		Sha:  from.CommitID,
		Link: from.HTMLURL,
		Author: scm.User{
			Login:  from.User.Login,
			Avatar: from.User.AvatarURL,
		},
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
	}
}
