package config

import (
	"madledger/common/crypto"
	"madledger/core/types"
	"testing"
)

var (
	admin    = newMember("admin")
	civilian = newMember("civilian")
	criminal = newMember("criminal")
)

func TestPublicPayload(t *testing.T) {
	// illegal channelID
	payload := Payload{
		ChannelID: "",
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}
	if payload.Verify() {
		t.Fatal()
	}
	// legal channelID
	payload = Payload{
		ChannelID: "public",
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}
	if !payload.Verify() {
		t.Fatal()
	}
	if payload.IsAdmin(admin) {
		t.Fatal()
	}
	if !payload.IsMember(admin) {
		t.Fatal()
	}
	if !payload.IsMember(civilian) {
		t.Fatal()
	}
	// then set admin
	payload.Profile.Admins = []*types.Member{admin}
	if !payload.Verify() {
		t.Fatal()
	}
	if !payload.IsAdmin(admin) {
		t.Fatal()
	}
	if !payload.IsMember(admin) {
		t.Fatal()
	}
	if !payload.IsMember(civilian) {
		t.Fatal()
	}
}

func TestPrivatePayload(t *testing.T) {
	// without members and admins
	payload := Payload{
		ChannelID: "private",
		Profile: &Profile{
			Public: false,
		},
		Version: 1,
	}
	if payload.Verify() {
		t.Fatal()
	}
	// with members
	payload = Payload{
		ChannelID: "private",
		Profile: &Profile{
			Public:  false,
			Members: []*types.Member{admin},
		},
		Version: 1,
	}
	if !payload.Verify() {
		t.Fatal()
	}
	if payload.IsAdmin(admin) || !payload.IsMember(admin) {
		t.Fatal()
	}
	if payload.IsMember(civilian) || payload.IsAdmin(civilian) {
		t.Fatal()
	}
	// with members and admins
	payload = Payload{
		ChannelID: "private",
		Profile: &Profile{
			Public:  false,
			Members: []*types.Member{civilian, admin},
			Admins:  []*types.Member{admin},
		},
		Version: 1,
	}
	if !payload.Verify() {
		t.Fatal()
	}
	if !payload.IsAdmin(admin) || !payload.IsMember(admin) {
		t.Fatal()
	}
	if !payload.IsMember(civilian) || payload.IsAdmin(civilian) {
		t.Fatal()
	}
	if payload.IsMember(criminal) || payload.IsAdmin(criminal) {
		t.Fatal()
	}
	// with members and admins, but admins are not contained in the members
	payload = Payload{
		ChannelID: "private",
		Profile: &Profile{
			Public:  false,
			Members: []*types.Member{civilian, criminal},
			Admins:  []*types.Member{admin},
		},
		Version: 1,
	}
	if payload.Verify() {
		t.Fatal()
	}
}

func newMember(name string) *types.Member {
	privKey, err := crypto.GeneratePrivateKey()
	if err != nil {
		return nil
	}
	pk := privKey.PubKey()
	return types.NewMember(pk, name)
}
