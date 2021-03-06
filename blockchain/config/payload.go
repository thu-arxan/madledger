// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package config

import (
	"madledger/common/util"
	"madledger/core"
)

// Payload config a channel
type Payload struct {
	ChannelID string
	Profile   *Profile
	Version   int32
}

// Profile is the profile in payload
type Profile struct {
	// Public is true means that everyone could access to the channel and
	// ignores the Members. But the Admins still works.
	Public bool
	// Dependencies includes all channels that the channel relies on.
	Dependencies []string
	// Members
	Members []*core.Member
	// Admins
	// Note: If the public is true, Admins is still works and may not be contained in the
	// Members. But if the public is false, Admins should be contained in the Members.
	Admins []*core.Member

	// GasPrice is the token needed per gas
	GasPrice uint64
	// number of tokens that one asset exchange
	AssetTokenRatio uint64
	// Maximum Gas spent in one evm execution
	MaxGas uint64

	// Block Storage Price
	BlockPrice uint64
}

// Verify returns if a payload is packed well
func (payload *Payload) Verify() bool {
	// verify ChannelID
	switch payload.ChannelID {
	case core.GLOBALCHANNELID:
	case core.CONFIGCHANNELID:
	case core.ASSETCHANNELID:
	default:
		if !util.IsLegalChannelName(payload.ChannelID) {
			return false
		}
	}

	if !payload.Profile.Public {
		if payload.Profile.Members == nil || len(payload.Profile.Members) == 0 {
			return false
		}
		for _, admin := range payload.Profile.Admins {
			if !payload.IsMember(admin) {
				return false
			}
		}
	}
	return true
}

// IsMember return if the member is contained in the channel
func (payload *Payload) IsMember(member *core.Member) bool {
	// If the channel is public, then returns true
	if payload.Profile.Public {
		return true
	}
	for _, m := range payload.Profile.Members {
		if member.Equal(m) {
			return true
		}
	}
	return false
}

// IsAdmin return if the member is the member of the channel
func (payload *Payload) IsAdmin(member *core.Member) bool {
	for _, m := range payload.Profile.Admins {
		if member.Equal(m) {
			return true
		}
	}
	return false
}
