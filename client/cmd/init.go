// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package cmd

import (
	"errors"
	"io/ioutil"
	cutil "madledger/client/util"
	"madledger/common/crypto"
	"madledger/common/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	initCmd = &cobra.Command{
		Use: "init",
	}
	initViper = viper.New()
)

func init() {
	initCmd.RunE = runInit
	initCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	initViper.BindPFlag("config", initCmd.Flags().Lookup("config"))
	initCmd.Flags().StringP("keyAlgo", "k", "sm2", "Crypto of private key, secp256k1 or sm2")
	initViper.BindPFlag("keyAlgo", initCmd.Flags().Lookup("keyAlgo"))
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	var err error
	cfgFile := initViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	keyStorePath, _ := util.MakeFileAbs(".keystore", homeDir)
	err = os.MkdirAll(keyStorePath, 0777)
	if err != nil {
		return err
	}
	var algo crypto.Algorithm
	switch initViper.GetString("keyAlgo") {
	case "secp256k1":
		algo = crypto.KeyAlgoSecp256k1
	default:
		algo = crypto.KeyAlgoSM2
	}
	keyPath, err := cutil.GeneratePrivateKey(keyStorePath, algo)
	if err != nil {
		return err
	}
	var cfg = cfgTemplate
	cfg = strings.Replace(cfg, "<<<KEYFILE>>>", keyPath, 1)

	return ioutil.WriteFile(cfgFile, []byte(cfg), 0777)
}
