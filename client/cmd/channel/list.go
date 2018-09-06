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
	listCmd.Flags().StringP("system", "s", "true", "If the system channel is contained")
	listViper.BindPFlag("system", listCmd.Flags().Lookup("system"))
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := orderer.NewClient()
	if err != nil {
		return err
	}

	var system = true
	if listViper.GetString("system") == "false" {
		system = false
	}

	return client.ListChannel(system)
}
