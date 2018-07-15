// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// bitbucket cloud access_token endpoint.
const tokenEndpoint = "https://bitbucket.org/site/oauth2/access_token"

type tokenService struct {
	refresher *oauth2.Refresher
}

func (t *tokenService) Refresh(ctx context.Context, token *scm.Token) (bool, error) {
	if oauth2.Expired(token) == false {
		return false, nil
	}
	t1 := token.Token
	err := t.refresher.Refresh(token)
	t2 := token.Token
	return t1 != t2, err
}
