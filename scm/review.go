// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"context"
	"time"
)

type Side int

const (
	SideUnspecified Side = iota
	SideRight
	SideLeft
)

func (s Side) String() string {
	switch s {
	case SideLeft:
		return "LEFT"
	default:
		return "RIGHT"
	}
}

type SubjectType int

const (
	SubjectTypeUnspecified SubjectType = iota
	SubjectTypeLine
	SubjectTypeFile
)

func (t SubjectType) String() string {
	if t == SubjectTypeFile {
		return "file"
	}
	return "line"
}

type (
	// Review represents a review comment.
	Review struct {
		ID      int
		Body    string
		Path    string
		Sha     string
		Line    int
		Link    string
		Author  User
		Created time.Time
		Updated time.Time
	}

	// ReviewInput provides the input fields required for
	// creating a review comment.
	ReviewInput struct {
		Body        string
		Sha         string
		Path        string
		Line        int         // 1-based absolute line number in the file (not diff position)
		Side        Side        // Which side of the diff (default RIGHT)
		StartLine   int         // Multi-line start; 0 means single-line
		StartSide   Side        // Side for StartLine; UNSPECIFIED mirrors Side
		SubjectType SubjectType // LINE or FILE
		InReplyTo   int         // Comment ID to reply to; 0 means root comment
	}

	// ReviewService provides access to review resources.
	ReviewService interface {
		// Find returns the review comment by id.
		Find(context.Context, string, int, int) (*Review, *Response, error)

		// List returns the review comment list.
		List(context.Context, string, int, ListOptions) ([]*Review, *Response, error)

		// Create creates a review comment.
		Create(context.Context, string, int, *ReviewInput) (*Review, *Response, error)

		// Delete deletes a review comment.
		Delete(context.Context, string, int, int) (*Response, error)
	}
)
