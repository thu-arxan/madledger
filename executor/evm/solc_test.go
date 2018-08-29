package evm

import (
	"bytes"
	"encoding/hex"
	"errors"
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

func TestHelloWorld(t *testing.T) {
	contractCodes, err := readCodes(getFilePath("HelloWorld.bin"))
	if err != nil {
		t.Fatal(err)
	}
	db := simulate.NewStateDB()
	user := newAccount(1)
	db.SetAccount(user)
	vm := NewEVM(newContext(), user.Address(), db)
	// before create the the contrat, the nonce should be add one, but this is not a good thing
	// because this should be done automic
	user.SetNonce(user.GetNonce() + 1)
	db.SetAccount(user)
	output, contractAddr, err := vm.Create(user, contractCodes, []byte{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	runtimeCodes, err := readCodes(getFilePath("HelloWorld.bin-runtime"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(output, runtimeCodes) {
		t.Fatal()
	}
	contractAccount, err := vm.cache.GetAccount(contractAddr)
	if err != nil {
		t.Fatal(err)
	}
	input, _ := hex.DecodeString("82ab890a0000000000000000000000000000000000000000000000000000000000000045")

	output, err = vm.Call(user, contractAccount, output, input, 0)
	values, err := abi.Unpacker(getFilePath("HelloWorld.abi"), "update", output)
	if err != nil {
		t.Fatal(err)
	}
	for i, value := range values {
		switch i {
		case 0:
			if value.Value != user.Address().String() {
				t.Fatal()
			}
		case 1:
			if value.Value != "69" {
				t.Fatal()
			}
		default:
			t.Fatal()
		}
	}
	// Then will check if the code or anything else store on the statedb
	contractAccount, err = db.GetAccount(contractAccount.GetAddress())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(contractAccount.GetCode(), runtimeCodes) {
		t.Fatal(errors.New("The code of user is not same with the runtime code"))
	}
	if contractAccount.GetNonce() != 1 {
		t.Fatal(errors.New("The nonce of contract is not 1"))
	}
	user, err = db.GetAccount(user.GetAddress())
	if user.GetNonce() != 2 {
		t.Fatal(errors.New("The nonce of user is not 2"))
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
