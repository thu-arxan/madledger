package config

import (
	"encoding/base64"
	"encoding/json"
	"madledger/core/types"
	"madledger/common/crypto"
)

// CreateGenesisBlock return the genesis block
// TODO: maybe there should includes some admins in the genesis block
func CreateGenesisBlock(admins []*types.Member) (*types.Block, error) {
	var payloads = []Payload{{
		ChannelID: types.CONFIGCHANNELID,
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}, Payload{
		ChannelID: types.GLOBALCHANNELID,
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}, Payload{ // this payload is used to record the info of  system admin
		Profile: &Profile{
			Public: true,
			Admins: admins,
		},
		Version: 1,
	}}
	var txs []*types.Tx
	for i, payload := range payloads {
		payloadBytes, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}

		accountNonce := uint64(i)
		tx := types.NewTxWithoutSig(types.CONFIGCHANNELID, payloadBytes, accountNonce)
		txs = append(txs, tx)
	}

	return types.NewBlock(types.CONFIGCHANNELID, 0, types.GenesisBlockPrevHash, txs), nil
}

func CreateAdmins() ([]*types.Member, error) {
	// get pubkey from string by base64 encoding
	data, err := base64.StdEncoding.DecodeString("BGXcjZ3bhemsoLP4HgBwnQ5gsc8VM91b3y8bW0b6knkWu8x" +
		"CSKO2qiJXARMHcbtZtvU7Jos2A5kFCD1haJ/hLdg=")
	if err != nil {
		return nil,err
	}
	// create PublicKey
	pk, err := crypto.NewPublicKey(data)
	if err != nil {
		return nil, err
	}
	// create member and append it to the admin
	member, err := types.NewMember(pk, "SystemAdmin")
	admins := make([]*types.Member, 0)
	admins = append(admins, member)
	return admins,nil
}
