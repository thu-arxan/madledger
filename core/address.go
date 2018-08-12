package core

import (
	"fmt"
)

// TODO: This is copied from ethereum.
// However, how to use these is still a question.
// Also, there maybe many things is missed.
// Besides, how to support different kind of encrypted functions
// is still a problem.

const (
	// AddressLength is the expected length of the adddress
	AddressLength = 20
)

// Address represents the 20 byte address of an MadLedger account.
type Address [AddressLength]byte

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

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// MadLedger address or not.
func IsHexAddress(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}
