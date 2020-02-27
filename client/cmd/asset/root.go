package asset

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	homeDir, _ = os.Getwd()
	assetCmd   = &cobra.Command{
		Use: "asset",
	}
)

// Cmd return the account command
func Cmd() *cobra.Command {
	assetCmd.AddCommand(issueCmd)
	assetCmd.AddCommand(transferCmd)
	// todo:@zhq, i have no instruction to know how many assets that i have.
	return assetCmd
}
