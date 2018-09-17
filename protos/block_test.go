package protos

import (
	"encoding/hex"
	"madledger/common"
	"madledger/common/crypto"
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
	txs = append(txs, types.NewTxWithoutSig(types.GLOBALCHANNELID, nil, 0))
	txs = append(txs, types.NewTxWithoutSig(types.GLOBALCHANNELID, []byte("Hello World"), 0))
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
