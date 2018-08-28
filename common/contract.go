package common

import (
	"golang.org/x/crypto/ripemd160"
)

// NewContractAddress return the address of the contract
func NewContractAddress(caller Address, sequence uint64) (newAddr Address) {
	temp := make([]byte, 32+8)
	copy(temp, caller[:])
	PutUint64BE(temp[32:], uint64(sequence))
	hasher := ripemd160.New()
	hasher.Write(temp) // does not error
	copy(newAddr[:], hasher.Sum(nil))
	return
}
