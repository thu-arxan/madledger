package account

import (
	"errors"
	"madledger/client/lib"
	"madledger/client/util"

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
	listCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	listViper.BindPFlag("config", listCmd.Flags().Lookup("config"))
}

func runList(cmd *cobra.Command, args []string) error {
	cfgFile := listViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	address, err := client.GetPrivKey().PubKey().Address()
	if err != nil {
		return err
	}

	//todo: ab can get other info?
	info, err := client.GetAccountBalance(address)
	if err != nil {
		return err
	}
	table := util.NewTable()
	table.SetHeader("Address", "balance")
	table.AddRow(address.String(), info)
	table.Render()

	return nil
}
