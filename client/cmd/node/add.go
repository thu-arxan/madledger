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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/raft/raftpb"
)

var (
	addCmd = &cobra.Command{
		Use: "add",
	}
	addViper = viper.New()
)

func init() {
	addCmd.RunE = runAdd
	addCmd.Flags().StringP("nodeID", "i", "4",
		"The ID of node joining the exiting cluster")
	addViper.BindPFlag("nodeID", addCmd.Flags().Lookup("nodeID"))
	addCmd.Flags().StringP("url", "u", "127.0.0.1:45680",
		"The url of node joining the exiting cluster")
	addViper.BindPFlag("url", addCmd.Flags().Lookup("url"))
	addCmd.Flags().StringP("config", "c", "explorer-client.yaml", "The config file of client")
	addViper.BindPFlag("config", addCmd.Flags().Lookup("config"))
	addCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	addViper.BindPFlag("channelID", addCmd.Flags().Lookup("channelID"))
}

func runAdd(cmd *cobra.Command, args []string) error {
	nodeID, err := strconv.ParseUint(addViper.GetString("nodeID"), 10, 64)
	if err != nil {
		return err
	}
	if nodeID <= 0 {
		return errors.New("The ID must be bigger than zero")
	}

	urlRaw := addViper.GetString("url")
	if !strings.Contains(urlRaw, ":") {
		return errors.New("The url of node must contains ip and port like 127.0.0.1:12345")
	}
	port, err := strconv.ParseUint(strings.Split(urlRaw, ":")[1], 10, 64)
	if err != nil {
		return err
	}
	if port > 65535 {
		return errors.New("The port can not be bigger than 65535")
	}

	cfgFile := addViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	channelID := addViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}

	// construct ConfChange
	cc, err := json.Marshal(raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  nodeID,
		Context: []byte(urlRaw),
	})
	if err != nil {
		return err
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}
	tx, err := coreTypes.NewTx(channelID, coreTypes.CfgRaftAddress, cc, 0, "", client.GetPrivKey())
	if err != nil {
		return err
	}
	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}
	// Then print the status
	table := util.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "NodeAddOK")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, "ok")
	}
	table.Render()
	return nil
}
