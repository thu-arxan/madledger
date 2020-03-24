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
	"madledger/common/abi"
	"madledger/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	callCmd = &cobra.Command{
		Use: "call",
	}
	callViper = viper.New()
)

func init() {
	callCmd.RunE = runCall
	callCmd.Flags().StringP("abi", "a", "", "The abi of tx")
	callViper.BindPFlag("abi", callCmd.Flags().Lookup("abi"))
	callCmd.Flags().StringP("func", "f", "", "The func of contract")
	callViper.BindPFlag("func", callCmd.Flags().Lookup("func"))
	callCmd.Flags().StringSliceP("inputs", "i", nil, "The inputs of function")
	callViper.BindPFlag("inputs", callCmd.Flags().Lookup("inputs"))
	callCmd.Flags().StringP("config", "c", "explorer-client.yaml", "The config file of client")
	callViper.BindPFlag("config", callCmd.Flags().Lookup("config"))
	callCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	callViper.BindPFlag("channelID", callCmd.Flags().Lookup("channelID"))
	callCmd.Flags().StringP("receiver", "r", "", "The contract address of the tx")
	callViper.BindPFlag("receiver", callCmd.Flags().Lookup("receiver"))
	callCmd.Flags().Int64P("value", "v", 0, "The value of the tx")
	callViper.BindPFlag("value", callCmd.Flags().Lookup("value"))
}

func runCall(cmd *cobra.Command, args []string) error {
	cfgFile := callViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	channelID := callViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}
	abiPath := callViper.GetString("abi")
	if abiPath == "" {
		return errors.New("The abi path can not be nil")
	}
	funcName := callViper.GetString("func")
	if funcName == "" {
		return errors.New("The name of func can not be nil")
	}
	receiver := callViper.GetString("receiver")
	value := uint64(callViper.GetInt64("value"))
	if receiver == "" {
		return errors.New("The address of receiver can not be nil")
	}
	inputs := callViper.GetStringSlice("inputs")
	payloadBytes, err := abi.Pack(abiPath, funcName, inputs...)
	if err != nil {
		return err
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	tx, err := core.NewTx(channelID, common.HexToAddress(receiver), payloadBytes, value, "", client.GetPrivKey())
	if err != nil {
		return err
	}
	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}

	table := util.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "Output")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		values, err := abi.Unpack(abiPath, funcName, status.Output)
		if err != nil {
			return err
		}
		table.AddRow(status.BlockNumber, status.BlockIndex, values)
	}
	table.Render()

	return nil
}
