package scm

import (
	"context"
)

// Repository represents a git repository.

type (
	Installation struct {
		ID                        string
		Organization              *Organization
		AppName                   string
		CurrentUserHasAdminAccess bool
		CurrentUserHasPushAccess  bool
	}
	InstallationService interface {
		ListRepositories(ctx context.Context, opts ListOptions) ([]*Repository, *Response, error)
		ListInstallationsForAuthenticatedUser(ctx context.Context, opts ListOptions) ([]*Installation, *Response, error)
		ListInstallationRepositoriesForAuthenticatedUser(ctx context.Context, installationId int, opts ListOptions) ([]*Repository, *Response, error)
	}
)
