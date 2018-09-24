package tx

import (
	"errors"
	"madledger/client/lib"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	historyCmd = &cobra.Command{
		Use: "history",
	}
	historyViper = viper.New()
)

func init() {
	historyCmd.RunE = runHistory
	historyCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	historyViper.BindPFlag("config", historyCmd.Flags().Lookup("config"))
}

func runHistory(cmd *cobra.Command, args []string) error {
	cfgFile := historyViper.GetString("config")
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

	history, err := client.GetHistory(address.Bytes())
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Channel", "TxID"})
	for channel, txs := range history.Txs {
		for _, id := range txs.Value {
			table.Append([]string{channel, id})
		}
	}
	table.Render()
	return nil
}
