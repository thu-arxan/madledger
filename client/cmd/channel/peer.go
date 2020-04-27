// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package channel

import (
	"errors"
	"madledger/client/lib"
	"madledger/client/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	peerCmd = &cobra.Command{
		Use: "peer",
	}
	peerViper = viper.New()
)

func init() {
	peerCmd.RunE = runPeer
	peerCmd.Flags().StringP("name", "n", "", "the name of channel")
	peerViper.BindPFlag("name", peerCmd.Flags().Lookup("name"))
	peerCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	peerViper.BindPFlag("config", peerCmd.Flags().Lookup("config"))
}

func runPeer(cmd *cobra.Command, args []string) error {
	name := peerViper.GetString("name")
	if name == "" {
		return errors.New("The name of channel should be [a-z0-9]{1,32} such as test, test01 and etc")
	}
	cfgFile := peerViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}
	peers, err := client.GetPeerAddress(name)
	if err != nil {
		return err
	}
	table := util.NewTable()
	table.SetHeader("peer")
	for _, peer := range peers {
		table.AddRow(peer)
	}
	table.Render()

	return nil
}
