package evm

import (
	"encoding/hex"
	"io/ioutil"
	"madledger/common/abi"
	"madledger/common/util"
	"madledger/executor/evm/simulate"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
* This file will load some sol files and run them on the evm.
 */

var (
	simulateDB = simulate.NewStateDB()
	wb         = simulateDB.NewWriteBatch()
	user       = newAccount(1)
)

func TestBalance(t *testing.T) {
	contractCodes, err := readCodes(getFilePath("Balance.bin"))
	require.NoError(t, err)

	simulateDB.SetAccount(user)
	vm := NewEVM(newContext(), user.GetAddress(), simulateDB, wb)
	code, contractAddr, err := vm.Create(user, contractCodes, []byte{}, 0)
	require.NoError(t, err)

	runtimeCodes, err := readCodes(getFilePath("Balance.bin-runtime"))
	require.NoError(t, err)
	require.Equal(t, runtimeCodes, code)

	contractAccount, err := vm.cache.GetAccount(contractAddr)
	require.NoError(t, err)

	abiFilePath := getFilePath("Balance.abi")

	// add 34
	input, _ := abi.GetPayloadBytes(abiFilePath, "add", []string{"34"})
	output, err := vm.Call(user, contractAccount, code, input, 0)
	values, err := abi.Unpacker(getFilePath("Balance.abi"), "add", output)
	require.NoError(t, err)
	require.Equal(t, "44", values[0].Value)

	// sub 5
	input, _ = abi.GetPayloadBytes(abiFilePath, "sub", []string{"5"})
	output, err = vm.Call(user, contractAccount, code, input, 0)
	values, err = abi.Unpacker(getFilePath("Balance.abi"), "sub", output)
	require.NoError(t, err)
	require.Equal(t, "39", values[0].Value)

	// set 1314
	input, _ = abi.GetPayloadBytes(abiFilePath, "set", []string{"1314"})
	output, err = vm.Call(user, contractAccount, code, input, 0)
	values, err = abi.Unpacker(getFilePath("Balance.abi"), "set", output)
	require.NoError(t, err)
	require.Equal(t, "true", values[0].Value)

	// get
	input, _ = abi.GetPayloadBytes(abiFilePath, "get", nil)
	output, err = vm.Call(user, contractAccount, code, input, 0)
	values, err = abi.Unpacker(getFilePath("Balance.abi"), "get", output)
	require.NoError(t, err)
	require.Equal(t, "1314", values[0].Value)

	// info
	input, _ = abi.GetPayloadBytes(abiFilePath, "info", nil)
	output, err = vm.Call(user, contractAccount, code, input, 0)
	values, err = abi.Unpacker(getFilePath("Balance.abi"), "info", output)
	require.NoError(t, err)
	require.Equal(t, user.GetAddress().String(), values[0].Value)
	require.Equal(t, "1314", values[1].Value)
	// Then will check if the code or anything else store on the statedb
	contractAccount, err = simulateDB.GetAccount(contractAccount.GetAddress())
	require.NoError(t, err)
	require.Equal(t, runtimeCodes, contractAccount.GetCode())

	// user, err = simulateDB.GetAccount(user.GetAddress())
}

func TestDuplicateAddress(t *testing.T) {
	contractCodes, _ := readCodes(getFilePath("Balance.bin"))
	vm := NewEVM(newContext(), user.GetAddress(), simulateDB, wb)
	_, _, err := vm.Create(user, contractCodes, []byte{}, 0)
	require.EqualError(t, err, "Duplicate address")
}

func getFilePath(name string) string {
	path, _ := util.MakeFileAbs("src/madledger/executor/evm/sols/output/"+name, os.Getenv("GOPATH"))
	return path
}

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}
