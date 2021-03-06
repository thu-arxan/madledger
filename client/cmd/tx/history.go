// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tx

import (
	"errors"
	"madledger/client/lib"
	"madledger/client/util"

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

	table := util.NewTable()
	table.SetHeader("Channel", "TxID")
	for channel, txs := range history.Txs {
		for _, id := range txs.Value {
			table.AddRow(channel, id)
		}
	}
	table.Render()
	return nil
}
