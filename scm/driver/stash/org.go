// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"

	"github.com/jenkins-x/go-scm/scm"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) ListTeams(ctx context.Context, org string, ops scm.ListOptions) ([]*scm.Team, *scm.Response, error) {
	panic("implement me")
}

func (s *organizationService) ListTeamMembers(ctx context.Context, id int, role string, ops scm.ListOptions) ([]*scm.TeamMember, *scm.Response, error) {
	panic("implement me")
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *organizationService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}
