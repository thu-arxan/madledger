package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/asn1"
	"errors"
	"fmt"
	"madledger/common"
	"math/big"
)

// The ecdsa implementation of sign.go

// These constants define the lengths of serialized public keys.
const (
	PubKeyBytesLenUncompressed = 65
)

const (
	pubkeyUncompressed byte = 0x4 // x coord + y coord
)

// ECDSASignature is a type representing an ecdsa signature.
type ECDSASignature struct {
	R *big.Int
	S *big.Int
}

// ECDSAPrivateKey wraps an ecdsa.PrivateKey as a convenience
// It should change to domestic encryption algorithm when needed
type ECDSAPrivateKey ecdsa.PrivateKey

// ECDSAPublicKey is an ecdsa.PublicKey as a convenience
// It should change to domestic encryption algorithm when needed
type ECDSAPublicKey ecdsa.PublicKey

// Sign sign the data using the privateKey
func (p ECDSAPrivateKey) Sign(hash []byte) (Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, (*ecdsa.PrivateKey)(&p), hash)
	if err != nil {
		return nil, err
	}
	raw, err := asn1.Marshal(ECDSASignature{r, s})
	if err != nil {
		return nil, err
	}
	return NewSignature(raw)
}

// PubKey returns the PublicKey corresponding to this private key.
func (p ECDSAPrivateKey) PubKey() PublicKey {
	pk := (ECDSAPublicKey)((ecdsa.PrivateKey)(p).PublicKey)
	return (PublicKey)(pk)
}

// Bytes is the implementation of PublicKey interface
func (p ECDSAPublicKey) Bytes() ([]byte, error) {
	bytes := p.SerializeUncompressed()
	return bytes, nil
	// return p.Bytes(), nil
}

// Address is not implementation yet
func (p ECDSAPublicKey) Address() (common.Address, error) {
	return common.ZeroAddress, errors.New("Not implementation yet")
}

// SerializeUncompressed serializes a public key in a 65-byte uncompressed format.
func (p ECDSAPublicKey) SerializeUncompressed() []byte {
	b := make([]byte, 0, PubKeyBytesLenUncompressed)
	b = append(b, pubkeyUncompressed)
	b = paddedAppend(32, b, p.X.Bytes())
	return paddedAppend(32, b, p.Y.Bytes())
}

// parseECDSAPublicKey parses a public key into a ecdsa.Publickey, verifying that it is valid.
// It only supports uncompressed yet.
func parseECDSAPublicKey(pubKeyStr []byte) (key PublicKey, err error) {
	pubkey := ECDSAPublicKey{}
	pubkey.Curve = elliptic.P256()

	if len(pubKeyStr) == 0 {
		return nil, errors.New("pubkey string is empty")
	}

	format := pubKeyStr[0]
	// ybit := (format & 0x1) == 0x1
	format &= ^byte(0x1)

	switch len(pubKeyStr) {
	case PubKeyBytesLenUncompressed:
		if format != pubkeyUncompressed {
			return nil, fmt.Errorf("invalid magic in pubkey str: "+"%d", pubKeyStr[0])
		}

		pubkey.X = new(big.Int).SetBytes(pubKeyStr[1:33])
		pubkey.Y = new(big.Int).SetBytes(pubKeyStr[33:])
	default: // wrong!
		return nil, fmt.Errorf("invalid pub key length %d",
			len(pubKeyStr))
	}

	if pubkey.X.Cmp(pubkey.Curve.Params().P) >= 0 {
		return nil, fmt.Errorf("pubkey X parameter is >= to P")
	}
	if pubkey.Y.Cmp(pubkey.Curve.Params().P) >= 0 {
		return nil, fmt.Errorf("pubkey Y parameter is >= to P")
	}
	if !pubkey.Curve.IsOnCurve(pubkey.X, pubkey.Y) {
		return nil, fmt.Errorf("pubkey isn't on secp256k1 curve")
	}
	return (PublicKey)(pubkey), nil
}

// parseECDSASignature return the signature of ECDSA
func parseECDSASignature(raw []byte) (Signature, error) {
	// Unmarshal
	sig := new(ECDSASignature)
	_, err := asn1.Unmarshal(raw, sig)
	if err != nil {
		return nil, fmt.Errorf("Failed unmashalling signature [%s]", err)
	}

	// Validate sig
	if sig.R == nil {
		return nil, errors.New("Invalid signature. R must be different from nil")
	}
	if sig.S == nil {
		return nil, errors.New("Invalid signature. S must be different from nil")
	}

	if sig.R.Sign() != 1 {
		return nil, errors.New("Invalid signature. R must be larger than zero")
	}
	if sig.S.Sign() != 1 {
		return nil, errors.New("Invalid signature. S must be larger than zero")
	}

	return (Signature)(sig), nil
}

// Verify is the implementation of interface
func (s ECDSASignature) Verify(hash []byte, pubKey PublicKey) bool {
	var key ECDSAPublicKey
	switch pubKey.(type) {
	case ECDSAPublicKey:
		key = pubKey.(ECDSAPublicKey)
	default:
		return false
	}
	return ecdsa.Verify((*ecdsa.PublicKey)(&key), hash, s.R, s.S)
}
