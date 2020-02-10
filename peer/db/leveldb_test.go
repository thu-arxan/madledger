package db

import (
	"encoding/hex"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	"os"
	"reflect"
	"testing"
	"time"

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

var (
	tx1, _    = core.NewTx("test", common.ZeroAddress, []byte("1"), 0, "", privKey)
	tx1Status = &TxStatus{
		Err:             "",
		BlockNumber:     3,
		BlockIndex:      1,
		Output:          []byte("tx1"),
		ContractAddress: "",
	}
	tx2, _    = core.NewTx("test", common.ZeroAddress, []byte("2"), 0, "", privKey)
	tx2Status = &TxStatus{
		Err:             "",
		BlockNumber:     3,
		BlockIndex:      1,
		Output:          []byte("tx2"),
		ContractAddress: "",
	}
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

func TestAccount(t *testing.T) {
	address, err := privKey.PubKey().Address()
	require.NoError(t, err)
	// The address should not exist
	require.False(t, db.AccountExist(address))
	// But if we GetAccount, we can get the default account
	account, err := db.GetAccount(address)
	require.NoError(t, err)
	defaultAccount := common.NewDefaultAccount(address)
	require.Equal(t, defaultAccount, account)
	// then set balance and code
	account.AddBalance(100)
	require.Equal(t, uint64(100), account.GetBalance())
	code := []byte("Hello world")
	account.SetCode(code)
	require.Equal(t, code, account.GetCode())
	// the set the account
	err = db.SetAccount(account)
	require.NoError(t, err)
	account, err = db.GetAccount(address)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(account.GetAddress().Bytes(), address.Bytes()))
	require.Equal(t, uint64(100), account.GetBalance())
	require.Equal(t, code, account.GetCode())
	require.True(t, db.AccountExist(account.GetAddress()))
	// then remove account
	err = db.RemoveAccount(account.GetAddress())
	require.NoError(t, err)
	require.False(t, db.AccountExist(account.GetAddress()))
}

func TestStorage(t *testing.T) {
	// first set an account
	address, _ := privKey.PubKey().Address()
	account, _ := db.GetAccount(address)
	db.SetAccount(account)
	// then get key and value
	key, err := common.BytesToWord256([]byte("I want a key which length is 32."))
	require.NoError(t, err)
	value, err := common.BytesToWord256([]byte("I need a value that length is 32"))
	require.NoError(t, err)
	// then test the storage
	_, err = db.GetStorage(address, key)
	require.Error(t, err, "not found")
	err = db.SetStorage(address, key, value)
	require.NoError(t, err)
	v, err := db.GetStorage(address, key)
	require.NoError(t, err)
	require.Equal(t, value, v)
}

func TestSetTxStatus(t *testing.T) {
	err := db.SetTxStatus(tx1, tx1Status)
	require.NoError(t, err)
}

func TestGetTxStatus(t *testing.T) {
	status, err := db.GetTxStatus(tx1.Data.ChannelID, tx1.ID)
	require.NoError(t, err)
	require.Equal(t, status, tx1Status)

	_, err = db.GetTxStatus(tx2.Data.ChannelID, tx2.ID)
	require.Error(t, err, "Not exist")
}

func TestGetTxStatusAsync(t *testing.T) {
	_, err := db.GetTxStatusAsync(tx1.Data.ChannelID, tx1.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then try to get tx2 status async
	var endChan = make(chan bool)
	go func() {
		defer func() {
			endChan <- true
		}()
		status, err := db.GetTxStatusAsync(tx2.Data.ChannelID, tx2.ID)
		require.NoError(t, err)
		if !reflect.DeepEqual(status, tx2Status) {
			t.Fatal()
		}
	}()
	time.Sleep(100 * time.Millisecond)
	go func() {
		err := db.SetTxStatus(tx2, tx2Status)
		require.NoError(t, err)
	}()
	<-endChan
}

func TestListTxHistory(t *testing.T) {
	address, err := privKey.PubKey().Address()
	if err != nil {
		t.Fatal(err)
	}
	history := db.ListTxHistory(address.Bytes())
	expectHistory := make(map[string][]string)
	expectHistory["test"] = []string{tx1.ID, tx2.ID}
	if !reflect.DeepEqual(history, expectHistory) {
		t.Fatal()
	}
}

func TestEnd(t *testing.T) {
	os.RemoveAll(dir)
}
