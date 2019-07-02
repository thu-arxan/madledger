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
	// ID is the hash of the tx while presented in hex
	ID                string
	IsValidatorUpdate bool
	Data              TxData
	Time              int64
}

// TxData is the data of Tx
type TxData struct {
	ChannelID    string
	AccountNonce uint64
	Recipient    []byte
	Payload      []byte
	Version      int32
	Sig          *TxSig
}

// TxSig is the sig of tx
type TxSig struct {
	PK  []byte
	Sig []byte
}

// NewTx is the constructor of Tx
func NewTx(channelID string, recipient common.Address, payload []byte, privKey crypto.PrivateKey, isValidatorUpdate bool) (*Tx, error) {
	if payload == nil || len(payload) == 0 {
		return nil, errors.New("The payload can not be empty")
	}
	var tx = &Tx{
		Data: TxData{
			ChannelID:    channelID,
			AccountNonce: util.RandUint64(),
			Recipient:    recipient.Bytes(),
			Payload:      payload,
			Version:      1,
			Sig:          nil,
		},
		IsValidatorUpdate: isValidatorUpdate,
		Time:              util.Now(),
	}
	hash := tx.HashWithoutSig()
	sig, err := privKey.Sign(hash)
	if err != nil {
		return nil, err
	}
	pkBytes, err := privKey.PubKey().Bytes()
	if err != nil {
		return nil, err
	}
	sigBytes, err := sig.Bytes()
	if err != nil {
		return nil, err
	}
	tx.Data.Sig = &TxSig{
		PK:  pkBytes,
		Sig: sigBytes,
	}
	tx.ID = util.Hex(tx.Hash())
	return tx, nil
}

// NewTxWithoutSig is a special kind of tx without sig, it
// is prepared for the genesis and global hash
func NewTxWithoutSig(channelID string, payload []byte, accountNonce uint64) *Tx {
	var tx = &Tx{
		Data: TxData{
			ChannelID:    channelID,
			AccountNonce: accountNonce,
			Recipient:    common.ZeroAddress.Bytes(),
			Payload:      payload,
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	}
	tx.ID = util.Hex(tx.Hash())
	return tx
}

// Verify return true if a tx is packed well, else return false
func (tx *Tx) Verify() bool {
	if util.Hex(tx.Hash()) != tx.ID {
		return false
	}
	if tx.Data.Sig == nil {
		return false
	}
	hash := tx.HashWithoutSig()
	pk, err := crypto.NewPublicKey(tx.Data.Sig.PK)
	if err != nil {
		return false
	}
	sig, err := crypto.NewSignature(tx.Data.Sig.Sig)
	if err != nil {
		return false
	}
	if !sig.Verify(hash, pk) {
		return false
	}

	return true
}

// GetSender return the sender of the tx
func (tx *Tx) GetSender() (common.Address, error) {
	// return tx.Data.Sig.PK.Address()
	if tx.Data.Sig == nil {
		return common.ZeroAddress, nil
	}
	pk, err := crypto.NewPublicKey(tx.Data.Sig.PK)
	if err != nil {
		return common.ZeroAddress, err
	}
	return pk.Address()
}

// GetReceiver return the receiver
func (tx *Tx) GetReceiver() common.Address {
	return common.BytesToAddress(tx.Data.Recipient)
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

	bytes, _ := json.Marshal(data)
	return crypto.Hash(bytes)
}

// Bytes return the bytes of tx, which is the wrapper of json.Marshal
func (tx *Tx) Bytes() ([]byte, error) {
	return json.Marshal(tx)
}

// BytesToTx convert bytes to tx, which is the wrapper of json.Unmarshal
func BytesToTx(data []byte) (*Tx, error) {
	var tx *Tx
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
