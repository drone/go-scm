// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha1

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

// Validate checks the hmac signature of the message.
func Validate(message, key, signature []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	sum := mac.Sum(nil)
	return hmac.Equal(signature, sum)
}

// ValidateEncoded checks the hmac signature of the mssasge
// using a hex encoded signature.
func ValidateEncoded(message, key []byte, signature string) bool {
	decoded, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return Validate(message, key, decoded)
}
