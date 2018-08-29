package evm

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"madledger/common"
	"madledger/common/util"
	"madledger/executor/evm/abi"
	"madledger/executor/evm/simulate"
	"os"
	"testing"
)

/*
* This file will load some sol files and run them on the evm.
 */

var (
	gopath        = os.Getenv("GOPATH")
	deployAccount = common.NewDefaultAccount(common.ZeroAddress)
)

func TestBalance(t *testing.T) {
	contractCodes, err := readCodes(getFilePath("Balance.bin"))
	if err != nil {
		t.Fatal(err)
	}
	db := simulate.NewStateDB()
	user := newAccount(1)
	db.SetAccount(user)
	vm := NewEVM(newContext(), user.Address(), db)
	code, contractAddr, err := vm.Create(user, contractCodes, []byte{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	runtimeCodes, err := readCodes(getFilePath("Balance.bin-runtime"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(code, runtimeCodes) {
		t.Fatal()
	}
	contractAccount, err := vm.cache.GetAccount(contractAddr)
	if err != nil {
		t.Fatal(err)
	}
	// add 34
	input, _ := hex.DecodeString("1003e2d20000000000000000000000000000000000000000000000000000000000000022")
	output, err := vm.Call(user, contractAccount, code, input, 0)
	values, err := abi.Unpacker(getFilePath("Balance.abi"), "add", output)
	if err != nil {
		t.Fatal(err)
	}
	if values[0].Value != "44" {
		t.Fatal(values[0].Value)
	}
	// sub 5
	input, _ = hex.DecodeString("27ee58a60000000000000000000000000000000000000000000000000000000000000005")
	output, err = vm.Call(user, contractAccount, code, input, 0)
	values, err = abi.Unpacker(getFilePath("Balance.abi"), "sub", output)
	if err != nil {
		t.Fatal(err)
	}
	if values[0].Value != "39" {
		t.Fatal(values[0].Value)
	}
	// set 1314
	input, _ = hex.DecodeString("60fe47b10000000000000000000000000000000000000000000000000000000000000522")
	output, err = vm.Call(user, contractAccount, code, input, 0)
	values, err = abi.Unpacker(getFilePath("Balance.abi"), "set", output)
	if err != nil {
		t.Fatal(err)
	}
	if values[0].Value != "true" {
		t.Fatal(values[0].Value)
	}
	// get
	input, _ = hex.DecodeString("6d4ce63c")
	output, err = vm.Call(user, contractAccount, code, input, 0)
	values, err = abi.Unpacker(getFilePath("Balance.abi"), "get", output)
	if err != nil {
		t.Fatal(err)
	}
	if values[0].Value != "1314" {
		t.Fatal(values[0].Value)
	}
	// info
	input, _ = hex.DecodeString("370158ea")
	output, err = vm.Call(user, contractAccount, code, input, 0)
	values, err = abi.Unpacker(getFilePath("Balance.abi"), "info", output)
	if err != nil {
		t.Fatal(err)
	}
	if values[0].Value != user.Address().String() {
		t.Fatal(values[0].Value)
	}
	if values[1].Value != "1314" {
		t.Fatal(values[1].Value)
	}
	// Then will check if the code or anything else store on the statedb
	contractAccount, err = db.GetAccount(contractAccount.GetAddress())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(contractAccount.GetCode(), runtimeCodes) {
		t.Fatal(errors.New("The code of user is not same with the runtime code"))
	}
	if contractAccount.GetNonce() != 0 {
		t.Fatal(fmt.Errorf("The nonce of contract is %d", contractAccount.GetNonce()))
	}
	user, err = db.GetAccount(user.GetAddress())
	if user.GetNonce() != 1 {
		t.Fatal(fmt.Errorf("The nonce of user is %d", user.GetNonce()))
	}
}

func getFilePath(name string) string {
	path, _ := util.MakeFileAbs("src/madledger/executor/evm/sols/output/"+name, gopath)
	return path
}

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}
