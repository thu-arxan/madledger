// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/crypto/hash"
	"madledger/common/util"
)

// Tx is the transaction, which structure is not decided yet
// Note: The Time is not important and will cause some consensus problems, so it won't
// be included while cacluating the hash
type Tx struct {
	// ID is the hash of the tx while presented in hex
	ID   string `json:"id,omitempty"`
	Data TxData `json:"data,omitempty"`
	Time int64  `json:"time,omitempty"`
	// Below are some caches
	sender *common.Address
}

// TxType is the type of consensus
type TxType int64

// Here define some kind of tx type
const (
	_ TxType = iota
	CREATECHANNEL
	// VALIDATOR is the tendermint cfgChange tx
	VALIDATOR
	// NODE is the raft cfgChange tx
	NODE
)

// TxData is the data of Tx
type TxData struct {
	ChannelID string `json:"channelID,omitempty"`
	Nonce     uint64 `json:"nonce,omitempty"`
	Recipient []byte `json:"recipient,omitempty"`
	Payload   []byte `json:"payload,omitempty"`
	Value     uint64 `json:"value,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Version   int32  `json:"version,omitempty"`
	Sig       TxSig  `json:"sig,omitempty"`
}

// TxSig is the sig of tx
type TxSig struct {
	PK   []byte           `json:"pk,omitempty"`
	Sig  []byte           `json:"sig,omitempty"`
	Algo crypto.Algorithm `json:"algo,omitempty"`
}

// NewTx is the constructor of Tx
func NewTx(channelID string, recipient common.Address, payload []byte, value uint64, msg string, privKey crypto.PrivateKey) (*Tx, error) {
	if payload == nil || len(payload) == 0 {
		return nil, errors.New("The payload can not be empty")
	}
	switch privKey.Algo() {
	case crypto.KeyAlgoSecp256k1, crypto.KeyAlgoSM2:
		// do nothing
	default:
		return nil, fmt.Errorf("unsupported key algo:%v", privKey.Algo())
	}
	var tx = &Tx{
		Data: TxData{
			ChannelID: channelID,
			Nonce:     util.RandUint64(),
			Recipient: recipient.Bytes(),
			Payload:   payload,
			Value:     value,
			Msg:       msg,
			Version:   1,
		},
		Time: util.Now(),
	}
	hash := tx.hashWithoutSig(privKey.Algo())
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
	tx.Data.Sig = TxSig{
		PK:   pkBytes,
		Sig:  sigBytes,
		Algo: privKey.Algo(),
	}
	tx.ID = util.Hex(tx.Hash(privKey.Algo()))
	return tx, nil
}

// NewTxWithoutSig is a special kind of tx without sig,
// it is prepared for the genesis and global hash
func NewTxWithoutSig(channelID string, payload []byte, nonce uint64) *Tx {
	var tx = &Tx{
		Data: TxData{
			ChannelID: channelID,
			Nonce:     nonce,
			Recipient: common.ZeroAddress.Bytes(),
			Payload:   payload,
			Version:   1,
		},
		Time: util.Now(),
	}
	tx.ID = util.Hex(tx.Hash(crypto.KeyAlgoSM2))
	return tx
}

// Verify return true if a tx is packed well, else return false
func (tx *Tx) Verify() bool {
	var algo = tx.Data.Sig.Algo
	switch algo {
	case crypto.KeyAlgoSM2, crypto.KeyAlgoSecp256k1:
		// do nothing
	default:
		return false
	}
	if util.Hex(tx.Hash(algo)) != tx.ID {
		return false
	}
	hash := tx.hashWithoutSig(algo)
	pk, err := crypto.NewPublicKey(tx.Data.Sig.PK, algo)
	if err != nil {
		return false
	}
	sig, err := crypto.NewSignature(tx.Data.Sig.Sig, algo)
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
	if tx.sender != nil {
		return *(tx.sender), nil
	}
	var algo = tx.Data.Sig.Algo
	switch algo {
	case crypto.KeyAlgoSecp256k1, crypto.KeyAlgoSM2:
		// do nothing
	default:
		return common.ZeroAddress, fmt.Errorf("unsupport algo:%v", algo)
	}
	pk, err := crypto.NewPublicKey(tx.Data.Sig.PK, algo)
	if err != nil {
		return common.ZeroAddress, err
	}
	sender, err := pk.Address()
	if err == nil {
		tx.sender = &sender
	}
	return sender, err
}

// GetReceiver return the receiver
func (tx *Tx) GetReceiver() common.Address {
	return common.BytesToAddress(tx.Data.Recipient)
}

// Hash return the hash of tx
// Note: Be careful to make sure tx.Data.Sig.Algo right
func (tx *Tx) Hash(algo ...crypto.Algorithm) []byte {
	if len(algo) != 0 {
		return tx.hash(true, algo[0])
	}
	return tx.hash(true, tx.Data.Sig.Algo)
}

// hashWithoutSig return the hash of tx without sig
func (tx *Tx) hashWithoutSig(algo crypto.Algorithm) []byte {
	return tx.hash(false, algo)
}

// hash implementation different hash
// Note: is algo is not secp256k1, regard it as sm3
func (tx *Tx) hash(withSig bool, algo crypto.Algorithm) []byte {
	var sig = tx.Data.Sig
	if !withSig {
		tx.Data.Sig = TxSig{}
	}

	bytes, _ := json.Marshal(tx.Data)
	tx.Data.Sig = sig

	switch algo {
	case crypto.KeyAlgoSecp256k1:
		return hash.SHA256(bytes)
	default:
		return hash.SM3(bytes)
	}
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
