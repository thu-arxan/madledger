package db

import (
	"encoding/hex"
	"errors"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	secp256k1String      = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	rawSecp256k1Bytes, _ = hex.DecodeString(secp256k1String)
	rawPrivKey           = rawSecp256k1Bytes
	privKey, _           = crypto.NewPrivateKey(rawPrivKey)
)

var (
	dir = ".leveldb"
	db  DB
)

func TestInit(t *testing.T) {
	err := os.RemoveAll(dir)
	require.NoError(t, err)

	err = os.MkdirAll(dir, 0777)
	require.NoError(t, err)
}

func TestNewLevelDB(t *testing.T) {
	var err error
	db, err = NewLevelDB(dir)
	require.NoError(t, err)
}

func TestListChannel(t *testing.T) {
	channels := db.ListChannel()
	require.Len(t, channels, 0)
}

func TestUpdateChannel(t *testing.T) {
	err := db.UpdateChannel("_config", &cc.Profile{
		Public: true,
	})
	require.NoError(t, err)
	var channels []string
	channels = db.ListChannel()
	require.Len(t, channels, 1)
	require.Equal(t, channels[0], "_config")
	// add _global
	err = db.UpdateChannel("_global", &cc.Profile{
		Public: true,
	})
	require.NoError(t, err)

	channels = db.ListChannel()
	require.Len(t, channels, 2)
	if !util.Contain(channels, "_global") {
		t.Fatal(errors.New("Channel _global is not contained"))
	}
	// add user channel
	admin, _ := core.NewMember(privKey.PubKey(), "admin")
	err = db.UpdateChannel("test", &cc.Profile{
		Public: true,
		Admins: []*core.Member{admin},
	})
	require.NoError(t, err)
	channels = db.ListChannel()
	require.Len(t, channels, 3)
	if !util.Contain(channels, "test") {
		t.Fatal(errors.New("Channel test is not contained"))
	}

	// add _asset
	err = db.UpdateChannel("_asset", &cc.Profile{
		Public: true,
	})
	require.NoError(t, err)

	channels = db.ListChannel()
	require.Len(t, channels, 4)
	if !util.Contain(channels, "_asset") {
		t.Fatal(errors.New("Channel _global is not contained"))
	}
}

func TestAddBlock(t *testing.T) {
	tx1, _ := core.NewTx("test", common.ZeroAddress, []byte("1"), 0, "", privKey)
	tx2, _ := core.NewTx("test", common.ZeroAddress, []byte("2"), 0, "", privKey)
	block1 := core.NewBlock("test", 0, core.GenesisBlockPrevHash, []*core.Tx{tx1, tx2})
	err := db.AddBlock(block1)
	require.NoError(t, err)

	block2 := core.NewBlock("test", 1, block1.Hash().Bytes(), []*core.Tx{tx1})
	err = db.AddBlock(block2)
	require.Error(t, err)

	if !db.HasTx(tx1) || !db.HasTx(tx2) {
		t.Fatal()
	}
}

func TestIsMember(t *testing.T) {
	member, _ := core.NewMember(privKey.PubKey(), "admin")
	require.True(t, db.IsMember("test", member))
	require.True(t, db.IsMember(core.GLOBALCHANNELID, member))
}

func TestIsAdmin(t *testing.T) {
	member, _ := core.NewMember(privKey.PubKey(), "admin")
	require.True(t, db.IsAdmin("test", member))
	require.False(t, db.IsAdmin(core.GLOBALCHANNELID, member))
}

func TestAssetAdmin(t *testing.T) {
	wb := db.NewWriteBatch()
	require.False(t, db.IsAssetAdmin(privKey.PubKey()))
	err := wb.SetAssetAdmin(privKey.PubKey())
	require.NoError(t, err)
	require.NoError(t, wb.Sync())
	require.True(t, db.IsAssetAdmin(privKey.PubKey()))
}

func TestAccount(t *testing.T) {
	wb := db.NewWriteBatch()
	address := common.BytesToAddress([]byte("12345678"))
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

func TestEnd(t *testing.T) {
	os.RemoveAll(dir)
}
