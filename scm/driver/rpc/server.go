// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/drone/go-scm/scm"
)

// Handler exposes the underlying scm.Client as a handler.
func Handler(client *scm.Client) http.Handler {
	return &handler{client: client}
}

type handler struct {
	client *scm.Client
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var out interface{}
	var res *scm.Response
	var err error

	switch r.URL.Path {
	case "/user":
		out, res, err = h.client.Users.Find(r.Context())
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("X-RateLimit-Limit", fmt.Sprint(res.Rate.Limit))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprint(res.Rate.Remaining))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprint(res.Rate.Reset))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Error{Message: err.Error()})
	} else {
		json.NewEncoder(w).Encode(out)
	}
}
