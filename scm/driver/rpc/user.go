// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(scm.User)
	res, err := s.client.do(ctx, "GET", "/user", nil, out)
	return out, res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	return "", nil, scm.ErrNotSupported
}
