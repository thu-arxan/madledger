// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package asset

import (
	"encoding/json"
	"errors"
	"madledger/blockchain/asset"
	"madledger/client/lib"
	"madledger/client/util"
	"madledger/common"
	coreTypes "madledger/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	transferCmd = &cobra.Command{
		Use: "transfer",
	}
	transferViper = viper.New()
)

func init() {
	transferCmd.RunE = runTransfer
	transferCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	transferViper.BindPFlag("config", transferCmd.Flags().Lookup("config"))
	transferCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	transferViper.BindPFlag("channelID", transferCmd.Flags().Lookup("channelID"))
	transferCmd.Flags().StringP("value", "v", "0", "value to be transfered")
	transferViper.BindPFlag("value", transferCmd.Flags().Lookup("value"))
	transferCmd.Flags().StringP("address", "a", "", "receiver's hex address to be transfered")
	transferViper.BindPFlag("address", transferCmd.Flags().Lookup("address"))
}

func runTransfer(cmd *cobra.Command, args []string) error {
	cfgFile := transferViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	//channelID can be empty
	channelID := transferViper.GetString("channelID")

	value := transferViper.GetInt("value")
	if value < 0 {
		return errors.New("cannot issue negative value")
	}
	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}
	receiver := transferViper.GetString("address")
	var recipient common.Address

	if channelID == "" && receiver != "" {
		recipient = common.HexToAddress(receiver)
	} else if channelID != "" && receiver == "" {
		recipient = coreTypes.TransferContractrAddress
	} else {
		return errors.New("only one of channelID and receiver can have value")
	}

	payload, err := json.Marshal(asset.Payload{
		ChannelID: channelID,
	})
	if err != nil {
		return err
	}

	tx, err := coreTypes.NewTx(coreTypes.ASSETCHANNELID, recipient, payload, uint64(value), "", client.GetPrivKey())
	if err != nil {
		return err
	}

	status, err := client.AddTx(tx)
	table := util.NewTable()
	table.SetHeader("Status", "Error")

	table.AddRow(status, err)

	table.Render()

	return err
}
