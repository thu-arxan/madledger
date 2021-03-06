// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package node

import (
	"encoding/json"
	"errors"
	"madledger/client/lib"
	"madledger/client/util"
	coreTypes "madledger/core"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/raft/raftpb"
)

var (
	removeCmd = &cobra.Command{
		Use: "remove",
	}
	removeViper = viper.New()
)

func init() {
	removeCmd.RunE = runRemove
	removeCmd.Flags().StringP("nodeID", "i", "4",
		"The ID of node you want to remove from the exiting cluster")
	removeViper.BindPFlag("nodeID", removeCmd.Flags().Lookup("nodeID"))
	removeCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	removeViper.BindPFlag("config", removeCmd.Flags().Lookup("config"))
	removeCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	removeViper.BindPFlag("channelID", removeCmd.Flags().Lookup("channelID"))
}

func runRemove(cmd *cobra.Command, args []string) error {
	nodeID, err := strconv.ParseUint(removeViper.GetString("nodeID"), 10, 64)
	if err != nil {
		return err
	}
	if nodeID <= 0 {
		return errors.New("The ID must be bigger than zero")
	}

	cfgFile := removeViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	channelID := removeViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}

	// construct ConfChange
	cc, err := json.Marshal(raftpb.ConfChange{
		Type:   raftpb.ConfChangeRemoveNode,
		NodeID: nodeID,
	})
	if err != nil {
		return err
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}
	tx, err := coreTypes.NewTx(channelID, coreTypes.CfgConsensusAddress, cc, 0, "", client.GetPrivKey())
	if err != nil {
		return err
	}
	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}
	// Then print the status
	table := util.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "NodeRemoveOK")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, "ok")
	}
	table.Render()
	return nil
}
