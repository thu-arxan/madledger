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
	"madledger/common"
	"madledger/common/crypto/hash"
	"madledger/common/crypto/openssl/sm2"
)

// SM2PrivateKey defines the sm2 private key
type SM2PrivateKey sm2.PrivateKey

// SM2PublicKey defines the sm2 public key
type SM2PublicKey sm2.PublicKey

// SM2Signature defines the sm2 signature
type SM2Signature sm2.Signature

// GenerateSM2PrivateKey return a new sesm2cp256k1 private key
func GenerateSM2PrivateKey() (PrivateKey, error) {
	privKey, err := sm2.GenerateKey()
	if err != nil {
		return nil, err
	}
	return (SM2PrivateKey)(*privKey), nil
}

// PubKey returns the PublicKey corresponding to this private key.
func (p SM2PrivateKey) PubKey() PublicKey {
	var privKey = (sm2.PrivateKey)(p)
	var pubKey = privKey.Public()
	return (SM2PublicKey)(*pubKey)
}

// Bytes is the implementation of interface
func (p SM2PrivateKey) Bytes() []byte {
	var privKey = (sm2.PrivateKey)(p)
	bs, _ := sm2.WritePrivateKeytoMem(&privKey, nil)
	return bs
}

// Sign sign the data using the privateKey
func (p SM2PrivateKey) Sign(hash []byte) (Signature, error) {
	var privKey = (*sm2.PrivateKey)(&p)
	sigData, err := privKey.Sign(hash)
	if err != nil {
		return nil, err
	}
	sig, err := sm2.NewSignature(sigData)
	if err != nil {
		return nil, err
	}
	return (SM2Signature)(*sig), nil
}

// Algo return Algo of private key
func (p SM2PrivateKey) Algo() Algorithm {
	return KeyAlgoSM2
}

// Bytes returns the bytes of Public key
func (p SM2PublicKey) Bytes() ([]byte, error) {
	var pubKey = (sm2.PublicKey)(p)
	// return sm2.CompressPubKey(&pubKey)
	// todo: fix the bug that pk.Bytes may panic
	// defer return error ???
	return pubKey.SerializeUncompressed()
}

// Address is the implementation of interface
func (p SM2PublicKey) Address() (common.Address, error) {
	var pubKey = (sm2.PublicKey)(p)
	bytes, err := pubKey.SerializeUncompressed()
	if err != nil {
		return common.ZeroAddress, nil
	}

	hash := hash.SM3(bytes[1:])
	addrBytes := hash[12:]
	return common.AddressFromBytes(addrBytes)
}

// Algo return Algo of private key
func (p SM2PublicKey) Algo() Algorithm {
	return KeyAlgoSM2
}

// Verify is the implementation of interface
func (s SM2Signature) Verify(hash []byte, pubKey PublicKey) bool {
	var sig = (sm2.Signature)(s)
	switch pubKey.(type) {
	case SM2PublicKey:
		var pk = (sm2.PublicKey)((pubKey).(SM2PublicKey))
		return sig.Verify(hash, &pk)
	default:
		return false
	}
}

// Bytes is the implementation of interface
func (s SM2Signature) Bytes() ([]byte, error) {
	var sig = (sm2.Signature)(s)
	return sig.Bytes(), nil
}

func newSM2PublicKey(raw []byte) (PublicKey, error) {
	pk, err := sm2.NewPublicKey(raw)
	if err != nil {
		return nil, err
	}
	return (SM2PublicKey)(*pk), nil
}

func newSM2Signature(raw []byte) (Signature, error) {
	sig, err := sm2.NewSignature(raw)
	if err != nil {
		return nil, err
	}
	return (SM2Signature)(*sig), nil
}

func toSM2PrivateKey(bs []byte) (PrivateKey, error) {
	priv, err := sm2.ReadPrivateKeyFromMem(bs, nil)
	if err != nil {
		return nil, err
	}
	return (SM2PrivateKey)(*priv), nil
}
