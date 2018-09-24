package channel

import (
	"encoding/json"
	"errors"
	cc "madledger/blockchain/config"
	"madledger/core/types"
	"madledger/peer/db"
)

// AddConfigBlock add a config block
// todo: private channel
func (m *Manager) AddConfigBlock(block *types.Block) error {
	for i, tx := range block.Transactions {
		status := &db.TxStatus{
			Err:         "",
			BlockNumber: block.Header.Number,
			BlockIndex:  i,
			Output:      nil,
		}
		payload, err := getConfigPayload(tx)
		if err == nil {
			channelID := payload.ChannelID
			if payload.Profile.Public {
				if !m.db.BelongChannel(channelID) {
					m.db.AddChannel(channelID)
				}
			}
		} else {
			status.Err = err.Error()
		}
		m.db.SetTxStatus(tx, status)
	}
	return nil
}

func getConfigPayload(tx *types.Tx) (*cc.Payload, error) {
	if tx.Data.ChannelID != types.CONFIGCHANNELID {
		return nil, errors.New("The tx does not belog to global channel")
	}
	var payload cc.Payload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
