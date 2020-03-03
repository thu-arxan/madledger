package crypto

import "madledger/common/crypto/hash"

// Hash return the hash of date
func Hash(data []byte) []byte {
	return hash.Hash(data)
}
