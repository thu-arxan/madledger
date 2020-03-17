// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tendermint

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
	"os"
)

func TestDB_GetDB(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/src/madledger/env/bft/orderers/0/data/leveldb", gopath)
	db, err := leveldb.OpenFile(path, nil)
	fmt.Printf("Get bft.db from %s\n", path)
	require.NoError(t, err)
	defer db.Close()
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		fmt.Printf("(%s, %s)\n", key, value)
	}
	iter.Release()
}
