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
	"encoding/base64"
	"encoding/json"
	"madledger/common/crypto"
	"madledger/core"
)

// CreateGenesisBlock return the genesis block
// TODO: maybe there should includes some admins in the genesis block
func CreateGenesisBlock(admins []*core.Member) (*core.Block, error) {
	var payloads = []Payload{{
		ChannelID: core.CONFIGCHANNELID,
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}, Payload{
		ChannelID: core.GLOBALCHANNELID,
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}, Payload{
		ChannelID: core.ASSETCHANNELID,
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}, Payload{ // this payload is used to record the info of  system admin
		// todo: modify here, choose a better way to record it, the nil channelID will confuse the peer/orderer
		Profile: &Profile{
			Public: true,
			Admins: admins,
		},
		Version: 1,
	}}
	var txs []*core.Tx
	for i, payload := range payloads {
		payloadBytes, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}

		accountNonce := uint64(i)
		tx := core.NewTxWithoutSig(core.CONFIGCHANNELID, payloadBytes, accountNonce)
		txs = append(txs, tx)
	}

	return core.NewBlock(core.CONFIGCHANNELID, 0, core.GenesisBlockPrevHash, txs), nil
}

// CreateAdmins create admins
// TODO: Hard code here
// TODO: Remove it
func CreateAdmins() ([]*core.Member, error) {
	// get pubkey from string by base64 encoding
	data, err := base64.StdEncoding.DecodeString("BGXcjZ3bhemsoLP4HgBwnQ5gsc8VM91b3y8bW0b6knkWu8x" +
		"CSKO2qiJXARMHcbtZtvU7Jos2A5kFCD1haJ/hLdg=")
	if err != nil {
		return nil, err
	}
	// create PublicKey
	// TODO: Should not do that
	pk, err := crypto.NewPublicKey(data, crypto.KeyAlgoSecp256k1)
	if err != nil {
		return nil, err
	}
	// create member and append it to the admin
	member, err := core.NewMember(pk, "SystemAdmin")
	admins := make([]*core.Member, 0)
	admins = append(admins, member)
	return admins, nil
}
