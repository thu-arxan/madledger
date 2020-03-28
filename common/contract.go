// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

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
