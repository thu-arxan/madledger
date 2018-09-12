package tx

import (
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

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
	return txCmd
}

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}
