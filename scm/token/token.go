// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token provided oauth token authorization.
package token

import "time"

// Token represents an authentication token.
type Token struct {
	Access  string
	Refresh string
	Expires time.Time
}

type tokenGrant struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh"`
	Expires int64  `json:"expires_in"`
}

type tokenError struct {
	Code    string `json:"error"`
	Message string `json:"error_description"`
}

func (t *tokenError) Error() string {
	return t.Message
}
