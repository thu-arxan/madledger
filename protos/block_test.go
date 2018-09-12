package protos

import (
	"encoding/hex"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core/types"
	"reflect"
	"testing"
)

func TestConvertBlock(t *testing.T) {
	// test block without tx
	block := types.NewBlock("test", 1, types.GenesisBlockPrevHash, nil)
	typesBlock, err := convertTypesBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(typesBlock, block) {
		t.Fatal()
	}
	// test block with tx but without sig
	var txs []*types.Tx
	txs = append(txs, &types.Tx{
		Data: types.TxData{
			ChannelID:    types.GLOBALCHANNELID,
			AccountNonce: 0,
			Recipient:    common.ZeroAddress.Bytes(),
			Payload:      nil,
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	})
	txs = append(txs, &types.Tx{
		Data: types.TxData{
			ChannelID:    types.GLOBALCHANNELID,
			AccountNonce: 0,
			Recipient:    common.ZeroAddress.Bytes(),
			Payload:      []byte("Hello World"),
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	})
	block = types.NewBlock(types.GLOBALCHANNELID, 1, types.GenesisBlockPrevHash, txs)
	typesBlock, err = convertTypesBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(typesBlock, block) {
		t.Fatal()
	}
	// test block with tx and with sig
	rawPrivKey, _ := hex.DecodeString("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032")
	privKey, _ := crypto.NewPrivateKey(rawPrivKey)
	sigTx, err := types.NewTx(types.GLOBALCHANNELID, common.ZeroAddress, []byte("Hello World again"), privKey)
	if err != nil {
		t.Fatal(err)
	}
	txs = append(txs, sigTx)
	block = types.NewBlock(types.GLOBALCHANNELID, 1, types.GenesisBlockPrevHash, txs)
	typesBlock, err = convertTypesBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(typesBlock, block) {
		t.Fatal()
	}
}

func convertTypesBlock(block *types.Block) (*types.Block, error) {
	pbBlock, err := NewBlock(block)
	if err != nil {
		return nil, err
	}
	typesBlock, err := pbBlock.ConvertToTypes()
	if err != nil {
		return nil, err
	}
	return typesBlock, nil
}
