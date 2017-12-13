// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// expiryDelta determines how earlier a token should be considered
// expired than its actual expiration time. It is used to avoid late
// expirations due to client-server time mismatches.
const expiryDelta = 30 * time.Second

// Refresher is an http.RoundTripper that refreshes oauth
// tokens, wrapping a base RoundTripper and refreshing the
// token if expired.
type Refresher struct {
	ClientID     string
	ClientSecret string
	Endpoint     string

	Base http.RoundTripper
}

// RoundTrip authorizes and authenticates the request with
// an access token from the request context.
func (t *Refresher) RoundTrip(r *http.Request) (*http.Response, error) {
	token, ok := FromContext(r.Context())
	if ok && expired(token) {
		if err := t.refresh(token); err != nil {
			return nil, err
		}
	}
	return t.Base.RoundTrip(r)
}

// refresh refreshes the expired token.
func (t *Refresher) refresh(token *Token) error {
	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", token.Refresh)

	reader := strings.NewReader(
		values.Encode(),
	)
	req, err := http.NewRequest("POST", t.Endpoint, reader)
	if err != nil {
		return err
	}
	req.SetBasicAuth(t.ClientID, t.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := t.Base.RoundTrip(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		out := new(tokenError)
		err = json.NewDecoder(res.Body).Decode(res)
		if err != nil {
			return err
		}
		return out
	}

	out := new(tokenGrant)
	err = json.NewDecoder(res.Body).Decode(out)
	if err != nil {
		return err
	}

	token.Access = out.Access
	token.Refresh = out.Refresh
	token.Expires = time.Now().Add(
		time.Duration(out.Expires) * time.Second,
	)
	return nil
}

// expired reports whether the token is expired.
func expired(token *Token) bool {
	if token.Expires.IsZero() && len(token.Access) != 0 {
		return false
	}
	return token.Expires.Add(-expiryDelta).
		Before(time.Now())
}
