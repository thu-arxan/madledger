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
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "cmd"})
)

var (
	rootCmd = &cobra.Command{
		Use:  "peer",
		Long: "This is the cli of MadLedger peer.",
	}
	homeDir, _ = os.Getwd()
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)
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

func registerStop() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	return
}
