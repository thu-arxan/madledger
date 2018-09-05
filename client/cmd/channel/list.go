package channel

import (
	"madledger/client/orderer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	listCmd = &cobra.Command{
		Use: "list",
	}
	listViper = viper.New()
)

func init() {
	listCmd.RunE = runList
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := orderer.NewClient()
	if err != nil {
		return err
	}

	return client.ListChannel()
}
