// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"

	"github.com/drone/go-scm/scm"
)

// RepositoryService implements the repository service for
// the GitHub driver.
type RepositoryService struct {
	client *wrapper
}

// Find returns the repository by name.
func (s *RepositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// FindHook returns a repository hook.
func (s *RepositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// FindPerms returns the repository permissions.
func (s *RepositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// List returns the user repository list.
func (s *RepositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// ListHooks returns a list or repository hooks.
func (s *RepositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// ListStatus returns a list of commit statuses.
func (s *RepositoryService) ListStatus(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// CreateHook creates a new repository webhook.
func (s *RepositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// CreateStatus creates a new commit status.
func (s *RepositoryService) CreateStatus(ctx context.Context, repo, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// CreateDeployStatus creates a new deployment status.
func (s *RepositoryService) CreateDeployStatus(ctx context.Context, repo string, input *scm.DeployStatus) (*scm.DeployStatus, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// UpdateHook updates a repository webhook.
func (s *RepositoryService) UpdateHook(ctx context.Context, repo, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// DeleteHook deletes a repository webhook.
func (s *RepositoryService) DeleteHook(ctx context.Context, repo, id string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}
