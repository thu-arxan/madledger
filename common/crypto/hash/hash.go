// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package hash

import (
	"crypto/sha256"
	"madledger/common/crypto/openssl/sm3"
)

// Hash is just wrapper now.
// TODO: Change it.
func Hash(data []byte) []byte {
	return SM3(data)
}

// SHA256 return the sha256 of data
func SHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// SM3 return the sm3 of data
func SM3(data []byte) []byte {
	return sm3.Sm3Sum(data)
}
