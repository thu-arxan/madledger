package node

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	homeDir, _ = os.Getwd()
	nodeCmd    = &cobra.Command{
		Use: "node",
	}
)

// Cmd return the channel command
func Cmd() *cobra.Command {
	nodeCmd.AddCommand(addCmd)
	nodeCmd.AddCommand(removeCmd)
	return nodeCmd
}
