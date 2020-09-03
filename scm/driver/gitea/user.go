// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"github.com/jenkins-x/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out, err := s.client.GiteaClient.GetMyUserInfo()
	return convertGiteaUser(out), nil, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	out, err := s.client.GiteaClient.GetUserInfo(login)
	return convertGiteaUser(out), nil, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	user, res, err := s.Find(ctx)
	if user != nil {
		return user.Email, res, err
	}
	return "", res, err
}

func (s *userService) ListInvitations(context.Context) ([]*scm.Invitation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *userService) AcceptInvitation(context.Context, int64) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

//
// native data structures
//

type user struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
}

//
// native data structure conversion
//

func convertGiteaUsers(src []*gitea.User) []scm.User {
	answer := []scm.User{}
	for _, u := range src {
		user := convertGiteaUser(u)
		if user.Login != "" {
			answer = append(answer, *user)
		}
	}
	if len(answer) == 0 {
		return nil
	}
	return answer
}

func convertUsers(src []user) []scm.User {
	answer := []scm.User{}
	for _, u := range src {
		user := convertUser(&u)
		if user.Login != "" {
			answer = append(answer, *user)
		}
	}
	if len(answer) == 0 {
		return nil
	}
	return answer
}

func convertGiteaUser(src *gitea.User) *scm.User {
	if src == nil || src.UserName == "" {
		return nil
	}
	return &scm.User{
		ID:     int(src.ID),
		Login:  src.UserName,
		Name:   src.FullName,
		Email:  src.Email,
		Avatar: src.AvatarURL,
	}
}
func convertUser(src *user) *scm.User {
	return &scm.User{
		Login:  userLogin(src),
		Avatar: src.Avatar,
		Email:  src.Email,
		Name:   src.Fullname,
	}
}

func userLogin(src *user) string {
	if src.Username != "" {
		return src.Username
	}
	return src.Login
}
