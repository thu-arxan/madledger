package types

import (
	"encoding/hex"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	rawECDSAPrivKey = `30770201010420d14ff3a516ac13545ae0ac64f28cf5e1121d36f295d982688851ac5ddb8f4032a00a06082a8648ce3d030107a144034200049e441d5fdbeaac2522b258c485dd47c50027d19ef058d19fc18e5f33f926930fd47f251aff754b33e9d53eafc6660abce92026ed721e3430251161c72bee2cd9`
	// rawSM2PrivKey   = `308193020100301306072a8648ce3d020106082a811ccf5501822d04793077020101042051d03a5dcc7262900d2ac8ad4ea0511bcb8c6f62444ab0afd438a2a2eebacc1ea00a06082a811ccf5501822da14403420004ad3dc92f8c65bbe3d22e84170321fcb63d79e402e84d817126ec998ecb29cc1e8fb56a49d8b0e7f63f4d3cd00b9636da00c9a044d6ad5d1d7c684c4f6edd6e54`
	secp256k1String = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
)

var (
	rawSecp256k1Bytes, _ = hex.DecodeString(secp256k1String)
	rawPrivKey           = rawSecp256k1Bytes
)
var (
	tx *Tx
)

func TestNewTx(t *testing.T) {
	var err error
	tx, err = NewTx("test", common.ZeroAddress, []byte("Hello World"), getPrivKey())
	require.NoError(t, err)
}

func TestVerify(t *testing.T) {
	if !tx.Verify() {
		t.Fatal()
	}
	// change the nonce of tx
	nonce := tx.Data.AccountNonce
	tx.Data.AccountNonce = nonce + 1
	if tx.Verify() {
		t.Fatal()
	}
	tx.Data.AccountNonce = nonce
	if !tx.Verify() {
		t.Fatal()
	}
	// However, the situation is more complicated than what you thought
	sig := TxSig{
		PK:  tx.Data.Sig.PK,
		Sig: tx.Data.Sig.Sig,
	}
	// 1. set the pk to nil
	tx.Data.Sig.PK = nil
	if tx.Verify() {
		t.Fatal()
	}
	// 2. set the pk to random bytes
	tx.Data.Sig.PK = []byte("Fake pk")
	if tx.Verify() {
		t.Fatal()
	}
	// 3. set the sig to nil
	tx.Data.Sig.PK = sig.PK
	tx.Data.Sig.Sig = nil
	if tx.Verify() {
		t.Fatal()
	}
	// 4. set the sig to random bytes
	tx.Data.Sig.Sig = []byte("Fake sig")
	if tx.Verify() {
		t.Fatal()
	}

	// then make everything to be right
	tx.Data.Sig = &sig
	if !tx.Verify() {
		t.Fatal()
	}
}

func TestGetSender(t *testing.T) {
	sender, err := tx.GetSender()
	require.NoError(t, err)
	require.Equal(t, sender.String(), "0x970e8128ab834e8eac17ab8e3812f010678cf791")
}

func TestGetReceiver(t *testing.T) {
	receiver := tx.GetReceiver()
	if !reflect.DeepEqual(common.ZeroAddress.Bytes(), receiver.Bytes()) {
		t.Fatal()
	}
	privKey := getPrivKey()
	addr, err := privKey.PubKey().Address()
	require.NoError(t, err)

	selfTx, err := NewTx("test", addr, []byte("Hello World"), getPrivKey())
	require.NoError(t, err)

	if !reflect.DeepEqual(selfTx.GetReceiver().Bytes(), addr.Bytes()) {
		t.Fatal()
	}
}

func TestMarshaAndUnmarshalWithSig(t *testing.T) {
	txBytes, err := tx.Bytes()
	require.NoError(t, err)

	newTx, err := BytesToTx(txBytes)
	require.NoError(t, err)

	if !reflect.DeepEqual(tx, newTx) {
		t.Fatal()
	}
}

func TestMarshaAndUnmarshalWithoutSig(t *testing.T) {
	var tx = &Tx{
		Data: TxData{
			ChannelID:    "test",
			AccountNonce: 1,
			Recipient:    common.ZeroAddress.Bytes(),
			Payload:      []byte("Hello World"),
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	}
	txBytes, err := tx.Bytes()
	require.NoError(t, err)

	newTx, err := BytesToTx(txBytes)
	require.NoError(t, err)

	if !reflect.DeepEqual(tx, newTx) {
		t.Fatal()
	}
}

func getPrivKey() crypto.PrivateKey {
	privKey, _ := crypto.NewPrivateKey(rawPrivKey)
	return privKey
}
