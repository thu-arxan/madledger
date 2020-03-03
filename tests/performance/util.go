package performance

import (
	"encoding/hex"
	"errors"
	"evm/abi"
	"fmt"
	"io/ioutil"
	client "madledger/client/lib"
	pb "madledger/protos"
	"madledger/tests/performance/bft"
	"madledger/tests/performance/raft"
	"madledger/tests/performance/solo"
	"os"
)

var (
	gopath  = os.Getenv("GOPATH")
	logPath = "performance.out"
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
		fmt.Println("here>>>", status.Output)
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

func writeLog(log string) error {
	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = file.WriteString(log)
	return err
}

func getClients() []*client.Client {
	var clients []*client.Client
	switch consensus {
	case "solo":
		clients = solo.GetClients()
	case "raft":
		clients = raft.GetClients()
	case "bft":
		clients = bft.GetClients()
	default:
		panic("Unsupport consensus")
	}
	return clients
}
