// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(user)
	res, err := s.client.do(ctx, "GET", "api/v1/user", nil, out)
	return convertUser(out), res, err
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

//
// native data structures
//

type user struct {
	Admin       bool   `json:"admin"`
	Blocked     bool   `json:"blocked"`
	Created     int    `json:"created"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	UID         string `json:"uid"`
	Updated     int    `json:"updated"`
}

//
// native data structure conversion
//

func convertUser(src *user) *scm.User {
	return &scm.User{
		Login: src.Email,
		Email: src.Email,
		Name:  src.DisplayName,
		ID:    src.UID,
	}
}
