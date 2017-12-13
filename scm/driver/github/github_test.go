// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"testing"

	"github.com/drone/go-scm/scm"
)

func testRate(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Rate.Limit, 60; got != want {
			t.Errorf("Want X-RateLimit-Limit %d, got %d", want, got)
		}
		if got, want := res.Rate.Remaining, 59; got != want {
			t.Errorf("Want X-RateLimit-Remaining %d, got %d", want, got)
		}
		if got, want := res.Rate.Reset, int64(1512076018); got != want {
			t.Errorf("Want X-RateLimit-Reset %d, got %d", want, got)
		}
	}
}

func testPage(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Page.Next, 2; got != want {
			t.Errorf("Want next page %d, got %d", want, got)
		}
		if got, want := res.Page.Prev, 1; got != want {
			t.Errorf("Want prev page %d, got %d", want, got)
		}
		if got, want := res.Page.First, 1; got != want {
			t.Errorf("Want first page %d, got %d", want, got)
		}
		if got, want := res.Page.Last, 5; got != want {
			t.Errorf("Want last page %d, got %d", want, got)
		}
	}
}

func testRequest(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.ID, "DD0E:6011:12F21A8:1926790:5A2064E2"; got != want {
			t.Errorf("Want X-GitHub-Request-Id %q, got %q", want, got)
		}
	}
}
