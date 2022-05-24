// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"errors"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(user)
	res, err := s.client.do(ctx, "GET", "2.0/user", nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := fmt.Sprintf("2.0/users/%s", login)
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	out := new(emails)
	res, err := s.client.do(ctx, "GET", "2.0/user/emails", nil, out)
	if err != nil {
		return "", res, err
	}

	for _, emailObj := range out.Values {
		if emailObj.IsPrimary {
			return emailObj.Email, res, nil
		}
	}
	return "", res, errors.New("no primary email found")
}

func (s *userService) ListEmail(context.Context, scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type user struct {
	// The `username` field is no longer available after 29 April 2019 in
	// accordance with GDPR regulations. See:
	// https://developer.atlassian.com/cloud/bitbucket/bitbucket-api-changes-gdpr/
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AccountID   string `json:"account_id"`
	Nickname    string `json:"nickname"`
	Links       struct {
		Self   link `json:"self"`
		HTML   link `json:"html"`
		Avatar link `json:"avatar"`
	} `json:"links"`
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

type emails struct {
	Values []struct {
		IsPrimary   bool   `json:"is_primary"`
		IsConfirmed bool   `json:"is_confirmed"`
		Email       string `json:"email"`
	} `json:"values"`
}

func convertUser(from *user) *scm.User {
	return &scm.User{
		Avatar: fmt.Sprintf("https://bitbucket.org/account/%s/avatar/32/", from.Username),
		Login:  from.Username,
		Name:   from.DisplayName,
	}
}
