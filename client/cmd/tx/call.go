package tx

import (
	"encoding/hex"
	"errors"
	"madledger/client/lib"
	"madledger/common"
	"madledger/common/abi"
	"madledger/core/types"

	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	callCmd = &cobra.Command{
		Use: "call",
	}
	callViper = viper.New()
)

func init() {
	callCmd.RunE = runcall
	callCmd.Flags().StringP("abi", "a", "", "The abi of tx")
	callViper.BindPFlag("abi", callCmd.Flags().Lookup("abi"))
	callCmd.Flags().StringP("func", "f", "", "The func of contract")
	callViper.BindPFlag("func", callCmd.Flags().Lookup("func"))
	callCmd.Flags().StringP("payload", "p", "", "The payload of tx")
	callViper.BindPFlag("payload", callCmd.Flags().Lookup("payload"))
	callCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	callViper.BindPFlag("config", callCmd.Flags().Lookup("config"))
	callCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	callViper.BindPFlag("channelID", callCmd.Flags().Lookup("channelID"))
	callCmd.Flags().StringP("receiver", "r", "", "The contract address of the tx")
	callViper.BindPFlag("receiver", callCmd.Flags().Lookup("receiver"))
}

type callTxStatus struct {
	BlockNumber uint64
	BlockIndex  int32
	Output      []string
}

func runcall(cmd *cobra.Command, args []string) error {
	cfgFile := callViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	channelID := callViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}
	abiPath := callViper.GetString("abi")
	if abiPath == "" {
		return errors.New("The abi path can not be nil")
	}
	funcName := callViper.GetString("func")
	if funcName == "" {
		return errors.New("The name of func can not be nil")
	}
	receiver := callViper.GetString("receiver")
	if receiver == "" {
		return errors.New("The address of receiver can not be nil")
	}
	payload := callViper.GetString("payload")
	if payload == "" {
		return errors.New("The payload of client can not be nil")
	}
	payloadBytes, err := hex.DecodeString(payload)
	if err != nil {
		return err
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	tx, err := types.NewTx(channelID, common.HexToAddress(receiver), payloadBytes, client.GetPrivKey())
	if err != nil {
		return err
	}
	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}
	values, err := abi.Unpacker(abiPath, funcName, status.Output)
	if err != nil {
		return err
	}
	var callStatus = callTxStatus{
		BlockNumber: status.BlockNumber,
		BlockIndex:  status.BlockIndex,
	}

	for _, value := range values {
		callStatus.Output = append(callStatus.Output, value.Value)
	}
	table.Output([]callTxStatus{callStatus})

	return nil
}