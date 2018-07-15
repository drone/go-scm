// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"context"
	"time"
)

type (
	// Token represents the credentials used to authorize
	// the requests to access protected resources.
	Token struct {
		Token   string
		Refresh string
		Expires time.Time
	}

	// TokenSource returns a token.
	TokenSource interface {
		Token(context.Context) (*Token, error)
	}

	// TokenService provides token refresh capabilities.
	TokenService interface {
		// Refresh refreshes the token if necessary. The
		// function returns a boolean value indicating
		// whether or not the token is refreshed.
		Refresh(context.Context, *Token) (bool, error)
	}

	// TokenKey is the key to use with the context.WithValue
	// function to associate an Token value with a context.
	TokenKey struct{}
)

// WithToken returns a copy of parent in which the token value is set
func WithToken(parent context.Context, token *Token) context.Context {
	return context.WithValue(parent, TokenKey{}, token)
}
