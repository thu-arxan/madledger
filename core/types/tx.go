package types

import (
	"encoding/json"
	"errors"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
)

// Tx is the transaction, which structure is not decided yet
// Note: The Time is not important and will cause some consensus problems, so it won't
// be included while cacluating the hash
type Tx struct {
	Data TxData
	Time int64
}

// TxData is the data of Tx
type TxData struct {
	ChannelID    string
	AccountNonce uint64
	Recipient    common.Address
	Payload      []byte
	Version      int32
	Sig          *TxSig
}

// TxSig is the sig of tx
type TxSig struct {
	PK  crypto.PublicKey
	Sig crypto.Signature
}

// NewTx is the constructor of Tx
// TODO: AccountNonce
func NewTx(channelID string, recipient common.Address, payload []byte, privKey crypto.PrivateKey) (*Tx, error) {
	if payload == nil || len(payload) == 0 {
		return nil, errors.New("The payload can not be empty")
	}
	var tx = &Tx{
		Data: TxData{
			ChannelID:    channelID,
			AccountNonce: 0,
			Recipient:    recipient,
			Payload:      payload,
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	}
	hash := tx.HashWithoutSig()
	sig, err := privKey.Sign(hash)
	if err != nil {
		return nil, err
	}
	tx.Data.Sig = &TxSig{
		PK:  privKey.PubKey(),
		Sig: sig,
	}
	return tx, nil
}

// Verify return true if a tx is packed well, else return false
func (tx *Tx) Verify() bool {
	if tx.Data.Sig == nil {
		return false
	}
	hash := tx.HashWithoutSig()
	if !tx.Data.Sig.Sig.Verify(hash, tx.Data.Sig.PK) {
		return false
	}
	return true
}

// GetSender return the sender of the tx
func (tx *Tx) GetSender() (common.Address, error) {
	return tx.Data.Sig.PK.Address()
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
		data.ChannelID = tx.Data.ChannelID
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
	return crypto.Hash(bytes)
}

type marshalTx struct {
	Data marshalTxData `json:"Data,omitempty"`
	Time int64         `json:"Time,omitempty"`
}

type marshalTxData struct {
	ChannelID    string         `json:"ChannelID,omitempty"`
	AccountNonce uint64         `json:"AccountNonce,omitempty"`
	Recipient    common.Address `json:"Recipient,omitempty"`
	Payload      []byte         `json:"Payload,omitempty"`
	Version      int32          `json:"Version,omitempty"`
	Sig          *marshalTxSig  `json:"Sig,omitempty"`
}

// TxSig is the sig of tx
type marshalTxSig struct {
	PK  []byte `json:"PK,omitempty"`
	Sig []byte `Sig:"Sig,omitempty"`
}

// Bytes return the bytes of tx
func (tx *Tx) Bytes() ([]byte, error) {
	var pk []byte
	var sig []byte
	var err error
	var mt marshalTx
	if tx.Data.Sig != nil {
		pk, err = tx.Data.Sig.PK.Bytes()
		if err != nil {
			return nil, err
		}
		sig, err = tx.Data.Sig.Sig.Bytes()
		if err != nil {
			return nil, err
		}
		txSig := marshalTxSig{
			PK:  pk,
			Sig: sig,
		}
		mt = marshalTx{
			Data: marshalTxData{
				ChannelID:    tx.Data.ChannelID,
				AccountNonce: tx.Data.AccountNonce,
				Recipient:    tx.Data.Recipient,
				Payload:      tx.Data.Payload,
				Version:      tx.Data.Version,
				Sig:          &txSig,
			},
			Time: tx.Time,
		}
	} else {
		mt = marshalTx{
			Data: marshalTxData{
				ChannelID:    tx.Data.ChannelID,
				AccountNonce: tx.Data.AccountNonce,
				Recipient:    tx.Data.Recipient,
				Payload:      tx.Data.Payload,
				Version:      tx.Data.Version,
				Sig:          nil,
			},
			Time: tx.Time,
		}
	}

	return json.Marshal(mt)
}

// BytesToTx convert bytes to tx
func BytesToTx(data []byte) (*Tx, error) {
	var mt marshalTx
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return nil, err
	}
	if mt.Data.Sig != nil {
		pk, err := crypto.NewPublicKey(mt.Data.Sig.PK)
		if err != nil {
			return nil, err
		}
		sig, err := crypto.NewSignature(mt.Data.Sig.Sig)
		if err != nil {
			return nil, err
		}
		txSig := TxSig{
			PK:  pk,
			Sig: sig,
		}
		return &Tx{
			Data: TxData{
				ChannelID:    mt.Data.ChannelID,
				AccountNonce: mt.Data.AccountNonce,
				Recipient:    mt.Data.Recipient,
				Payload:      mt.Data.Payload,
				Version:      mt.Data.Version,
				Sig:          &txSig,
			},
			Time: mt.Time,
		}, nil
	}

	return &Tx{
		Data: TxData{
			ChannelID:    mt.Data.ChannelID,
			AccountNonce: mt.Data.AccountNonce,
			Recipient:    mt.Data.Recipient,
			Payload:      mt.Data.Payload,
			Version:      mt.Data.Version,
			Sig:          nil,
		},
		Time: mt.Time,
	}, nil
}
