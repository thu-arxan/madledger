package global

import (
	"encoding/hex"
	"encoding/json"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core/types"
)

// TODO: This must read from config files
var (
	secp256k1String      = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	rawSecp256k1Bytes, _ = hex.DecodeString(secp256k1String)
	genesisPrivKey, _    = crypto.NewPrivateKey(rawSecp256k1Bytes)
)

// CreateGenesisBlock return the genesis block
// maybe the address should be a special addr rather than all zero
// also the data is still need to be discussed
// TODO:
func CreateGenesisBlock(payloads []*Payload) (*types.Block, error) {
	var txs []*types.Tx
	for _, payload := range payloads {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		// all zero
		var addr common.Address
		tx, err := types.NewTx(addr, payloadBytes, genesisPrivKey)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return types.NewBlock(0, nil, txs), nil
}
