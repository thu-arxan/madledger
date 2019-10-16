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
	addCmd.Flags().StringP("url", "u", "127.0.0.1:45678",
		"The url of node joining the exiting cluster")
	addViper.BindPFlag("url", addCmd.Flags().Lookup("url"))
	addCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	addViper.BindPFlag("config", addCmd.Flags().Lookup("config"))
	addCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	addViper.BindPFlag("channelID", addCmd.Flags().Lookup("channelID"))
}

func runAdd(cmd *cobra.Command, args []string) error {
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
	tx, err := coreTypes.NewTx(channelID, coreTypes.ValidatorUpdateAddress, validatorUpdate,
		client.GetPrivKey(), coreTypes.NODE)
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
