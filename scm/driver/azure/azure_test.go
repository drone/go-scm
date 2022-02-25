// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"testing"

	"github.com/drone/go-scm/scm"
)

var mockHeaders = map[string]string{
	"X-Request-Id": "7cb049dbf1fafae67b4f0aa81ca7e870",
}

var mockPageHeaders = map[string]string{
	"total_page": `3`,
}

func TestClient(t *testing.T) {
	client, err := New("https://gitee.com/api/v5")
	if err != nil {
		t.Error(err)
	}
	if got, want := client.BaseURL.String(), "https://gitee.com/api/v5/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Base(t *testing.T) {
	client, err := New("https://gitee.com/api/v5")
	if err != nil {
		t.Error(err)
	}
	got, want := client.BaseURL.String(), "https://gitee.com/api/v5/"
	if got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Default(t *testing.T) {
	client := NewDefault()
	if got, want := client.BaseURL.String(), "https://gitee.com/api/v5/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Error(t *testing.T) {
	_, err := New("http://a b.com/")
	if err == nil {
		t.Errorf("Expect error when invalid URL")
	}
}

func testPage(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Page.Prev, 1; got != want {
			t.Errorf("Want prev page %d, got %d", want, got)
		}
		if got, want := res.Page.Next, 2; got != want {
			t.Errorf("Want next page %d, got %d", want, got)
		}
		if got, want := res.Page.Last, 3; got != want {
			t.Errorf("Want last page %d, got %d", want, got)
		}
		if got, want := res.Page.First, 1; got != want {
			t.Errorf("Want first page %d, got %d", want, got)
		}
	}
}

func testRequest(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.ID, "7cb049dbf1fafae67b4f0aa81ca7e870"; got != want {
			t.Errorf("Want X-Request-Id %q, got %q", want, got)
		}
	}
}
