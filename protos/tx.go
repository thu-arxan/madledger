// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package protos

import (
	"madledger/common/util"
	"madledger/core"
)

// ToCore convert pb.Tx to core.Tx
func (tx *Tx) ToCore() (*core.Tx, error) {
	t := &core.Tx{
		ID:   tx.ID,
		Time: tx.Time,
	}
	if tx.Data != nil {
		t.Data = *(tx.Data.ToCore())
	}
	return t, nil
}

// NewTx is the constructor of Tx
func NewTx(tx *core.Tx) (*Tx, error) {
	if tx == nil {
		return nil, nil
	}
	return &Tx{
		ID:   tx.ID,
		Data: NewTxData(&(tx.Data)),
		Time: tx.Time,
	}, nil
}

// NewTxData convert core.TxData to TxData
func NewTxData(txData *core.TxData) *TxData {
	if txData == nil {
		return nil
	}
	var td = &TxData{
		ChannelID: txData.ChannelID,
		Nonce:     txData.Nonce,
		Recipient: util.CopyBytes(txData.Recipient),
		Payload:   util.CopyBytes(txData.Payload),
		Value:     txData.Value,
		Msg:       txData.Msg,
		Version:   txData.Version,
		Sig: &TxSig{
			PK:   util.CopyBytes(txData.Sig.PK),
			Sig:  util.CopyBytes(txData.Sig.Sig),
			Algo: txData.Sig.Algo,
		},
	}

	return td
}

// ToCore convert TxData to core.TxData
func (data *TxData) ToCore() *core.TxData {
	var td = &core.TxData{
		ChannelID: data.ChannelID,
		Nonce:     data.Nonce,
		Recipient: util.CopyBytes(data.Recipient),
		Payload:   util.CopyBytes(data.Payload),
		Value:     data.Value,
		Msg:       data.Msg,
		Version:   data.Version,
	}
	if data.Sig != nil {
		td.Sig = core.TxSig{
			PK:   util.CopyBytes(data.Sig.PK),
			Sig:  util.CopyBytes(data.Sig.Sig),
			Algo: data.Sig.Algo,
		}
	}
	return td
}
