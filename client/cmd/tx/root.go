// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tx

import (
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var (
	homeDir, _ = os.Getwd()
	txCmd      = &cobra.Command{
		Use: "tx",
	}
)

// Cmd return the tx command
func Cmd() *cobra.Command {
	txCmd.AddCommand(createCmd)
	txCmd.AddCommand(callCmd)
	txCmd.AddCommand(historyCmd)
	return txCmd
}

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}
