package channel

import (
	"errors"
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
	listCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	listViper.BindPFlag("config", listCmd.Flags().Lookup("config"))
}

func runList(cmd *cobra.Command, args []string) error {
	cfgFile := listViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	client, err := orderer.NewClient(cfgFile)
	if err != nil {
		return err
	}

	var system = true
	if listViper.GetString("system") == "false" {
		system = false
	}

	return client.ListChannel(system)
}