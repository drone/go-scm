// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import "net/http"

// Options provides Bitbucket client options.
type Options struct {
	clientID     string
	clientSecret string
	client       *http.Client
}

// Option provides a Bitbucket client option.
type Option func(*Options)

// WithClient returns an option to set the http.Client
// used to refresh authorization tokens.
func WithClient(client *http.Client) Option {
	return func(opts *Options) {
		opts.client = client
	}
}

// WithClientID returns an option to set the Bitbucket
// oauth2 client identifier.
func WithClientID(clientID string) Option {
	return func(opts *Options) {
		opts.clientID = clientID
	}
}

// WithClientSecret returns an option to set the
// Bitbucket oauth2 client secret.
func WithClientSecret(clientSecret string) Option {
	return func(opts *Options) {
		opts.clientSecret = clientSecret
	}
}
