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

const (
	// Word160Length define the length of Word160
	Word160Length = 20
	// Word256Word160Delta define the delta between the Word160 and Word256
	Word256Word160Delta = 12
)

// ZeroWord160 is the zero of Word160
var ZeroWord160 = Word160{}

// Word160 is the bytes which length is 20
type Word160 [Word160Length]byte

// Word256 return the Word256 of a Word160
func (w Word160) Word256() (word256 Word256) {
	copy(word256[Word256Word160Delta:], w[:])
	return
}

// Bytes return the bytes of Word160
func (w Word160) Bytes() []byte {
	return w[:]
}
