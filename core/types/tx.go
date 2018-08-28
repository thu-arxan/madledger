package types

import (
	"encoding/json"
	"madledger/common"
	"madledger/util"
	"math/big"
)

// Tx is the transaction, which structure is not decided yet
type Tx struct {
	Data TxData
	Time int64
}

// TxData is the data of Tx
// todo: maybe should contain ChannelID
type TxData struct {
	AccountNonce uint64
	Recipient    common.Address
	Payload      []byte
	Version      int32
	Sig          *TxSig
}

// TxSig is the sig of tx
type TxSig struct {
	// maybe change to bytes is better
	V    *big.Int
	R, S *big.Int
}

// NewTx is the constructor of Tx
// TODO: AccountNonce & sig
func NewTx(recipient common.Address, payload []byte) (*Tx, error) {
	var tx = Tx{
		Data: TxData{
			AccountNonce: 0,
			Recipient:    recipient,
			Payload:      payload,
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	}
	return &tx, nil
}

// Hash return the hash of tx
func (tx *Tx) Hash() []byte {
	return tx.hash(true)
}

// HashWithoutSig return the hash of tx without sig
func (tx *Tx) HashWithoutSig() []byte {
	return tx.hash(false)
}

// hash implementation different hash
func (tx *Tx) hash(withSig bool) []byte {
	var data TxData
	if withSig {
		data = tx.Data
	} else { // clone
		data.AccountNonce = tx.Data.AccountNonce
		data.Recipient = tx.Data.Recipient
		data.Payload = tx.Data.Payload
		data.Version = tx.Data.Version
		data.Sig = nil
	}

	// todo: should we ignore the error?
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return util.Hash(bytes)
}
