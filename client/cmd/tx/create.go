package tx

import (
	"errors"
	"madledger/client/lib"
	"madledger/common"
	"madledger/core/types"

	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	createCmd = &cobra.Command{
		Use: "create",
	}
	createViper = viper.New()
)

func init() {
	createCmd.RunE = runCreate
	createCmd.Flags().StringP("bin", "b", "", "The bin of tx")
	createViper.BindPFlag("bin", createCmd.Flags().Lookup("bin"))
	createCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	createViper.BindPFlag("config", createCmd.Flags().Lookup("config"))
	createCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	createViper.BindPFlag("channelID", createCmd.Flags().Lookup("channelID"))
}

type createTxStatus struct {
	BlockNumber     uint64
	BlockIndex      int32
	ContractAddress string
}

func runCreate(cmd *cobra.Command, args []string) error {
	cfgFile := createViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	channelID := createViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}
	binPath := createViper.GetString("bin")
	if binPath == "" {
		return errors.New("The bin path can not be nil")
	}
	contractCodes, err := readCodes(binPath)
	if err != nil {
		return err
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	tx, err := types.NewTx(channelID, common.ZeroAddress, contractCodes, client.GetPrivKey())
	if err != nil {
		return err
	}
	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}
	// Then print the status
	table.Output([]createTxStatus{createTxStatus{
		BlockNumber:     status.BlockNumber,
		BlockIndex:      status.BlockIndex,
		ContractAddress: status.ContractAddress,
	}})

	return nil
}
