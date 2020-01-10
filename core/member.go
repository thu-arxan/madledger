package core

import (
	"madledger/common/crypto"
	"reflect"
)

// Member is the member in the system
type Member struct {
	// PK is the public key, which defines what the member is
	PK []byte
	// Name is used to present the member
	Name string
}

// NewMember is the constructor of Member
func NewMember(pk crypto.PublicKey, name string) (*Member, error) {
	pkBytes, err := pk.Bytes()
	if err != nil {
		return nil, err
	}
	return &Member{
		PK:   pkBytes,
		Name: name,
	}, nil
}

// Equal return if the two member is same
func (m *Member) Equal(m1 *Member) bool {
	return reflect.DeepEqual(m.PK, m1.PK)
}
