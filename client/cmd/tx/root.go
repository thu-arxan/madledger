package tx

import (
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
	txCmd.AddCommand(sendCmd)
	return txCmd
}
