// Package sslrandom includes functions that make the use of randomly generated contents easier for me
package sslrandom

// <---------------------------------------------------------------------------------------------------->

import (
	"crypto/rand"
	"encoding/base64"
)

// <---------------------------------------------------------------------------------------------------->

// String generates a random string of Base64 characters
func String(length int) (string, error) {
	// create a byte slice
	bytes := make([]byte, length)

	// populate it with random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return base64.StdEncoding.EncodeToString(bytes), err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}
