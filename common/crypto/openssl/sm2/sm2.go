// +build cgo

package sm2

/*
#include "util.h"
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"math/big"
	"unsafe"

	"github.com/tjfoc/gmsm/sm2"
)

// Sign sign with uid, return bytes sign
func Sign(priv *PrivateKey, msg, uid []byte) ([]byte, error) {
	var siglen C.size_t
	sig := make([]byte, 128)

	cmsg := C.CString(string(msg))
	cuid := C.CString(string(uid))

	defer func() {
		C.free(unsafe.Pointer(cmsg))
		C.free(unsafe.Pointer(cuid))
	}()

	if C.Sm2Sign(((*C.char)(unsafe.Pointer(&priv.Key[0]))), C.int(len(priv.Key)), ((*C.char)(unsafe.Pointer(cmsg))), C.size_t(len(msg)), ((*C.char)(unsafe.Pointer(cuid))), C.size_t(len(uid)), ((*C.uchar)(unsafe.Pointer(&sig[0]))), &siglen) != 1 {
		return nil, errors.New("Sign Failed: " + OpenError())
	}

	return sig[:siglen], nil
}

// Verify verify sm2 sig data with uid
func Verify(pub *PublicKey, msg, uid []byte, sig []byte) bool {
	cmsg := C.CString(string(msg))
	cuid := C.CString(string(uid))
	csig := C.CString(string(sig))

	defer func() {
		C.free(unsafe.Pointer(cmsg))
		C.free(unsafe.Pointer(cuid))
		C.free(unsafe.Pointer(csig))
	}()

	if C.Sm2Verify(((*C.char)(unsafe.Pointer(&pub.Key[0]))), C.int(len(pub.Key)), ((*C.char)(unsafe.Pointer(cmsg))), C.size_t(len(msg)), ((*C.char)(unsafe.Pointer(cuid))), C.size_t(len(uid)), ((*C.uchar)(unsafe.Pointer(csig))), C.size_t(len(sig))) != 1 {
		return false
	}

	return true
}

// ParseSm2PublicKey parse the public key from pkcs8 encoded bytes
func ParseSm2PublicKey(der []byte) (*PublicKey, error) {
	if len(der) != 91 {
		return nil, errors.New("public Key length invalid")
	}
	p := &PublicKey{Key: der}
	return p, nil
}

// MarshalSm2PublicKey encoding public key
func MarshalSm2PublicKey(key *PublicKey) ([]byte, error) {
	return key.Key, nil
}

// ParsePKCS8UnecryptedPrivateKey parse pkcs8 priv
func ParsePKCS8UnecryptedPrivateKey(der []byte) (*PrivateKey, error) {
	return ParseSm2PrivateKey(der)
}

// ParseSm2PrivateKey parse private key from pkcs8 encoded bytes
func ParseSm2PrivateKey(der []byte) (*PrivateKey, error) {
	if len(der) == 0 {
		return nil, errors.New("empty der block")
	}

	pubByte := make([]byte, 256)
	var publen C.int

	if C.GetPkFromPriv((*C.char)(unsafe.Pointer(&der[0])), C.int(len(der)), (*C.char)(unsafe.Pointer(&pubByte[0])), &publen) != 1 {
		return nil, errors.New("Failed to get pk from priv: " + OpenError())
	}

	pk := &PublicKey{
		Key: pubByte[:publen],
	}
	priv := &PrivateKey{
		PublicKey: pk,
		Key:       der,
	}

	return priv, nil
}

// MarshalSm2UnecryptedPrivateKey encodes private key
func MarshalSm2UnecryptedPrivateKey(key *PrivateKey) ([]byte, error) {
	return key.Key, nil
}

// GenerateKey return sm2 private key with public key
func GenerateKey() (*PrivateKey, error) {
	privByte := make([]byte, 256)
	pubByte := make([]byte, 256)
	var privlen C.int
	var publen C.int

	if C.X_GenerateKey((*C.char)(unsafe.Pointer(&privByte[0])), &privlen, (*C.char)(unsafe.Pointer(&pubByte[0])), &publen) != 1 {
		return nil, errors.New("Failed to generate key: " + OpenError())
	}

	pub := &PublicKey{
		Key: pubByte[:publen],
	}
	priv := &PrivateKey{
		PublicKey: pub,
		Key:       privByte[:privlen],
	}

	return priv, nil
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// CompressPubKey PC || X, Reference to GB/T 32918.5-2017
//
// y_bit is the lastbit of y, y_bit=0, then PC=02, y_bit=1, then PC=03
func CompressPubKey(key *PublicKey) ([]byte, error) {
	pub, err := OpensslPubToTjfoc(key)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 0, pubKeyBytesLenCompressed)

	pc := pubkeyCompressed
	if isOdd(pub.Y) {
		pc |= 0x1
	}
	b = append(b, pc)

	return paddedAppend(32, b, pub.X.Bytes()), nil
}

// SerializeUncompressed serializes a public key in a 65-byte uncompressed
// format.
func (pub *PublicKey) SerializeUncompressed() ([]byte, error) {
	p, err := OpensslPubToTjfoc(pub)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 0, pubKeyBytesLenUncompressed)
	b = append(b, pubkeyUncompressed)
	b = paddedAppend(32, b, p.X.Bytes())
	return paddedAppend(32, b, p.Y.Bytes()), nil
}

// TjfocPubToOpenssl convert tjfoc pubkey to openssl sm2 pubkey
func TjfocPubToOpenssl(key *sm2.PublicKey) (*PublicKey, error) {
	data, err := sm2.MarshalSm2PublicKey(key)
	if err != nil {
		return nil, err
	}
	pub, err := ParseSm2PublicKey(data)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

// OpensslPubToTjfoc convert openssl sm2 pubkey to tjfoc pubkey
func OpensslPubToTjfoc(key *PublicKey) (*sm2.PublicKey, error) {
	data, err := MarshalSm2PublicKey(key)
	if err != nil {
		return nil, err
	}
	pub, err := sm2.ParseSm2PublicKey(data)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

// ParseUncompressedPubKey parse uncompressed publickey(04 + X + Y format)
func ParseUncompressedPubKey(data []byte) (*PublicKey, error) {
	if l := len(data); l != pubKeyBytesLenUncompressed {
		return nil, fmt.Errorf("invalid uncompressed pubkey length %d", l)
	}
	format := data[0]
	if format != pubkeyUncompressed {
		return nil, fmt.Errorf("invalid magic in pubkey str: 0x%x", format)
	}
	pk := &sm2.PublicKey{
		Curve: sm2.P256Sm2(),
		X:     new(big.Int).SetBytes(data[1:33]),
		Y:     new(big.Int).SetBytes(data[33:]),
	}
	return TjfocPubToOpenssl(pk)
}

// DecompressPubKey decompress compressed publickey(PC + X format)
func DecompressPubKey(data []byte) (*PublicKey, error) {
	if len(data) != pubKeyBytesLenCompressed {
		return nil, fmt.Errorf("invalid compressed pubkey length %d, want: %d", len(data), pubKeyBytesLenCompressed)
	}
	format := data[0]
	// if is odd
	ybit := (format & 0x1) == 0x1
	format &= ^byte(0x1)
	if format != pubkeyCompressed {
		return nil, fmt.Errorf("invalid magic in compressed pubkey string %d", data[0])
	}

	pubStr := make([]byte, 0, pubKeyBytesLenCompressed)
	pubStr = append(pubStr, data...)
	pubStr[0] = pubStr[0] - 2

	tjpub := sm2.Decompress(pubStr)

	// verify y-coord hash expected parity
	if ybit != isOdd(tjpub.Y) {
		return nil, fmt.Errorf("ybit doesn't match oddness")
	}
	return TjfocPubToOpenssl(tjpub)
}
