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
)

// Hash return the hash of date
func Hash(data []byte) []byte {
	return hashWithSHA256(data)
}

// hashWithSHA256 is the implementation of SHA256
func hashWithSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
