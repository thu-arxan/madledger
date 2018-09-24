package tests

import (
	"encoding/hex"
	"io/ioutil"
	"madledger/common/abi"
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
	values, err := abi.Unpacker(abiPath, funcName, status.Output)
	if err != nil {
		return nil, err
	}
	var txStatus = &txStatus{
		BlockNumber: status.BlockNumber,
		BlockIndex:  status.BlockIndex,
	}

	for _, value := range values {
		txStatus.Output = append(txStatus.Output, value.Value)
	}
	return txStatus, nil
}