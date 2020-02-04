package protos

import (
	"madledger/common/util"
	"madledger/core/types"
)

// ConvertToTypes convert pb.Tx to types.Tx
func (tx *Tx) ConvertToTypes() (*types.Tx, error) {
	return &types.Tx{
		ID: tx.ID,
		Data: types.TxData{
			ChannelID:    tx.Data.ChannelID,
			AccountNonce: tx.Data.AccountNonce,
			Recipient:    util.CopyBytes(tx.Data.Recipient),
			Payload:      util.CopyBytes(tx.Data.Payload),
			Version:      tx.Data.Version,
			Sig: &types.TxSig{
				PK:  util.CopyBytes(tx.Data.Sig.PK),
				Sig: util.CopyBytes(tx.Data.Sig.Sig),
			},
		},
		Time: tx.Time,
	}, nil
}

// NewTx is the constructor of Tx
func NewTx(tx *types.Tx) (*Tx, error) {
	return &Tx{
		ID: tx.ID,
		Data: &TxData{
			ChannelID:    tx.Data.ChannelID,
			AccountNonce: tx.Data.AccountNonce,
			Recipient:    util.CopyBytes(tx.Data.Recipient),
			Payload:      util.CopyBytes(tx.Data.Payload),
			Version:      tx.Data.Version,
			Sig: &TxSig{
				PK:  util.CopyBytes(tx.Data.Sig.PK),
				Sig: util.CopyBytes(tx.Data.Sig.Sig),
			},
		},
		Time: tx.Time,
	}, nil
}
