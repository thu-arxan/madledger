package common

import (
	"madledger/common/util"

	"golang.org/x/crypto/ripemd160"
)

// NewContractAddress return the address of the contract
func NewContractAddress(channelID string, caller Address, code []byte) (newAddr Address) {
	// temp := make([]byte, 32+8)
	// copy(temp, caller[:])
	// PutUint64BE(temp[32:], uint64(sequence))
	temp := util.BytesCombine([]byte(channelID), caller.Bytes(), code)
	hasher := ripemd160.New()
	hasher.Write(temp) // does not error
	copy(newAddr[:], hasher.Sum(nil))
	return
}
