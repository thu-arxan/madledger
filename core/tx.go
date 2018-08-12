package core

import (
	"encoding/json"
	"madledger/util"
	"math/big"
)

// Tx is the transaction, which structure is not decided yet
// TODO
type Tx struct {
	Data txData
	Time uint64
}

// txData is the data of Tx
type txData struct {
	AccountNonce uint64
	Recipient    Address
	Payload      []byte
	Version      int
	Sig          *txSig
}

// txSig is the sig of tx
type txSig struct {
	// maybe change to bytes is better
	V    *big.Int
	R, S *big.Int
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
	var core txData
	if withSig {
		core = tx.Data
	} else { // clone
		core.AccountNonce = tx.Data.AccountNonce
		core.Recipient = tx.Data.Recipient
		core.Payload = tx.Data.Payload
		core.Version = tx.Data.Version
		core.Sig = nil
	}

	// todo: should we ignore the error?
	bytes, err := json.Marshal(core)
	if err != nil {
		return nil
	}
	return util.Hash(bytes)
}
