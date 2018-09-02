package types

import (
	"encoding/json"
	"fmt"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
)

// Tx is the transaction, which structure is not decided yet
// Note: The Time is not important and will cause some consensus problems, so it won't
// be included while cacluating the hash
type Tx struct {
	Data *TxData
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
	var tx = Tx{
		Data: &TxData{
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
	return &tx, nil
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
	pk, err := tx.Data.Sig.PK.Bytes()
	if err != nil {
		return common.ZeroAddress, err
	}
	fmt.Println(util.Hex(pk))
	return common.ZeroAddress, nil
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
	var data *TxData
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
