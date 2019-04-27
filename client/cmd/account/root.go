package account

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	homeDir, _ = os.Getwd()
	accountCmd = &cobra.Command{
		Use: "account",
	}
)

// Cmd return the account command
func Cmd() *cobra.Command {
	accountCmd.AddCommand(listCmd)
	return accountCmd
}
