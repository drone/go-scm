// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import "context"

// The key type is unexported to prevent collisions
type key int

// tokenKey is the context key for the request token.
const tokenKey key = iota

// WithContext returns a copy of parent in which the token value is set
func WithContext(parent context.Context, token *Token) context.Context {
	return context.WithValue(parent, tokenKey, token)
}

// FromContext returns the value of the token key on the ctx
func FromContext(ctx context.Context) (*Token, bool) {
	token, ok := ctx.Value(tokenKey).(*Token)
	return token, ok
}
