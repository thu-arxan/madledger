package config

import (
	"encoding/json"
	"madledger/core"
)

// CreateGenesisBlock return the genesis block
// maybe the address should be a special addr rather than all zero
// also the data is still need to be discussed
// TODO
func CreateGenesisBlock() (*core.Block, error) {
	var payload = Payload{
		ChannelID: core.GLOBALCHANNELID,
		Profile: Profile{
			Open: true,
		},
		Version: 1,
	}
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	// all zero
	var addr core.Address
	tx, err := core.NewTx(addr, payloadBytes)
	if err != nil {
		return nil, err
	}

	return core.NewBlock(0, nil, []*core.Tx{tx}), nil
}
