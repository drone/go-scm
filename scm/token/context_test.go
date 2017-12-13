// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := context.Background()

	_, ok := FromContext(ctx)
	if ok {
		t.Errorf("Expect FromContext returns a nil token")
	}

	want := new(Token)
	ctx = WithContext(ctx, want)

	got, _ := FromContext(ctx)
	if want != got {
		t.Errorf("Expect FromContext returns the Token")
	}
}
