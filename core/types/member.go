package types

import (
	"madledger/common/crypto"
	"reflect"
)

// Member is the member in the system
type Member struct {
	// PK is the public key, which defines what the member is
	PK crypto.PublicKey
	// Name is used to present the member
	Name string
}

// NewMember is the constructor of Member
func NewMember(pk crypto.PublicKey, name string) *Member {
	return &Member{
		PK:   pk,
		Name: name,
	}
}

// Equal return if the two member is same
func (m *Member) Equal(m1 *Member) bool {
	pk, err := m.PK.Bytes()
	if err != nil {
		return false
	}
	pk1, err := m1.PK.Bytes()
	if err != nil {
		return false
	}
	return reflect.DeepEqual(pk, pk1)
}
