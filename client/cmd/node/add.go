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
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/tendermint/tendermint/abci/types"
	//"google.golang.org/grpc/status"
	"madledger/client/lib"
	"madledger/client/util"
	"madledger/common/crypto"
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
	addCmd.Flags().StringP("pubkey", "k", "", "The pubkey of validator")
	addViper.BindPFlag("pubkey", addCmd.Flags().Lookup("pubkey"))
	addCmd.Flags().StringP("power", "p", "10", "The power of validator")
	addViper.BindPFlag("power", addCmd.Flags().Lookup("power"))
	addCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	addViper.BindPFlag("config", addCmd.Flags().Lookup("config"))
}

func runAdd(cmd *cobra.Command, args []string) error {
	isRaft := false
	isTendermint := false
	if addViper.GetString("pubkey") != "" {
		isTendermint = true
	}
	if addViper.GetString("nodeID") != "" {
		isRaft = true
	}

	if isRaft && isTendermint {
		return errors.New("Too many arguments. Please specify consensus type.")
	}
	if !isRaft && !isTendermint {
		return errors.New("Too few arguments. Please specify consensus type.")
	}

	cfgFile := addViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	var tx *coreTypes.Tx
	if isRaft {
		tx, err = getRaftConfChangeTx(client.GetPrivKey())
	} else {
		tx, err = getTendermintConfChangeTx(client.GetPrivKey())
	}
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

func getRaftConfChangeTx(privkey crypto.PrivateKey) (*coreTypes.Tx, error){
	nodeID, err := strconv.ParseUint(addViper.GetString("nodeID"), 10, 64)
	if err != nil {
		return nil, err
	}
	if nodeID <= 0 {
		return nil, errors.New("The ID must be bigger than zero")
	}

	urlRaw := addViper.GetString("url")
	if !strings.Contains(urlRaw, ":") {
		return nil, errors.New("The url of node must contains ip and port like 127.0.0.1:12345")
	}
	port, err := strconv.ParseUint(strings.Split(urlRaw, ":")[1], 10, 64)
	if err != nil {
		return nil, err
	}
	if port > 65535 {
		return nil, errors.New("The port can not be bigger than 65535")
	}

	// construct ConfChange
	cc, err := json.Marshal(raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  nodeID,
		Context: []byte(urlRaw),
	})
	if err != nil {
		return nil, err
	}

	tx, err := coreTypes.NewTx(coreTypes.CONFIGCHANNELID, coreTypes.CfgConsensusAddress, cc, 0, "", privkey)
	return tx, err
}

func getTendermintConfChangeTx(privKey crypto.PrivateKey) (*coreTypes.Tx, error){
	dataS := addViper.GetString("pubkey")
	if dataS == "" {
		return nil, errors.New("Tendermint pubkey cannot be empty")
	}
	// construct PubKey
	data, err := base64.StdEncoding.DecodeString(dataS)
	if err != nil {
		return nil, err
	}

	pubkey := types.PubKey{
		Type: "ed25519",
		Data: data,
	}

	power := addViper.GetInt64("power")
	if power < 0 {
		return nil, errors.New("The power of validator power must be non-negative")
	}

	// construct ValidatorUpdate
	validatorUpdate, err := json.Marshal(types.ValidatorUpdate{
		PubKey: pubkey,
		Power:  power,
	})
	if err != nil {
		return nil, err
	}

	tx, err := coreTypes.NewTx(coreTypes.CONFIGCHANNELID, coreTypes.CfgConsensusAddress, validatorUpdate, 0, "", privKey)
	return tx, err
}