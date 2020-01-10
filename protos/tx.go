package protos

import (
	"madledger/common/util"
	"madledger/core/types"
)

// ConvertToTypes convert pb.Tx to types.Tx
func (tx *Tx) ConvertToTypes() (*types.Tx, error) {
	t := &types.Tx{
		ID:   tx.ID,
		Time: tx.Time,
	}
	if tx.Data != nil {
		t.Data = *(tx.Data.ToTypes())
	}
	return t, nil
}

// NewTx is the constructor of Tx
func NewTx(tx *types.Tx) (*Tx, error) {
	if tx == nil {
		return nil, nil
	}
	return &Tx{
		ID:   tx.ID,
		Data: NewTxData(&(tx.Data)),
		Time: tx.Time,
	}, nil
}

// NewTxData convert types.TxData to TxData
func NewTxData(txData *types.TxData) *TxData {
	if txData == nil {
		return nil
	}
	var td = &TxData{
		ChannelID: txData.ChannelID,
		Nonce:     txData.Nonce,
		Recipient: util.CopyBytes(txData.Recipient),
		Payload:   util.CopyBytes(txData.Payload),
		Version:   txData.Version,
	}
	if txData.Sig != nil {
		td.Sig = &TxSig{
			PK:  util.CopyBytes(txData.Sig.PK),
			Sig: util.CopyBytes(txData.Sig.Sig),
		}
	}
	return td
}

// ToTypes convert TxData to types.TxData
func (data *TxData) ToTypes() *types.TxData {
	var td = &types.TxData{
		ChannelID: data.ChannelID,
		Nonce:     data.Nonce,
		Recipient: util.CopyBytes(data.Recipient),
		Payload:   util.CopyBytes(data.Payload),
		Version:   data.Version,
	}
	if data.Sig != nil {
		td.Sig = &types.TxSig{
			PK:  util.CopyBytes(data.Sig.PK),
			Sig: util.CopyBytes(data.Sig.Sig),
		}
	}
	return td
}
