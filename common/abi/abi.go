package abi

import (
	eabi "evm/abi"
	"fmt"
	"madledger/common/util"
)

// This file is a wrapper of madledger/common/abi

func init() {
	eabi.SetAddressParser(20, func(bs []byte) string {
		return "0x" + fmt.Sprintf("%x", bs)
	}, func(addr string) ([]byte, error) {
		return util.HexToBytes(addr)
	})
}

// Pack provide a easy way to pack
func Pack(abiFile, funcName string, inputs ...string) ([]byte, error) {
	return eabi.Pack(abiFile, funcName, inputs...)
}

// Unpack provide a easy way to unpack
func Unpack(abiFile, funcName string, data []byte) (values []string, err error) {
	return eabi.Unpack(abiFile, funcName, data)
}
