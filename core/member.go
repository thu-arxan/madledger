// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package core

import (
	"madledger/common/crypto"
	"reflect"
)

// Member is the member in the system
type Member struct {
	// PK is the public key, which defines what the member is
	PK   []byte
	Algo crypto.Algorithm
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
		Algo: pk.Algo(),
		Name: name,
	}, nil
}

// Equal return if the two member is same
func (m *Member) Equal(m1 *Member) bool {
	return reflect.DeepEqual(m.PK, m1.PK)
}
