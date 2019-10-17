package node

import (
	"encoding/base64"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"errors"
	"github.com/tendermint/tendermint/abci/types"
	coreTypes "madledger/core/types"
	"madledger/client/lib"
	"madledger/client/util"
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
		"The ID of node leaving the exiting cluster")
	removeViper.BindPFlag("nodeID", addCmd.Flags().Lookup("nodeID"))
	removeCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	removeViper.BindPFlag("config", addCmd.Flags().Lookup("config"))
	removeCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	removeViper.BindPFlag("channelID", addCmd.Flags().Lookup("channelID"))
}

func runRemove(cmd *cobra.Command, args []string) error {
	dataS := addViper.GetString("pubkey")
	if dataS == "" {
		return errors.New("The pubkey.data of validator can not be nil")
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
	tx, err := coreTypes.NewTx(channelID, coreTypes.CfgRaftAddress, validatorUpdate,
		client.GetPrivKey(), coreTypes.VALIDATOR)
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
