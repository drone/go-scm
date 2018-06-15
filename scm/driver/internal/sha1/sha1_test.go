// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha1

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		msg string
		key string
		sig string
		res bool
	}{
		{
			msg: `{"yo":true}`,
			key: "0123456789abcdef",
			sig: "126f2c800419c60137ce748d7672e77b65cf16d6",
			res: true,
		},
	}

	for _, test := range tests {
		res := ValidateEncoded(
			[]byte(test.msg),
			[]byte(test.key),
			test.sig,
		)
		if res != test.res {
			t.Errorf("Want valid %v for message %q",
				test.res, test.msg)
		}
	}
}
