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
	path := fmt.Sprintf("2.0/user/workspaces?%s", encodeListRoleOptions(opts))
	out := new(workspaceAccessList)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertWorkspaceAccessList(out), res, err
}

func convertWorkspaceAccessList(from *workspaceAccessList) []*scm.Organization {
	to := []*scm.Organization{}
	for _, values := range from.Values {
		if values.Workspace != nil {
			to = append(to, convertWorkspace(values.Workspace))
		}
	}
	return to
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

func convertWorkspace(workspace *workspace) *scm.Organization {
	avatar := ""
	if workspace.Links.Avatar.Href != "" {
		avatar = workspace.Links.Avatar.Href
	} else {
		avatar = fmt.Sprintf("https://bitbucket.org/account/%s/avatar/32", workspace.Slug)
	}
	return &scm.Organization{
		Name:   workspace.Slug,
		Avatar: avatar,
	}
}
