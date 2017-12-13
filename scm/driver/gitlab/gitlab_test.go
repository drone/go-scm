// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"testing"

	"github.com/drone/go-scm/scm"
)

func testRate(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Rate.Limit, 600; got != want {
			t.Errorf("Want RateLimit-Limit %d, got %d", want, got)
		}
		if got, want := res.Rate.Remaining, 599; got != want {
			t.Errorf("Want RateLimit-Remaining %d, got %d", want, got)
		}
		if got, want := res.Rate.Reset, int64(1512454441); got != want {
			t.Errorf("Want RateLimit-Reset %d, got %d", want, got)
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
		if got, want := res.ID, "0d511a76-2ade-4c34-af0d-d17e84adb255"; got != want {
			t.Errorf("Want X-Request-Id: %q, got %q", want, got)
		}
	}
}
