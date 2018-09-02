package crypto

import (
	"madledger/common"
	"math/big"
)

// PrivateKey is the interface of privateKey
// It may support ecdsa or sm2
type PrivateKey interface {
	Sign(hash []byte) (Signature, error)
	PubKey() PublicKey
}

// PublicKey is the interface of publicKey
// It may support ecdsa or sm2
type PublicKey interface {
	Bytes() ([]byte, error)
	// GetSerializeLength() int
	Address() (common.Address, error)
}

// Signature interface is the interface of signature
// It may support ecdsa or sm2
type Signature interface {
	Verify(hash []byte, pubKey PublicKey) bool
	Bytes() ([]byte, error)
}

// NewPrivateKey return a PrivateKey
// Only support ECDSAPrivateKey yet
func NewPrivateKey(raw []byte) (PrivateKey, error) {
	return toSECP256K1PrivateKey(raw)
	// first try to parse
	// ecPrivKey, err := x509.ParseECPrivateKey(raw)
	// if err == nil {
	// 	return (PrivateKey)((*ECDSAPrivateKey)(ecPrivKey)), nil
	// }
	// return nil, err
}

// NewPublicKey return a PublicKey from []byte
func NewPublicKey(raw []byte) (PublicKey, error) {
	// return parseECDSAPublicKey(raw)
	return newSECP256K1PublicKey(raw)
}

// NewSignature return a signature from []byte
func NewSignature(raw []byte) (Signature, error) {
	// return parseECDSASignature(raw)
	return newSECP256K1Signature(raw)
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
