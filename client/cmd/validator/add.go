// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package validator

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"madledger/client/lib"
	"madledger/client/util"
	"madledger/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/abci/types"
)

var (
	addCmd = &cobra.Command{
		Use: "add",
	}
	addViper = viper.New()
)

func init() {
	addCmd.RunE = runAdd
	addCmd.Flags().StringP("pubkey", "k", "", "The pubkey of validator")
	addViper.BindPFlag("pubkey", addCmd.Flags().Lookup("pubkey"))
	addCmd.Flags().StringP("power", "p", "10", "The power of validator")
	addViper.BindPFlag("power", addCmd.Flags().Lookup("power"))
	addCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	addViper.BindPFlag("config", addCmd.Flags().Lookup("config"))
	addCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	addViper.BindPFlag("channelID", addCmd.Flags().Lookup("channelID"))
}

func runAdd(cmd *cobra.Command, args []string) error {
	dataS := addViper.GetString("pubkey")
	if dataS == "" {
		return errors.New("The pubkey of validator can not be nil")
	}
	// construct PubKey
	data, err := base64.StdEncoding.DecodeString(dataS)
	if err != nil {
		return err
	}
	pubkey := types.PubKey{
		Type: "ed25519",
		Data: data,
	}

	power := addViper.GetInt64("power")
	if power < 0 {
		return errors.New("The power of validator power must be non-negative")
	}

	cfgFile := addViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	channelID := addViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}
	// construct ValidatorUpdate
	validatorUpdate, err := json.Marshal(types.ValidatorUpdate{
		PubKey: pubkey,
		Power:  power,
	})
	if err != nil {
		return err
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}
	tx, err := core.NewTx(channelID, core.CfgTendermintAddress, validatorUpdate, 0, "", client.GetPrivKey())
	if err != nil {
		return err
	}
	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}
	// Then print the status
	table := util.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "ValidatorAddOk")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, "ok")
	}
	table.Render()
	return nil
}
