// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

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

	info, err := client.GetAccountBalance(address)
	if err != nil {
		return err
	}
	token, err := client.GetTokenInfo(address)
	if err != nil {
		return err
	}
	table := util.NewTable()
	table.SetHeader("Address", "balance", "token")
	table.AddRow(address.String(), info, token)
	table.Render()

	return nil
}
