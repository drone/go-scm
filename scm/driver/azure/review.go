// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"fmt"
	"strings"

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
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-threads/create?view=azure-devops-rest-6.0
	endpoint := fmt.Sprintf("%s/%s/_apis/git/repositories/%s/pullRequests/%d/threads?api-version=6.0",
		s.client.owner, s.client.project, repo, number)
	in := buildThreadInput(input)
	out := new(thread)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertReview(out, input), res, err
}

func (s *reviewService) Delete(ctx context.Context, repo string, number, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func buildThreadInput(input *scm.ReviewInput) *threadInput {
	in := &threadInput{
		Status: 1,
		Comments: []threadCommentInput{
			{
				ParentCommentID: 0,
				Content:         input.Body,
				CommentType:     1,
			},
		},
	}

	path := input.Path
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	endLine := input.Line
	startLine := input.Line
	if input.StartLine > 0 {
		startLine = input.StartLine
	}

	tc := &threadContext{
		FilePath: path,
	}
	start := &commentPosition{Line: startLine, Offset: 1}
	end := &commentPosition{Line: endLine, Offset: 1}
	if input.Side == scm.SideLeft {
		tc.LeftFileStart = start
		tc.LeftFileEnd = end
	} else {
		tc.RightFileStart = start
		tc.RightFileEnd = end
	}
	in.ThreadContext = tc
	return in
}

type commentPosition struct {
	Line   int `json:"line"`
	Offset int `json:"offset"`
}

type threadContext struct {
	FilePath       string           `json:"filePath,omitempty"`
	LeftFileStart  *commentPosition `json:"leftFileStart,omitempty"`
	LeftFileEnd    *commentPosition `json:"leftFileEnd,omitempty"`
	RightFileStart *commentPosition `json:"rightFileStart,omitempty"`
	RightFileEnd   *commentPosition `json:"rightFileEnd,omitempty"`
}

type threadCommentInput struct {
	ParentCommentID int    `json:"parentCommentId"`
	Content         string `json:"content"`
	CommentType     int    `json:"commentType"`
}

type threadInput struct {
	Comments      []threadCommentInput `json:"comments"`
	Status        int                  `json:"status"`
	ThreadContext *threadContext       `json:"threadContext,omitempty"`
}

type thread struct {
	ID       int `json:"id"`
	Comments []struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
	} `json:"comments"`
	ThreadContext struct {
		FilePath string `json:"filePath"`
	} `json:"threadContext"`
}

func convertReview(from *thread, input *scm.ReviewInput) *scm.Review {
	path := from.ThreadContext.FilePath
	if path == "" {
		path = input.Path
	}
	return &scm.Review{
		ID:   from.ID,
		Body: input.Body,
		Path: path,
		Line: input.Line,
	}
}
