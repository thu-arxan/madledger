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
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"madledger/common"
)

// Algorithm identifies algorithm for asymmetric cryptography
type Algorithm = int32

// These defines asymmetric cryptography algorithm
const (
	KeyAlgoSM2 Algorithm = iota
	KeyAlgoSecp256k1
	// KeyAlgoED25519
)

// PrivateKey is the interface of privateKey
// It may support secp256k1 or sm2
type PrivateKey interface {
	Sign(hash []byte) (Signature, error)
	PubKey() PublicKey
	Bytes() ([]byte, error)
	Algo() Algorithm
}

// PublicKey is the interface of publicKey
// It  support secp256k1 or sm2
type PublicKey interface {
	Bytes() ([]byte, error)
	Address() (common.Address, error)
	Algo() Algorithm
}

// GeneratePrivateKey try to generate a private key
func GeneratePrivateKey(algo ...Algorithm) (PrivateKey, error) {
	if len(algo) != 0 {
		switch algo[0] {
		case KeyAlgoSecp256k1:
			return GenerateSECP256K1PrivateKey()
		default:
			return GenerateSM2PrivateKey()
		}
	}
	return GenerateSM2PrivateKey()
}

// LoadPrivateKeyFromFile load private key from file
func LoadPrivateKeyFromFile(file string) (PrivateKey, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// Note: secp256k1 private key we store on disk using hex string, so we can distinguish these two keys
	key, err := hex.DecodeString(string(data))
	if err != nil {
		return NewPrivateKey(data, KeyAlgoSM2)
	}
	return NewPrivateKey(key, KeyAlgoSecp256k1)
}

// NewPrivateKey return a PrivateKey
// Support secp256k1 and sm2
func NewPrivateKey(raw []byte, algo Algorithm) (PrivateKey, error) {
	switch algo {
	case KeyAlgoSecp256k1:
		return toSECP256K1PrivateKey(raw)
	case KeyAlgoSM2:
		return toSM2PrivateKey(raw)
	default:
		return nil, fmt.Errorf("unsupport algo:%v", algo)
	}
}

// NewPublicKey return a PublicKey from []byte
// Support secp256k1 and sm2
func NewPublicKey(raw []byte, algo ...Algorithm) (PublicKey, error) {
	if len(algo) != 0 {
		switch algo[0] {
		case KeyAlgoSecp256k1:
			return newSECP256K1PublicKey(raw)
		default:
			return newSM2PublicKey(raw)
		}
	}
	return newSM2PublicKey(raw)
}
