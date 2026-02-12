// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("2.0/workspaces/%s", name)
	out := new(organization)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertOrganization(out), res, err
}

func (s *organizationService) FindMembership(ctx context.Context, name, username string) (*scm.Membership, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *organizationService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	// Use /2.0/user/workspaces endpoint (replaces deprecated /2.0/workspaces)
	path := fmt.Sprintf("2.0/user/workspaces?%s", encodeListOptions(opts))
	out := new(workspaceAccessList)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	if err != nil {
		return nil, res, err
	}
	copyPagination(out.pagination, res)
	return convertWorkspaceAccessList(out), res, err
}

func convertWorkspaceAccessList(from *workspaceAccessList) []*scm.Organization {
	to := []*scm.Organization{}
	for _, v := range from.Values {
		if v.Workspace != nil {
			to = append(to, convertWorkspace(v.Workspace))
		}
	}
	return to
}

func convertWorkspace(from *workspace) *scm.Organization {
	avatar := ""
	if from.Links.Avatar.Href != "" {
		avatar = from.Links.Avatar.Href
	} else {
		avatar = fmt.Sprintf("https://bitbucket.org/account/%s/avatar/32/", from.Slug)
	}
	return &scm.Organization{
		Name:   from.Slug,
		Avatar: avatar,
	}
}

type organization struct {
	Login string `json:"slug"`
}

func convertOrganization(from *organization) *scm.Organization {
	return &scm.Organization{
		Name:   from.Login,
		Avatar: fmt.Sprintf("https://bitbucket.org/account/%s/avatar/32/", from.Login),
	}
}
