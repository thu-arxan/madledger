// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package crypto

import (
	"errors"
	"math/big"
)

// Signature interface is the interface of signature
// It may support ecdsa or sm2
type Signature interface {
	Verify(hash []byte, pubKey PublicKey) bool
	Bytes() ([]byte, error)
}

// NewSignature return a signature from []byte
func NewSignature(raw []byte, algo Algorithm) (Signature, error) {
	// return parseECDSASignature(raw)
	switch algo {
	case KeyAlgoSM2:
		return newSM2Signature(raw)
	case KeyAlgoSecp256k1:
		return newSECP256K1Signature(raw)
	default:
		return nil, errors.New("unsupport algo")
	}
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
