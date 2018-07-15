// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bitbucket implements a Bitbucket Cloud client.
package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// ClientID    string
// ClientSecret    string
// Endpoint    string

// Source    scm.TokenSource
// Client    *http.Client

// New returns a new Bitbucket API client.
func New(uri string, opt ...Option) (*scm.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	opts := new(Options)
	for _, o := range opt {
		o(opts)
	}
	refresher := &oauth2.Refresher{
		Endpoint:     tokenEndpoint,
		ClientID:     opts.clientID,
		ClientSecret: opts.clientSecret,
		Client:       opts.client,
	}
	client := &wrapper{new(scm.Client)}
	client.BaseURL = base
	// initialize services
	client.Driver = scm.DriverGithub
	client.Contents = &contentService{client}
	client.Git = &gitService{client}
	client.Issues = &issueService{client}
	client.Organizations = &organizationService{client}
	client.PullRequests = &pullService{&issueService{client}}
	client.Repositories = &repositoryService{client}
	client.Reviews = &reviewService{client}
	client.Tokens = &tokenService{refresher}
	client.Users = &userService{client}
	client.Webhooks = &webhookService{client}
	return client.Client, nil
}

// NewDefault returns a new Bitbucket API client using the
// default api.bitbucket.org address.
func NewDefault(opt ...Option) *scm.Client {
	client, _ := New("https://api.bitbucket.org", opt...)
	return client
}

// wraper wraps the Client to provide high level helper functions
// for making http requests and unmarshaling the response.
type wrapper struct {
	*scm.Client
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*scm.Response, error) {
	req := &scm.Request{
		Method: method,
		Path:   path,
	}
	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(in)
		req.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		req.Body = buf
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status == 401 {
		return res, scm.ErrNotAuthorized
	} else if res.Status > 300 {
		err := new(Error)
		json.NewDecoder(res.Body).Decode(err)
		return res, err
	}

	if out == nil {
		return res, nil
	}

	// if raw output is expected, copy to the provided
	// buffer and exit.
	if w, ok := out.(io.Writer); ok {
		io.Copy(w, res.Body)
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	return res, json.NewDecoder(res.Body).Decode(out)
}

// pagination represents Bitbucket pagination properties
// embedded in list responses.
type pagination struct {
	PageLen int    `json:"pagelen"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Next    string `json:"next"`
}

// Error represents a Bitbucket error.
type Error struct {
	Type string `json:"type"`
	Data struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *Error) Error() string {
	return e.Data.Message
}
