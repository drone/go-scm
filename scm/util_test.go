// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import "testing"

func TestSplit(t *testing.T) {
	tests := []struct {
		value, owner, name string
	}{
		{"octocat/hello-world", "octocat", "hello-world"},
		{"octocat/hello/world", "octocat", "hello/world"},
		{"hello-world", "", "hello-world"},
		{value: ""}, // empty value returns nothing
	}
	for _, test := range tests {
		owner, name := Split(test.value)
		if got, want := owner, test.owner; got != want {
			t.Errorf("Got repository owner %s, want %s", got, want)
		}
		if got, want := name, test.name; got != want {
			t.Errorf("Got repository name %s, want %s", got, want)
		}
	}
}

func TestJoin(t *testing.T) {
	got, want := Join("octocat", "hello-world"), "octocat/hello-world"
	if got != want {
		t.Errorf("Got repository name %s, want %s", got, want)
	}
}
