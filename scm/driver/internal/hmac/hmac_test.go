// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hmac

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
			msg: "hello world",
			key: "",
			sig: "512d57985a1de49af20cfa7b785c07a701a0b7b5f2ff06726b3ac077066c1992",
			res: false,
		},
		{
			msg: "hola mundo",
			key: "",
			sig: "",
			res: false,
		},
		{
			msg: "bonjour monde",
			key: "root",
			sig: "512d57985a1de49af20cfa7b785c07a701a0b7b5f2ff06726b3ac077066c1992",
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
