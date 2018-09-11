package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"madledger/common"
	"madledger/common/crypto/secp256k1"
	"madledger/common/math"
)

// The secp256k1 implementation of sign.go

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

// SECP256K1PrivateKey defines the secp256k1 private key in ecdsa
type SECP256K1PrivateKey ecdsa.PrivateKey

// SECP256K1PublicKey defines the secp256k1 public key in ecdsa
type SECP256K1PublicKey ecdsa.PublicKey

// SECP256K1SignatureLength defines the length of SECP256K1Signature
const SECP256K1SignatureLength = 64

// SECP256K1Signature defines the secp256k1 signature
type SECP256K1Signature [SECP256K1SignatureLength]byte

// HexToSECP256K1PrivateKey convert a hex string to PrivateKey
func HexToSECP256K1PrivateKey(h string) (PrivateKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return toSECP256K1PrivateKey(b)
}

// GenerateSECP256K1PrivateKey return a new secp256k1 private key
func GenerateSECP256K1PrivateKey() (PrivateKey, error) {
	ecdsaPrivateKey, err := ecdsa.GenerateKey(S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return (SECP256K1PrivateKey)(*ecdsaPrivateKey), nil
}

// PubKey returns the PublicKey corresponding to this private key.
func (p SECP256K1PrivateKey) PubKey() PublicKey {
	pk := (SECP256K1PublicKey)((ecdsa.PrivateKey)(p).PublicKey)
	return (PublicKey)(pk)
}

// Bytes is the implementation of interface
func (p SECP256K1PrivateKey) Bytes() ([]byte, error) {
	return math.PaddedBigBytes(p.D, p.Params().BitSize/8), nil
}

// Bytes returns the bytes of Public key
func (p SECP256K1PublicKey) Bytes() ([]byte, error) {
	pk := (*(ecdsa.PublicKey))(&p)
	if pk == nil || pk.X == nil || pk.Y == nil {
		return nil, errors.New("The public key is not formatted")
	}
	return elliptic.Marshal(S256(), pk.X, pk.Y), nil
}

// Address is the implementation of interface
func (p SECP256K1PublicKey) Address() (common.Address, error) {
	bytes, err := p.Bytes()
	if err != nil {
		return common.ZeroAddress, nil
	}
	hash := Keccak256Hash(bytes[1:])
	addrBytes := hash.Bytes()[12:]
	return common.AddressFromBytes(addrBytes)
}

// Sign sign the data using the privateKey
func (p SECP256K1PrivateKey) Sign(hash []byte) (Signature, error) {
	var prv = (*ecdsa.PrivateKey)(&p)
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
	defer zeroBytes(seckey)
	sig, err := secp256k1.Sign(hash, seckey)
	if err != nil {
		return nil, err
	}
	return newSECP256K1Signature(sig)
}

// Verify is the implementation of interface
func (s SECP256K1Signature) Verify(hash []byte, pubKey PublicKey) bool {
	switch pubKey.(type) {
	case SECP256K1PublicKey:
	default:
		return false
	}
	keyBytes, err := pubKey.Bytes()
	if err != nil {
		return false
	}
	return secp256k1.VerifySignature(keyBytes, hash, s[:])
}

// Bytes is the implementation of interface
func (s SECP256K1Signature) Bytes() ([]byte, error) {
	return s[:], nil
}

func newSECP256K1PublicKey(raw []byte) (PublicKey, error) {
	x, y := elliptic.Unmarshal(S256(), raw)
	if x == nil {
		return nil, errors.New("Marshal failed")
	}
	return (SECP256K1PublicKey)(ecdsa.PublicKey{Curve: S256(), X: x, Y: y}), nil
}

func newSECP256K1Signature(raw []byte) (Signature, error) {
	var sig SECP256K1Signature
	copy(sig[:], raw)
	return sig, nil
}

func toSECP256K1PrivateKey(bs []byte) (PrivateKey, error) {
	priv, err := toECDSA(bs, false)
	// return priv
	if err != nil {
		return nil, err
	}
	return (SECP256K1PrivateKey)(*priv), nil
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
