package db

import (
	"encoding/hex"
	"errors"
	"fmt"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core/types"
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
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewLevelDB(t *testing.T) {
	var err error
	db, err = NewLevelDB(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListChannel(t *testing.T) {
	channels := db.ListChannel()
	if len(channels) != 0 {
		t.Fatal()
	}
}

func TestUpdateChannel(t *testing.T) {
	err := db.UpdateChannel("_config", &cc.Profile{
		Public: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	var channels []string
	channels = db.ListChannel()
	if len(channels) != 1 {
		t.Fatal(fmt.Errorf("Should contain one channel rather than %d channels", len(channels)))
	}
	if channels[0] != "_config" {
		t.Fatal(fmt.Errorf("Should contain channel _config rather than %s", channels[0]))
	}
	// add _global
	err = db.UpdateChannel("_global", &cc.Profile{
		Public: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	channels = db.ListChannel()
	if len(channels) != 2 {
		t.Fatal(fmt.Errorf("Should contain two channels rather than %d channels", len(channels)))
	}
	if !util.Contain(channels, "_global") {
		t.Fatal(errors.New("Channel _global is not contained"))
	}
	// add user channel
	admin, _ := types.NewMember(privKey.PubKey(), "admin")
	err = db.UpdateChannel("test", &cc.Profile{
		Public: true,
		Admins: []*types.Member{admin},
	})
	require.NoError(t, err)
	channels = db.ListChannel()
	if len(channels) != 3 {
		t.Fatal(fmt.Errorf("Should contain three channels rather than %d channels", len(channels)))
	}
	if !util.Contain(channels, "test") {
		t.Fatal(errors.New("Channel test is not contained"))
	}
	// todo: maybe illegal channel id

}

func TestAddBlock(t *testing.T) {
	tx1, _ := types.NewTx("test", common.ZeroAddress, []byte("1"), privKey)
	tx2, _ := types.NewTx("test", common.ZeroAddress, []byte("2"), privKey)
	block1 := types.NewBlock("test", 0, types.GenesisBlockPrevHash, []*types.Tx{tx1, tx2})
	err := db.AddBlock(block1)
	if err != nil {
		t.Fatal(err)
	}
	block2 := types.NewBlock("test", 1, block1.Hash().Bytes(), []*types.Tx{tx1})
	err = db.AddBlock(block2)
	if err == nil {
		t.Fatal()
	}
	if !db.HasTx(tx1) || !db.HasTx(tx2) {
		t.Fatal()
	}
}

func TestIsMember(t *testing.T) {
	member, _ := types.NewMember(privKey.PubKey(), "admin")
	require.True(t, db.IsMember("test", member))
	require.True(t, db.IsAdmin("test", member))
	require.True(t, db.IsMember(types.GLOBALCHANNELID, member))
	require.False(t, db.IsAdmin(types.GLOBALCHANNELID, member))
}

func TestEnd(t *testing.T) {
	os.RemoveAll(dir)
}
