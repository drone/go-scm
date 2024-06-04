// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"gopkg.in/h2non/gock.v1"
)

func TestTransport(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/user").
		MatchHeader("Authorization", "Bearer mF_9.B5f-4.1JqM").
		Reply(200)

	client := &http.Client{
		Transport: &Transport{
			Source: StaticTokenSource(
				&scm.Token{
					Token: "mF_9.B5f-4.1JqM",
				},
			),
		},
	}

	res, err := client.Get("https://api.github.com/user")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
}

func TestTransport_CustomScheme(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gogs.io").
		Get("/api/v1/user").
		MatchHeader("Authorization", "token mF_9.B5f-4.1JqM").
		Reply(200)

	client := &http.Client{
		Transport: &Transport{
			Scheme: "token",
			Source: StaticTokenSource(
				&scm.Token{
					Token: "mF_9.B5f-4.1JqM",
				},
			),
		},
	}

	res, err := client.Get("https://try.gogs.io/api/v1/user")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
}

func TestTransport_NoToken(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/user").
		Reply(200)

	client := &http.Client{
		Transport: &Transport{
			Source: ContextTokenSource(),
		},
	}

	res, err := client.Get("https://api.github.com/user")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
}

func TestTransport_TokenError(t *testing.T) {
	want := errors.New("Cannot retrieve token")
	client := &http.Client{
		Transport: &Transport{
			Source: mockErrorSource{want},
		},
	}

	resp, err := client.Get("https://api.github.com/user")
	if err == nil {
		defer resp.Body.Close()
		t.Errorf("Expect token source error, got nil")
	}
}

type mockErrorSource struct {
	err error
}

func (s mockErrorSource) Token(_ context.Context) (*scm.Token, error) {
	return nil, s.err
}
