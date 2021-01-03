package scm

import (
	"context"
	"time"
)

type (
	CommitStatus struct {
		Status       string
		Created      time.Time
		Started      time.Time
		Name         string
		AllowFailure bool
		Author       CommitStatusAuthor
		Description  string
		Sha          string
		TargetURL    string
		Finished     time.Time
		ID           int
		Ref          string
		Coverage     float64
	}

	CommitStatusAuthor struct {
		Username  string
		State     string
		WebUrl    string
		AvatarUrl string
		ID        int
		Name      string
	}

	CommitStatusUpdateOptions struct {
		ID          string
		Sha         string
		State       string
		Ref         string
		Name        string
		TargetUrl   string
		Description string
		Coverage    float64
		PipelineID  int
	}
)

type CommitService interface {
	UpdateCommitStatus(ctx context.Context,
		repo string, sha string, options CommitStatusUpdateOptions) (*CommitStatus, *Response, error)
}
