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
	"encoding/binary"
	"fmt"
	"math/big"
)

var (
	// ZeroWord256 is zero of Word256
	ZeroWord256 = Word256{}
	// OneWord256 is one of Word256
	OneWord256 = LeftPadWord256([]byte{1})
	// BigWord256Length define the int which length is Word256
	BigWord256Length = big.NewInt(Word256Length)
	trimCutSet       = string([]byte{0})
)

// Word256Length define the length of Word256
const Word256Length = 32

// Word256 define the Word256
type Word256 [Word256Length]byte

func (w Word256) String() string {
	return string(w[:])
}

// Copy copy the Word256
func (w Word256) Copy() Word256 {
	return w
}

// Bytes return the bytes of Word256
func (w Word256) Bytes() []byte {
	return w[:]
}

// Prefix return the prefix of Word256
func (w Word256) Prefix(n int) []byte {
	return w[:n]
}

// Postfix return the postfix of Word256
func (w Word256) Postfix(n int) []byte {
	return w[32-n:]
}

// Word160 return a Word160 embedded a Word256 and padded on the left (as it is for account addresses in EVM)
func (w Word256) Word160() (w160 Word160) {
	copy(w160[:], w[Word256Word160Delta:])
	return
}

// IsZero return is a Word256 all zero
func (w Word256) IsZero() bool {
	accum := byte(0)
	for _, byt := range w {
		accum |= byt
	}
	return accum == 0
}

// Is64BitOverflow return if a word is int64 overflow
func (w Word256) Is64BitOverflow() bool {
	for i := 0; i < len(w)-8; i++ {
		if w[i] != 0 {
			return true
		}
	}
	return false
}

// LeftPadWord256 copy bz to the left of word
func LeftPadWord256(bz []byte) (word Word256) {
	copy(word[32-len(bz):], bz)
	return
}

// Uint64FromWord256 convert word to uint64
func Uint64FromWord256(word Word256) uint64 {
	return GetUint64BE(word.Postfix(8))
}

// Int64FromWord256 onvert word to int64
func Int64FromWord256(word Word256) int64 {
	return GetInt64BE(word.Postfix(8))
}

// GetInt64BE convert bytes to int64
func GetInt64BE(src []byte) int64 {
	return int64(binary.BigEndian.Uint64(src))
}

// GetUint64BE convert bytes to uint64
func GetUint64BE(src []byte) uint64 {
	return binary.BigEndian.Uint64(src)
}

// Int64ToWord256 convert int64 to word256
func Int64ToWord256(i int64) (word Word256) {
	PutInt64BE(word[24:], i)
	return
}

// Uint64ToWord256 convert uint64 to word256
func Uint64ToWord256(i uint64) (word Word256) {
	PutUint64BE(word[24:], i)
	return
}

// BytesToWord256 convert bytes to Word256
func BytesToWord256(data []byte) (Word256, error) {
	var word Word256
	if len(data) != Word256Length {
		return ZeroWord256, fmt.Errorf("The length of data is %d other than %d", len(data), Word256Length)
	}
	copy(word[:], data)
	return word, nil
}
