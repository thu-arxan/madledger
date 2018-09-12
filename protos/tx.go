package protos

import (
	"encoding/json"
	"madledger/core/types"
)

// NewTx is the constructor of tx
// func NewTx(channelID string, recipient, payload []byte, privKey crypto.PrivateKey) (*Tx, error) {
// 	txData := TxData{
// 		ChannelID:    channelID,
// 		AccountNonce: 0,
// 		Recipient:    recipient,
// 		Payload:      payload,
// 		Version:      1,
// 	}
// 	txDataBytes, err := json.Marshal(txData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	sig, err := privKey.Sign(crypto.Hash(txDataBytes))
// 	if err != nil {
// 		return nil, err
// 	}
// 	pk, err := privKey.PubKey().Bytes()
// 	if err != nil {
// 		return nil, err
// 	}
// 	sigBytes, err := sig.Bytes()
// 	if err != nil {
// 		return nil, err
// 	}
// 	txData.Sig = &TxSig{
// 		PK:  pk,
// 		Sig: sigBytes,
// 	}
// 	return &Tx{
// 		Data: &txData,
// 		Time: util.Now(),
// 	}, nil
// }

// ConvertToTypes convert pb.Tx to types.Tx
func (tx *Tx) ConvertToTypes() (*types.Tx, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	var typesTx types.Tx
	err = json.Unmarshal(data, &typesTx)
	if err != nil {
		return nil, err
	}
	return &typesTx, nil
}

// NewTx is the constructor of Tx
func NewTx(typesTx *types.Tx) (*Tx, error) {
	data, err := json.Marshal(typesTx)
	if err != nil {
		return nil, err
	}
	var tx Tx
	err = json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// // ConvertTxFromPbToTypes convert Tx from pb.Tx to types.Tx
// func ConvertTxFromPbToTypes(pbTx *pb.Tx) (*types.Tx, error) {
// 	data, err := json.Marshal(pbTx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var typesTx types.Tx
// 	err = json.Unmarshal(data, &typesTx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &typesTx, nil
// }

// // ConvertTxFromTypesToPb convert Tx from types.Tx to pb.Tx
// func ConvertTxFromTypesToPb(typesTx *types.Tx) (*pb.Tx, error) {
// 	data, err := json.Marshal(typesTx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var pbTx pb.Tx
// 	err = json.Unmarshal(data, &pbTx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pbTx, nil
// }
