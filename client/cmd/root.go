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
	"madledger/client/cmd/asset"
	"madledger/client/cmd/account"
	"madledger/client/cmd/channel"
	"madledger/client/cmd/node"
	"madledger/client/cmd/tx"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "client", "package": "cmd"})
)

var (
	rootCmd = &cobra.Command{
		Use:  "client",
		Long: "This is the cli of MadLedger client.",
	}
	homeDir, _ = os.Getwd()
	// ordererHome, _ = util.MakeFileAbs(".orderer", homeDir)
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)

	rootCmd.AddCommand(channel.Cmd())
	rootCmd.AddCommand(tx.Cmd())
	rootCmd.AddCommand(account.Cmd())
	rootCmd.AddCommand(node.Cmd())
	rootCmd.AddCommand(asset.Cmd())
}

// Execute exec the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func setLog(debug bool) error {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	return nil
}
