// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	// connectionData works on both Azure DevOps Services and Azure DevOps Server,
	// unlike the Services-only vssps profile API.
	endpoint := fmt.Sprintf("%s/_apis/connectionData", s.client.owner)
	out := new(connectionData)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	if err != nil {
		return nil, res, err
	}
	return convertConnectionData(out), res, nil
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	return "", nil, scm.ErrNotSupported
}

func (s *userService) ListEmail(context.Context, scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type connectionData struct {
	AuthenticatedUser struct {
		ID                  string `json:"id"`
		ProviderDisplayName string `json:"providerDisplayName"`
		Properties          struct {
			Account struct {
				Value string `json:"$value"`
			} `json:"Account"`
		} `json:"properties"`
	} `json:"authenticatedUser"`
}

func convertConnectionData(from *connectionData) *scm.User {
	return &scm.User{
		ID:    from.AuthenticatedUser.ID,
		Name:  from.AuthenticatedUser.ProviderDisplayName,
		Login: from.AuthenticatedUser.Properties.Account.Value,
	}
}
