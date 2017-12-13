// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import "net/http"

// Transport is an http.RoundTripper that makes oauth
// requests, wrapping a base RoundTripper and adding an
// Authorization header with a token from the request
// Context.
type Transport struct {
	Base http.RoundTripper

	// SetToken defines an optional func to write the token
	// to the http.Request.
	SetToken func(*http.Request)
}

// RoundTrip authorizes and authenticates the request with
// an access token from the request context.
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Do not overwrite the authorization header if exists.
	if r.Header.Get("Authorization") != "" {
		return t.base().RoundTrip(r)
	}
	r2 := cloneRequest(r)
	if t.SetToken != nil {
		t.SetToken(r2)
	} else {
		setToken(r2)
	}
	return t.base().RoundTrip(r2)
}

// base returns the base transport. If no base transport
// is configured, the default transport is returned.
func (t *Transport) base() http.RoundTripper {
	switch {
	case t.Base == nil:
		return http.DefaultTransport
	default:
		return t.Base
	}
}

// setToken sets the Authorization header for the current
// request using the token from the context.
func setToken(r *http.Request) {
	token, ok := FromContext(r.Context())
	if ok {
		r.Header.Set("Authorization", "Bearer "+token.Access)
	}
}

// cloneRequest returns a clone of the provided
// http.Request. The clone is a shallow copy of the struct
// and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
