package util

import "crypto/sha256"

// Hash return the hash of date
func Hash(data []byte) []byte {
	return hashWithSHA256(data)
}

// hashWithSHA256 is the implementation of SHA256
func hashWithSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
