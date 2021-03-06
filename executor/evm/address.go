// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package evm

import (
	"madledger/common"

	"github.com/thu-arxan/evm/util"
)

// Address is the address
type Address [20]byte

// Bytes is the implementation of interface
func (a *Address) Bytes() []byte {
	return a[:]
}

// NewAddressFromCommon ...
func NewAddressFromCommon(addr common.Address) *Address {
	return BytesToAddress(addr.Bytes())
}

// BytesToAddress convert bytes to address
func BytesToAddress(bytes []byte) *Address {
	var a Address
	copy(a[:], util.FixBytesLength(bytes, 20))
	return &a
}

// HexToAddress convert hex string to address, string may begin with 0x, 0X or nothing
func HexToAddress(hex string) *Address {
	var a Address
	if bytes, err := util.HexToBytes(hex); err == nil {
		copy(a[:], util.FixBytesLength(bytes, 20))
	}
	return &a
}

// ZeroAddress return zero address
func ZeroAddress() *Address {
	var a Address
	return &a
}
