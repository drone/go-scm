// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"github.com/jenkins-x/go-scm/scm"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) IsMember(ctx context.Context, org string, user string) (bool, *scm.Response, error) {
	isMember, err := s.client.GiteaClient.CheckOrgMembership(org, user)
	return isMember, nil, err
}

func (s *organizationService) IsAdmin(ctx context.Context, org string, user string) (bool, *scm.Response, error) {
	members, res, err := s.ListOrgMembers(ctx, org, scm.ListOptions{})
	if err != nil {
		return false, res, err
	}
	for _, m := range members {
		if m.Login == user && m.IsAdmin {
			return true, res, nil
		}
	}
	return false, res, nil
}

func (s *organizationService) ListTeams(ctx context.Context, org string, ops scm.ListOptions) ([]*scm.Team, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *organizationService) ListTeamMembers(ctx context.Context, id int, role string, ops scm.ListOptions) ([]*scm.TeamMember, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *organizationService) ListOrgMembers(ctx context.Context, org string, ops scm.ListOptions) ([]*scm.TeamMember, *scm.Response, error) {
	out, err := s.client.GiteaClient.ListOrgMembership(org, gitea.ListOrgMembershipOption{})
	return convertMemberList(out), nil, err
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	out, err := s.client.GiteaClient.GetOrg(name)
	return convertOrg(out), nil, err
}

func (s *organizationService) List(ctx context.Context, _ scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	out, err := s.client.GiteaClient.ListMyOrgs(gitea.ListOrgsOptions{})
	return convertOrgList(out), nil, err
}

//
// native data structure conversion
//

func convertOrgList(from []*gitea.Organization) []*scm.Organization {
	to := []*scm.Organization{}
	for _, v := range from {
		to = append(to, convertOrg(v))
	}
	return to
}

func convertOrg(from *gitea.Organization) *scm.Organization {
	if from == nil || from.UserName == "" {
		return nil
	}
	return &scm.Organization{
		Name:   from.UserName,
		Avatar: from.AvatarURL,
	}
}

func convertMemberList(from []*gitea.User) []*scm.TeamMember {
	var to []*scm.TeamMember
	for _, v := range from {
		to = append(to, convertMember(v))
	}
	return to
}

func convertMember(from *gitea.User) *scm.TeamMember {
	return &scm.TeamMember{
		Login:   from.UserName,
		IsAdmin: from.IsAdmin,
	}
}
