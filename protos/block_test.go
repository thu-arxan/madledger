package protos

import (
	"encoding/hex"
	fmt "fmt"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertBlock(t *testing.T) {
	// test block without tx
	block := core.NewBlock("test", 1, core.GenesisBlockPrevHash, nil)
	typesBlock, err := convertTypesBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	require.EqualValues(t, block, typesBlock)
	// test block with tx but without sig
	var txs []*core.Tx
	txs = append(txs, core.NewTxWithoutSig(core.GLOBALCHANNELID, nil, 0))
	txs = append(txs, core.NewTxWithoutSig(core.GLOBALCHANNELID, []byte("Hello World"), 0))
	block = core.NewBlock(core.GLOBALCHANNELID, 1, core.GenesisBlockPrevHash, txs)
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
	sigTx, err := core.NewTx(core.GLOBALCHANNELID, common.ZeroAddress, []byte("Hello World again"), 0, "", privKey)
	if err != nil {
		t.Fatal(err)
	}
	txs = append(txs, sigTx)
	block = core.NewBlock(core.GLOBALCHANNELID, 1, core.GenesisBlockPrevHash, txs)
	typesBlock, err = convertTypesBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(typesBlock, block) {
		fmt.Printf("%s\n", string(typesBlock.Bytes()))
		fmt.Println("")
		fmt.Printf("%s\n", string(block.Bytes()))
		t.Fatal()
	}
}

func convertTypesBlock(block *core.Block) (*core.Block, error) {
	pbBlock, err := NewBlock(block)
	if err != nil {
		return nil, err
	}
	typesBlock, err := pbBlock.ToCore()
	if err != nil {
		return nil, err
	}
	return typesBlock, nil
}
