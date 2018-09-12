package tx

import (
	"encoding/hex"
	"errors"
	"madledger/client/orderer"
	"madledger/common"
	"madledger/core/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sendCmd = &cobra.Command{
		Use: "send",
	}
	sendViper = viper.New()
)

func init() {
	sendCmd.RunE = runsend
	sendCmd.Flags().StringP("payload", "p", "", "The payload of tx")
	sendViper.BindPFlag("payload", sendCmd.Flags().Lookup("payload"))
	sendCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	sendViper.BindPFlag("config", sendCmd.Flags().Lookup("config"))
	sendCmd.Flags().StringP("receiver", "r", "", "The receiver of the tx")
	sendViper.BindPFlag("receiver", sendCmd.Flags().Lookup("receiver"))
	sendCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	sendViper.BindPFlag("channelID", sendCmd.Flags().Lookup("channelID"))
}

func runsend(cmd *cobra.Command, args []string) error {
	cfgFile := sendViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	payload := sendViper.GetString("payload")
	if payload == "" {
		return errors.New("The payload of client can not be nil")
	}
	payloadBytes, err := hex.DecodeString(payload)
	if err != nil {
		return err
	}
	channelID := sendViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}
	receiver := sendViper.GetString("receiver")
	if receiver == "" {
		return errors.New("The address of receiver can not be nil")
	}
	client, err := orderer.NewClient(cfgFile)
	if err != nil {
		return err
	}

	tx, err := types.NewTx(channelID, common.HexToAddress(receiver), []byte(payloadBytes), client.GetPrivKey())
	if err != nil {
		return err
	}
	return client.AddTx(tx)
}
