// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tests

import (
	"encoding/hex"
	"errors"
	"madledger/common/abi"
	"io/ioutil"
	pb "madledger/protos"
	"os"
)

var (
	gopath = os.Getenv("GOPATH")
)

// Help fulfillment the test.

func initDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	return nil
}

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}

type txStatus struct {
	BlockNumber uint64
	BlockIndex  int32
	Output      []string
}

func getTxStatus(abiPath, funcName string, status *pb.TxStatus) (*txStatus, error) {
	if status.Err != "" {
		return nil, errors.New(status.Err)
	}
	values, err := abi.Unpack(abiPath, funcName, status.Output)
	if err != nil {
		return nil, err
	}
	var txStatus = &txStatus{
		BlockNumber: status.BlockNumber,
		BlockIndex:  status.BlockIndex,
	}

	for _, value := range values {
		txStatus.Output = append(txStatus.Output, value)
	}
	return txStatus, nil
}
