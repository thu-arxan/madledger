package validator

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	homeDir, _ = os.Getwd()
	validatorCmd = &cobra.Command{
		Use: "validator",
	}
)

// Cmd return the channel command
func Cmd() *cobra.Command {
	validatorCmd.AddCommand(addCmd)
	//channelCmd.AddCommand(listCmd)
	return validatorCmd
}

