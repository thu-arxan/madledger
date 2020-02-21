package asset

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"madledger/blockchain/asset"
	"madledger/client/lib"
	"madledger/common"
	coreTypes "madledger/core"
)

var (
	transferCmd = &cobra.Command{
		Use:"transfer",
	}
	transferViper = viper.New()
)

func init() {
	transferCmd.RunE = runTransfer
	transferCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	transferViper.BindPFlag("config", transferCmd.Flags().Lookup("config"))
	transferCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	transferViper.BindPFlag("channelID", transferCmd.Flags().Lookup("channelID"))
	transferCmd.Flags().StringP("value", "v", "0",
		"value to be transfered")
	transferViper.BindPFlag("value", transferCmd.Flags().Lookup("value"))
	transferCmd.Flags().StringP("address", "a", "0",
		"receiver's hex address to be transfered")
	transferViper.BindPFlag("address", transferCmd.Flags().Lookup("address"))
}

func runTransfer(cmd *cobra.Command, args []string) error {
	cfgFile := transferViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	channelID := transferViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}

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
		Action: "transfer",
		ChannelID: channelID,
	})
	if err != nil {
		return err
	}

	tx, err := coreTypes.NewTx(coreTypes.ASSETCHANNELID, recipient, payload, uint64(value), "", client.GetPrivKey())
	if err != nil {
		return err
	}

	_, err = client.AddTx(tx)

	return err
}