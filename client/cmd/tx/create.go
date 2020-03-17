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
	"madledger/common"
	"madledger/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	createCmd = &cobra.Command{
		Use: "create",
	}
	createViper = viper.New()
)

func init() {
	createCmd.RunE = runCreate
	createCmd.Flags().StringP("bin", "b", "", "The bin of tx")
	createViper.BindPFlag("bin", createCmd.Flags().Lookup("bin"))
	createCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	createViper.BindPFlag("config", createCmd.Flags().Lookup("config"))
	createCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	createViper.BindPFlag("channelID", createCmd.Flags().Lookup("channelID"))
	createCmd.Flags().Int64P("value", "v", 0, "The value of the tx")
	createViper.BindPFlag("value", createCmd.Flags().Lookup("value"))
}

type createTxStatus struct {
	BlockNumber     uint64
	BlockIndex      int32
	ContractAddress string
}

func runCreate(cmd *cobra.Command, args []string) error {
	cfgFile := createViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	channelID := createViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}
	value := uint64(createViper.GetInt64("value"))
	binPath := createViper.GetString("bin")
	if binPath == "" {
		return errors.New("The bin path can not be nil")
	}
	contractCodes, err := readCodes(binPath)
	if err != nil {
		return err
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	tx, err := core.NewTx(channelID, common.ZeroAddress, contractCodes, value, "", client.GetPrivKey())
	if err != nil {
		return err
	}
	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}
	// Then print the status
	table := util.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "ContractAddress")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.ContractAddress)
	}
	table.Render()

	return nil
}
