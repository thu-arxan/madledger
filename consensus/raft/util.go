package raft

import "madledger/common/crypto/hash"

// Hash is a wrapper of sha256
func Hash(data []byte) []byte {
	return hash.SHA256(data)
}
