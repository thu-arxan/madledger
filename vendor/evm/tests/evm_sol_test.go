package tests

import (
	"evm"
	"evm/db"
	"evm/example"
	"evm/util"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	evmCodeBin = "sols/Ethereum_sol_OpCodes.bin"
	evmCodeAbi = "sols/Ethereum_sol_OpCodes.abi"
	evmCode []byte
	evmCodeAddress evm.Address
)

func TestEvm(t *testing.T) {
	binBytes, err := util.ReadBinFile(evmCodeBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	evmCode, evmCodeAddress = deployContract(t, memoryDB, bc, origin, binBytes, "", "", 388049)
	input := mustPack(evmCodeAbi, "test")
	var gas uint64 = 1000000
	output, err := evm.New(bc, memoryDB, &evm.Context{
		Input: input,
		Value: 0,
		Gas: &gas,
		BlockHeight: 1,
	}).Call(origin, evmCodeAddress, evmCode)
	require.NoError(t, err)
	t.Log(output)
	t.Log(gas)
}