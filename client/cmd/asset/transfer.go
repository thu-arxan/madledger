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
	transferCmd.Flags().StringP("address", "a", "0", "receiver's hex address to be transfered")
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

	receiver := transferViper.GetString("address")
	recipient := common.HexToAddress(receiver)

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(asset.Payload{
		Action:    "transfer",
		ChannelID: channelID,
	})
	if err != nil {
		return err
	}

	tx, err := coreTypes.NewTx(coreTypes.ASSETCHANNELID, recipient, payload, uint64(value), "", client.GetPrivKey())
	if err != nil {
		return err
	}

	status, err := client.AddTxInOrderer(tx)
	table := util.NewTable()
	table.SetHeader("Status", "Error")

	table.AddRow(status, err)

	table.Render()

	return err
}