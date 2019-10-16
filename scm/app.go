package scm

import (
	"context"
	"time"
)

type (
	InstallationToken struct {
		Token     string
		ExpiresAt *time.Time
	}

	// AppService for GitHub App support
	AppService interface {
		CreateInstallationToken(ctx context.Context, id int64) (*InstallationToken, *Response, error)
	}
)
