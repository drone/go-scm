package scm

import (
	"context"
)

type (
	InstallationService interface {
		ListRepositories(ctx context.Context, opts ListOptions) ([]*Repository, *Response, error)
	}
)
