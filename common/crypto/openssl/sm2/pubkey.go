package sm2

import (
	"errors"
)

// PublicKey is the public key of sm2
type PublicKey struct {
	Key []byte
}

// // Bytes returns pubkey marshaled with the specific format
// // Note: It return uncompressed bytes
// func (pub *PublicKey) Bytes() ([]byte, error) {
// 	// switch opt {
// 	// case crypto.Compressed:
// 	// 	return CompressPubKey(pub)
// 	// case crypto.ASN1Uncompressed:
// 	// 	return MarshalSm2PublicKey(pub)
// 	// case crypto.Uncompressed:
// 	// 	return pub.SerializeUncompressed()
// 	// default:
// 	// 	return nil, fmt.Errorf("unknown compress opt 0x%x", opt)
// 	// }
// 	return pub.SerializeUncompressed()
// }

// NewPublicKey parse sm2 public key from asn.1 or compressed format serialized pk
func NewPublicKey(raw []byte) (*PublicKey, error) {
	if len(raw) == 0 {
		return nil, errors.New("nil pk raw")
	}
	switch raw[0] {
	case pubkeyCompressed, pubkeyCompressed + 1:
		return DecompressPubKey(raw)
	case pubkeyUncompressed:
		return ParseUncompressedPubKey(raw)
	case pubkeyASN1:
		return ParseSm2PublicKey(raw)
	default:
		return nil, errors.New("sm2: invalid prefix")
	}
}
