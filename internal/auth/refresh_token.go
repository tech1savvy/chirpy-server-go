package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// MakeRefreshToken makes a random 256 bit token
// encoded in hex
func MakeRefreshToken() string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes)
}
