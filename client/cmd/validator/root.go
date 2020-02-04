package validator

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	homeDir, _   = os.Getwd()
	validatorCmd = &cobra.Command{
		Use: "validator",
	}
)

// Cmd return the channel command
func Cmd() *cobra.Command {
	validatorCmd.AddCommand(addCmd)
	return validatorCmd
}
