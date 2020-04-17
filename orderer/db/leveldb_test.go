// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package db

import (
	"encoding/hex"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	secp256k1String      = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	rawSecp256k1Bytes, _ = hex.DecodeString(secp256k1String)
	rawPrivKey           = rawSecp256k1Bytes
	privKey, _           = crypto.NewPrivateKey(rawPrivKey, crypto.KeyAlgoSecp256k1)
)

var (
	dir = ".leveldb"
)

func TestLevelDB(t *testing.T) {
	db := initDB(t)
	testChannel(t, db)
	testAddBlock(t, db)
	testIsMember(t, db)
	testIsAdmin(t, db)
	testAssetAdmin(t, db)
	testAccount(t, db)
	testConsensusBlock(t, db)
	os.RemoveAll(dir)
}

func initDB(t *testing.T) DB {
	require.NoError(t, os.RemoveAll(dir))
	require.NoError(t, os.MkdirAll(dir, 0777))
	db, err := NewLevelDB(dir)
	require.NoError(t, err)
	return db
}

func testChannel(t *testing.T, db DB) {
	// size should be 0
	size := len(db.ListChannel())
	require.Equal(t, 0, size)
	// add _config channel
	wb := db.NewWriteBatch()
	require.NoError(t, wb.UpdateChannel("_config", &cc.Profile{
		Public: true,
	}))
	require.NoError(t, wb.Sync())
	size++
	// should contain _config
	channels := db.ListChannel()
	require.Len(t, channels, size)
	require.Contains(t, channels, "_config")
	// add _global and user channel
	wb = db.NewWriteBatch()
	require.NoError(t, wb.UpdateChannel("_global", &cc.Profile{
		Public: true,
	}))
	admin, _ := core.NewMember(privKey.PubKey(), "admin")
	// note: we add test twice, but different profile
	require.NoError(t, wb.UpdateChannel("test", &cc.Profile{
		Public: true,
		Admins: []*core.Member{admin},
	}))
	require.NoError(t, wb.UpdateChannel("test", &cc.Profile{
		Public: false,
		Admins: []*core.Member{admin},
	}))
	require.NoError(t, wb.Sync())
	size += 2
	channels = db.ListChannel()
	require.Len(t, channels, size)
	require.Contains(t, channels, "_global")
	// test channel should be private
	profile, err := db.GetChannelProfile("test")
	require.NoError(t, err)
	require.False(t, profile.Public)
	require.Len(t, profile.Admins, 1)
}

func testAddBlock(t *testing.T, db DB) {
	tx1, _ := core.NewTx("test", common.ZeroAddress, []byte("1"), 0, "", privKey)
	tx2, _ := core.NewTx("test", common.ZeroAddress, []byte("2"), 0, "", privKey)
	block1 := core.NewBlock("test", 0, core.GenesisBlockPrevHash, []*core.Tx{tx1, tx2})
	wb := db.NewWriteBatch()
	err := wb.AddBlock(block1)
	require.NoError(t, err)
	require.NoError(t, wb.Sync())
	block2 := core.NewBlock("test", 1, block1.Hash().Bytes(), []*core.Tx{tx1})
	wb = db.NewWriteBatch()
	require.Error(t, wb.AddBlock(block2))
	require.True(t, db.HasTx(tx1) && db.HasTx(tx2))
}

func testIsMember(t *testing.T, db DB) {
	member, _ := core.NewMember(privKey.PubKey(), "admin")
	require.True(t, db.IsMember("test", member))
	require.True(t, db.IsMember(core.GLOBALCHANNELID, member))
}

func testIsAdmin(t *testing.T, db DB) {
	member, _ := core.NewMember(privKey.PubKey(), "admin")
	require.True(t, db.IsAdmin("test", member))
	require.False(t, db.IsAdmin(core.GLOBALCHANNELID, member))
}

func testAssetAdmin(t *testing.T, db DB) {
	wb := db.NewWriteBatch()
	err := wb.SetAssetAdmin(privKey.PubKey())
	require.NoError(t, err)
	require.NoError(t, wb.Sync())
	pk, err := crypto.NewPublicKey(db.GetAssetAdminPKBytes(), crypto.KeyAlgoSecp256k1)
	require.NoError(t, err)
	require.Equal(t, pk, privKey.PubKey())
}

func testAccount(t *testing.T, db DB) {
	wb := db.NewWriteBatch()
	address := common.BytesToAddress([]byte("channelname"))
	account, err := db.GetOrCreateAccount(address)
	require.NoError(t, err)
	require.Equal(t, account.GetBalance(), uint64(0))
	require.NoError(t, account.AddBalance(10))
	require.NoError(t, wb.UpdateAccounts(account))
	require.NoError(t, wb.Sync())
	account, err = db.GetOrCreateAccount(address)
	require.NoError(t, err)
	require.Equal(t, account.GetBalance(), uint64(10))
}

func testConsensusBlock(t *testing.T, db DB) {
	require.EqualValues(t, 1, db.GetConsensusBlock("test"))
	wb := db.NewWriteBatch()
	wb.SetConsensusBlock("test", 0)
	require.NoError(t, wb.Sync())
	require.EqualValues(t, 1, db.GetConsensusBlock("test"))
	wb = db.NewWriteBatch()
	wb.SetConsensusBlock("test", 2)
	require.NoError(t, wb.Sync())
	require.EqualValues(t, 2, db.GetConsensusBlock("test"))
	wb = db.NewWriteBatch()
	wb.SetConsensusBlock("test", 3)
	wb.SetConsensusBlock("test", 4)
	require.NoError(t, wb.Sync())
	require.EqualValues(t, 4, db.GetConsensusBlock("test"))
}
