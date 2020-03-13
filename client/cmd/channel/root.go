package channel

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	homeDir, _ = os.Getwd()
	channelCmd = &cobra.Command{
		Use: "channel",
	}
)

// Cmd return the channel command
func Cmd() *cobra.Command {
	channelCmd.AddCommand(createCmd)
	channelCmd.AddCommand(listCmd)
	channelCmd.AddCommand(tokenCmd)
	return channelCmd
}
