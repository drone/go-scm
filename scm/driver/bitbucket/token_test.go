// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"testing"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/h2non/gock"
)

func TestTokenRefresh(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.org").
		Post("/site/oauth2/access_token").
		MatchHeader("Authorization", "Basic NTU5OTE4YTgwODowMmJiYTUwMTJm").
		Reply(200).
		BodyString(`
		{
			"access_token": "9698fa6a8113b3",
			"expires_in": 7200,
			"refresh_token": "3a2bfce4cb9b0f",
			"token_type": "bearer"
		}
	`)

	token := &scm.Token{
		Token:   "ae215a0a8223a9",
		Refresh: "3a2bfce4cb9b0f",
		Expires: time.Now().Add(-time.Second),
	}

	client, _ := New("https://api.bitbucket.org",
		WithClientID("559918a808"),
		WithClientSecret("02bba5012f"),
	)
	refreshed, err := client.Tokens.Refresh(context.Background(), token)
	if err != nil {
		t.Error(err)
	}

	if !refreshed {
		t.Errorf("Expected token refresh")
	}
	if got, want := token.Token, "9698fa6a8113b3"; got != want {
		t.Errorf("Expected refresh token %s, got %s", want, got)
	}
}

func TestTokenRefresh_NoRefresh(t *testing.T) {
	token := &scm.Token{
		Token:   "ae215a0a8223a9",
		Refresh: "3a2bfce4cb9b0f",
		Expires: time.Now().Add(time.Hour),
	}

	client, _ := New("https://api.bitbucket.org",
		WithClientID("559918a808"),
		WithClientSecret("02bba5012f"),
	)
	refreshed, err := client.Tokens.Refresh(context.Background(), token)
	if err != nil {
		t.Error(err)
	}

	if refreshed {
		t.Errorf("Expected token not refreshed")
	}
}
