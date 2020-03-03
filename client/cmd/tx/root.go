package tx

import (
	"encoding/hex"
	"evm/abi"
	"fmt"
	"io/ioutil"
	"madledger/common/util"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	abi.SetAddressParser(20, func(bs []byte) string {
		return "0x" + fmt.Sprintf("%x", bs)
	}, func(addr string) ([]byte, error) {
		return util.HexToBytes(addr)
	})
}

var (
	homeDir, _ = os.Getwd()
	txCmd      = &cobra.Command{
		Use: "tx",
	}
)

// Cmd return the tx command
func Cmd() *cobra.Command {
	txCmd.AddCommand(createCmd)
	txCmd.AddCommand(callCmd)
	txCmd.AddCommand(historyCmd)
	return txCmd
}

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}
