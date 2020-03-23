package sm2

import (
	"errors"
	"strconv"
)

// Signature sm2 signature signed with uid
type Signature struct {
	sig []byte
}

// Verify verifies the signature with the public key and default uid, digest is the origin msg
func (sig *Signature) Verify(digest []byte, pub *PublicKey) bool {
	return Verify(pub, digest, []byte(DefaultUID), sig.sig)
}

// Bytes returns the der-serialized, ASN.1 format signature
func (sig *Signature) Bytes() []byte {
	data := make([]byte, 0)
	data = append(data, sig.sig...)
	return data
}

// NewSignature parse signature from raw sig(asn.1 sig or prefix+sig)
func NewSignature(raw []byte) (*Signature, error) {
	l := len(raw)
	// 70~72, with prefix
	if l < minSigLen || l > maxSigLen+1 {
		return nil, errors.New("sm2: bad signature len: " + strconv.Itoa(l))
	}
	sig := &Signature{}
	data := make([]byte, 0, l)
	data = append(data, raw...)
	sig.sig = data
	return sig, nil
}
