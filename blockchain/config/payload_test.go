package config

import (
	"madledger/common/crypto"
	"madledger/core"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.Equal(t, payload.Verify(), false)

	// legal channelID
	payload = Payload{
		ChannelID: "public",
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}
	require.Equal(t, payload.Verify(), true)
	require.Equal(t, payload.IsAdmin(admin), false)
	require.Equal(t, payload.IsMember(admin), true)
	require.Equal(t, payload.IsMember(civilian), true)
	// then set admin
	payload.Profile.Admins = []*core.Member{admin}
	require.Equal(t, payload.Verify(), true)
	require.Equal(t, payload.IsAdmin(admin), true)
	require.Equal(t, payload.IsMember(admin), true)
	require.Equal(t, payload.IsMember(civilian), true)
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
	require.Equal(t, payload.Verify(), false)
	// with members
	payload = Payload{
		ChannelID: "private",
		Profile: &Profile{
			Public:  false,
			Members: []*core.Member{admin},
		},
		Version: 1,
	}
	require.Equal(t, payload.Verify(), true)
	require.Equal(t, payload.IsAdmin(admin), false)
	require.Equal(t, payload.IsMember(admin), true)
	require.Equal(t, payload.IsMember(civilian) || payload.IsAdmin(civilian), false)
	// with members and admins
	payload = Payload{
		ChannelID: "private",
		Profile: &Profile{
			Public:  false,
			Members: []*core.Member{civilian, admin},
			Admins:  []*core.Member{admin},
		},
		Version: 1,
	}
	require.Equal(t, payload.Verify(), true)
	require.Equal(t, !payload.IsAdmin(admin) || !payload.IsMember(admin), false)
	require.Equal(t, !payload.IsMember(civilian) || payload.IsAdmin(civilian), false)
	require.Equal(t, payload.IsMember(criminal) || payload.IsAdmin(criminal), false)
	// with members and admins, but admins are not contained in the members
	payload = Payload{
		ChannelID: "private",
		Profile: &Profile{
			Public:  false,
			Members: []*core.Member{civilian, criminal},
			Admins:  []*core.Member{admin},
		},
		Version: 1,
	}
	require.Equal(t, payload.Verify(), false)
}

func newMember(name string) *core.Member {
	privKey, err := crypto.GeneratePrivateKey()
	if err != nil {
		return nil
	}
	pk := privKey.PubKey()
	member, _ := core.NewMember(pk, name)
	return member
}
