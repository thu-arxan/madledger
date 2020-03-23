package sm2

// PrivateKey is the private key of sm2
type PrivateKey struct {
	*PublicKey
	Key []byte
}

// Public returns the public key corresponding to the private key
func (priv *PrivateKey) Public() *PublicKey {
	return priv.PublicKey
}

// Bytes returns asn.1 serialized private key
func (priv *PrivateKey) Bytes() ([]byte, error) {
	return MarshalSm2UnecryptedPrivateKey(priv)
}

// Sign signs digest(no need hash) with the private key and default uid
func (priv *PrivateKey) Sign(digest []byte) ([]byte, error) {
	hash := digest
	data, err := Sign(priv, hash, []byte(DefaultUID))
	if err != nil {
		return nil, err
	}
	return data, nil
}
