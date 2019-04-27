package common

import (
	"encoding/binary"
	"math/big"
)

// Common big integers often used
var (
	Big1   = big.NewInt(1)
	Big2   = big.NewInt(2)
	Big3   = big.NewInt(3)
	Big0   = big.NewInt(0)
	Big32  = big.NewInt(32)
	Big256 = big.NewInt(256)
	Big257 = big.NewInt(257)
)

var big1 = big.NewInt(1)
var tt256 = new(big.Int).Lsh(big1, 256)
var tt256m1 = new(big.Int).Sub(new(big.Int).Lsh(big1, 256), big1)
var tt255 = new(big.Int).Lsh(big1, 255)

// U256 converts a possibly negative big int x into a positive big int encoding a twos complement representation of x
// truncated to 32 bytes
func U256(x *big.Int) *big.Int {
	// Note that the And operation induces big.Int to hold a positive representation of a negative number
	return new(big.Int).And(x, tt256m1)
}

// S256 interprets a positive big.Int as a 256-bit two's complement signed integer
func S256(x *big.Int) *big.Int {
	// Sign bit not set, value is its positive self
	if x.Cmp(tt255) < 0 {
		return x
	}
	// negative value is represented
	return new(big.Int).Sub(x, tt256)
}

// PutUint64BE set an uint64
func PutUint64BE(dest []byte, i uint64) {
	binary.BigEndian.PutUint64(dest, i)
}

// PutInt64BE set an int64
func PutInt64BE(dest []byte, i int64) {
	binary.BigEndian.PutUint64(dest, uint64(i))
}

// SignExtend treats the positive big int x as if it contains an embedded a back + 1 byte signed integer in its least significant
// bits and extends that sign
func SignExtend(back uint64, x *big.Int) *big.Int {
	// we assume x contains a signed integer of back + 1 bytes width
	// most significant bit of the back'th byte,
	signBit := back*8 + 7
	// single bit set at sign bit position
	mask := new(big.Int).Lsh(big1, uint(signBit))
	// all bits below sign bit set to 1 all above (including sign bit) set to 0
	mask.Sub(mask, big1)
	if x.Bit(int(signBit)) == 1 {
		// Number represented is negative - set all bits above sign bit (including sign bit)
		return x.Or(x, mask.Not(mask))
	}
	// Number represented is positive - clear all bits above sign bit (including sign bit)
	return x.And(x, mask)
}
