package protos

import (
	"encoding/json"
	"madledger/common/crypto"
	"madledger/common/util"
)

// NewTx is the constructor of tx
func NewTx(channelID string, recipient, payload []byte, privKey crypto.PrivateKey) (*Tx, error) {
	txData := TxData{
		ChannelID:    channelID,
		AccountNonce: 0,
		Recipient:    recipient,
		Payload:      payload,
		Version:      1,
	}
	txDataBytes, err := json.Marshal(txData)
	if err != nil {
		return nil, err
	}
	sig, err := privKey.Sign(crypto.Hash(txDataBytes))
	if err != nil {
		return nil, err
	}
	pk, err := privKey.PubKey().Bytes()
	if err != nil {
		return nil, err
	}
	sigBytes, err := sig.Bytes()
	if err != nil {
		return nil, err
	}
	txData.Sig = &TxSig{
		PK:  pk,
		Sig: sigBytes,
	}
	return &Tx{
		Data: &txData,
		Time: util.Now(),
	}, nil
}
