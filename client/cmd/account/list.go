package account

import (
	"errors"
	"madledger/client/lib"
	"os"

	"github.com/olekukonko/tablewriter"
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
	listCmd.Flags().StringP("language", "l", "en", "The language of client")
	listViper.BindPFlag("language", listCmd.Flags().Lookup("language"))
}

func runList(cmd *cobra.Command, args []string) error {
	cfgFile := listViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	language := listViper.GetString("language")
	if language != "zh" {
		language = "en"
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	address, err := client.GetPrivKey().PubKey().Address()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	if language == "en" {
		table.SetHeader([]string{"Address"})
	} else {
		table.SetHeader([]string{"地址"})
	}
	table.Append([]string{address.String()})
	table.Render()

	return nil
}
