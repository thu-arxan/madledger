package crypto

import (
	"madledger/common"
	"madledger/common/crypto/sha3"

	"github.com/decred/dcrd/dcrec/secp256k1"
)

// The secp256k1 implementation of sign.go
// If need some help about secp256k1, visit https://godoc.org/github.com/decred/dcrd/dcrec/secp256k1

// SECP256K1PrivateKey defines the secp256k1 private key in ecdsa
type SECP256K1PrivateKey secp256k1.PrivateKey

// SECP256K1PublicKey defines the secp256k1 public key in ecdsa
type SECP256K1PublicKey secp256k1.PublicKey

// SECP256K1Signature defines the secp256k1 signature
type SECP256K1Signature secp256k1.Signature

// GenerateSECP256K1PrivateKey return a new secp256k1 private key
func GenerateSECP256K1PrivateKey() (PrivateKey, error) {
	privKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	return (SECP256K1PrivateKey)(*privKey), nil
}

// PubKey returns the PublicKey corresponding to this private key.
func (p SECP256K1PrivateKey) PubKey() PublicKey {
	var privKey = (secp256k1.PrivateKey)(p)
	var pubKey = privKey.PubKey()
	return (SECP256K1PublicKey)(*pubKey)
}

// Bytes is the implementation of interface
func (p SECP256K1PrivateKey) Bytes() ([]byte, error) {
	var privKey = (secp256k1.PrivateKey)(p)
	return privKey.Serialize(), nil
}

// Sign sign the data using the privateKey
func (p SECP256K1PrivateKey) Sign(hash []byte) (Signature, error) {
	var privKey = (*secp256k1.PrivateKey)(&p)
	sig, err := privKey.Sign(hash)
	if err != nil {
		return nil, err
	}
	return (SECP256K1Signature)(*sig), nil
}

// Bytes returns the bytes of Public key
func (p SECP256K1PublicKey) Bytes() ([]byte, error) {
	var pubKey = (secp256k1.PublicKey)(p)
	return pubKey.SerializeUncompressed(), nil
}

// Address is the implementation of interface
func (p SECP256K1PublicKey) Address() (common.Address, error) {
	bytes, err := p.Bytes()
	if err != nil {
		return common.ZeroAddress, nil
	}
	hash := sha3.Sha3(bytes[1:])
	addrBytes := hash[12:]
	return common.AddressFromBytes(addrBytes)
}

// Verify is the implementation of interface
func (s SECP256K1Signature) Verify(hash []byte, pubKey PublicKey) bool {
	var sig = (secp256k1.Signature)(s)
	switch pubKey.(type) {
	case SECP256K1PublicKey:
		var pk = (secp256k1.PublicKey)((pubKey).(SECP256K1PublicKey))
		return sig.Verify(hash, &pk)
	default:
		return false
	}
}

// Bytes is the implementation of interface
func (s SECP256K1Signature) Bytes() ([]byte, error) {
	var sig = (secp256k1.Signature)(s)
	return sig.Serialize(), nil
}

func newSECP256K1PublicKey(raw []byte) (PublicKey, error) {
	pk, err := secp256k1.ParsePubKey(raw)
	if err != nil {
		return nil, err
	}
	return (SECP256K1PublicKey)(*pk), nil
}

func newSECP256K1Signature(raw []byte) (Signature, error) {
	sig, err := secp256k1.ParseSignature(raw)
	if err != nil {
		return nil, err
	}
	return (SECP256K1Signature)(*sig), nil
}

func toSECP256K1PrivateKey(bs []byte) (PrivateKey, error) {
	priv, _ := secp256k1.PrivKeyFromBytes(bs)
	return (SECP256K1PrivateKey)(*priv), nil
}
