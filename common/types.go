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
	"fmt"
	"madledger/common/util"
	"math/big"

	"github.com/tmthrgd/go-hex"
)

// This is copied from ethereum.
// However, there maybe many things is missed like
// how to support different kind of encrypted functions.

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the adddress
	AddressLength = 20
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// ZeroHash is the zero of Hash
var ZeroHash = Hash{}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Bytes return the bytes of hash
func (h Hash) Bytes() []byte {
	return h[:]
}

// Word256 return the Word256 format of hash
func (h Hash) Word256() Word256 {
	return Word256(h)
}

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Address represents the 20 byte address of an MadLedger account.
type Address [AddressLength]byte

// ZeroAddress is the zero of Address
var ZeroAddress = Address{}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address {
	return BytesToAddress(FromHex(s))
}

// AddressFromHexString convert hex string to address
// Note: The str can not contain the prefix 0x
func AddressFromHexString(str string) (Address, error) {
	bs, err := hex.DecodeString(str)
	if err != nil {
		return ZeroAddress, err
	}
	return AddressFromBytes(bs)
}

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// MadLedger address or not.
func IsHexAddress(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// Word256 return the Word256 format of address
func (a Address) Word256() Word256 {
	return Word160(a).Word256()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// Bytes return the bytes of address
func (a Address) Bytes() []byte {
	return a[:]
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// String return the upper 0x
func (a Address) String() string {
	// return hex.EncodeUpperToString(a[:])
	return "0x" + util.Hex(a[:])
}

// AddressFromWord256 decode address from Word256
func AddressFromWord256(addr Word256) Address {
	return Address(addr.Word160())
}

// AddressFromBytes returns an address consisting of the first 20 bytes of bs, return an error if the bs does not have length exactly 20
// but will still return either: the bytes in bs padded on the right or the first 20 bytes of bs truncated in any case.
func AddressFromBytes(bs []byte) (address Address, err error) {
	if len(bs) != Word160Length {
		err = fmt.Errorf("slice passed as address '%X' has %d bytes but should have %d bytes",
			bs, len(bs), Word160Length)
		// It is caller's responsibility to check for errors. If they ignore the error we'll assume they want the
		// best-effort mapping of the bytes passed to an address so we don't return here
	}
	copy(address[:], bs)
	return
}

// AddressFromChannelID return a hex encoded version of channelid
// encoded bytes would be truncated if encoded len is bigger than AddressLength
func AddressFromChannelID(channelID string) Address {
	var address Address
	address.SetBytes(util.Hex(channelID[:]))
	return address
}
