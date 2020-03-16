// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package abi

import (
	"fmt"
	"madledger/common/util"

	eabi "github.com/thu-arxan/evm/abi"
)

// This file is a wrapper of github.com/thu-arxan/evm/abi
// Also, it set the address parser of abi

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
