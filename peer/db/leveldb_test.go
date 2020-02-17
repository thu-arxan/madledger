package db

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TODO: This test shoule be redone.

var (
	secp256k1String      = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	rawSecp256k1Bytes, _ = hex.DecodeString(secp256k1String)
	rawPrivKey           = rawSecp256k1Bytes
	privKey, _           = crypto.NewPrivateKey(rawPrivKey)
	benckmark            = true
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

// TODO: This test should be reconstructed
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
	wb := db.NewWriteBatch()
	require.NoError(t, wb.RemoveAccount(account.GetAddress()))
	require.NoError(t, wb.Sync())
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
	wb := db.NewWriteBatch()
	err = wb.SetStorage(address, key, value)
	require.NoError(t, err)
	require.NoError(t, wb.Sync())
	v, err := db.GetStorage(address, key)
	require.NoError(t, err)
	require.Equal(t, value, v)
}

func TestSetTxStatus(t *testing.T) {
	wb := db.NewWriteBatch()
	require.NoError(t, wb.SetTxStatus(tx1, tx1Status))
	require.NoError(t, wb.Sync())
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
	wb := db.NewWriteBatch()
	go func() {
		err := wb.SetTxStatus(tx2, tx2Status)
		require.NoError(t, err)
	}()
	<-endChan
	require.NoError(t, wb.Sync())
}

func TestListTxHistory(t *testing.T) {
	address, err := privKey.PubKey().Address()
	if err != nil {
		t.Fatal(err)
	}
	history := db.ListTxHistory(address.Bytes())
	exceptHistory := make(map[string][]string)
	exceptHistory["test"] = []string{tx1.ID, tx2.ID}
	if !reflect.DeepEqual(history, exceptHistory) {
		t.Fatal()
	}
}

func TestBenchmark(t *testing.T) {
	if !benckmark {
		return
	}
	var size = 10000
	var accounts = make([]common.Account, size)
	for i := 0; i < size; i++ {
		accounts[i] = newAccount()
	}
	var begin = time.Now()
	for i := 0; i < size; i++ {
		accounts[i].Bytes()
	}
	fmt.Printf("marshal %d accounts cost %v\n", size, time.Since(begin))
	begin = time.Now()
	for i := 0; i < size; i++ {
		MarshalAccount(accounts[i])
	}
	fmt.Printf("fast marshal %d accounts cost %v\n", size, time.Since(begin))
	begin = time.Now()
	bytes, _ := accounts[0].Bytes()
	for i := 0; i < size; i++ {
		var account common.DefaultAccount
		json.Unmarshal(bytes, &account)
	}
	fmt.Printf("unmarshal %d accounts cost %v\n", size, time.Since(begin))
	begin = time.Now()
	bytes = MarshalAccount(accounts[0])
	for i := 0; i < size; i++ {
		UnmarshalAccount(bytes)
	}
	fmt.Printf("fast unmarshal %d accounts cost %v\n", size, time.Since(begin))
	begin = time.Now()
	for i := 0; i < size; i++ {
		db.SetAccount(accounts[i])
	}
	fmt.Printf("add %d accounts cost %v\n", size, time.Since(begin))
	begin = time.Now()
	for i := 0; i < size; i++ {
		db.GetAccount(accounts[i].GetAddress())
	}
	fmt.Printf("get %d accounts cost %v\n", size, time.Since(begin))
}

func TestEnd(t *testing.T) {
	os.RemoveAll(dir)
}

func newAccount() common.Account {
	priv, _ := crypto.GeneratePrivateKey()
	addr, _ := priv.PubKey().Address()
	return common.NewDefaultAccount(addr)
}
