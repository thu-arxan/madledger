// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package core

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

func TestNewTx(t *testing.T) {
	tx, err := NewTx("test", common.ZeroAddress, []byte("Hello World"), 0, "", getPrivKey())
	require.NoError(t, err)
	require.EqualValues(t, 0, tx.Data.Value)
	var algos = []crypto.Algorithm{crypto.KeyAlgoSM2, crypto.KeyAlgoSecp256k1}
	for i := range algos {
		privKey, err := crypto.GeneratePrivateKey(algos[i])
		require.NoError(t, err)
		tx, err = NewTx("test", common.ZeroAddress, []byte("Hello World"), 0, "", privKey)
		require.NoError(t, err)
		require.True(t, tx.Verify())
	}
}

func TestVerify(t *testing.T) {
	tx, err := NewTx("test", common.ZeroAddress, []byte("Hello World"), 0, "", getPrivKey())
	require.NoError(t, err)
	if !tx.Verify() {
		t.Fatal()
	}
	// change the nonce of tx
	nonce := tx.Data.Nonce
	tx.Data.Nonce = nonce + 1
	if tx.Verify() {
		t.Fatal()
	}
	tx.Data.Nonce = nonce
	if !tx.Verify() {
		t.Fatal()
	}
	// However, the situation is more complicated than what you thought
	sig := TxSig{
		PK:   tx.Data.Sig.PK,
		Sig:  tx.Data.Sig.Sig,
		Algo: tx.Data.Sig.Algo,
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
	tx.Data.Sig = sig
	require.True(t, tx.Verify())
}

func TestGetSender(t *testing.T) {
	tx, err := NewTx("test", common.ZeroAddress, []byte("Hello World"), 0, "", getPrivKey())
	require.NoError(t, err)
	sender, err := tx.GetSender()
	require.NoError(t, err)
	require.Equal(t, sender.String(), "0x970e8128ab834e8eac17ab8e3812f010678cf791")
}

func TestGetReceiver(t *testing.T) {
	tx, err := NewTx("test", common.ZeroAddress, []byte("Hello World"), 0, "", getPrivKey())
	require.NoError(t, err)
	receiver := tx.GetReceiver()
	if !reflect.DeepEqual(common.ZeroAddress.Bytes(), receiver.Bytes()) {
		t.Fatal()
	}
	privKey := getPrivKey()
	addr, err := privKey.PubKey().Address()
	require.NoError(t, err)

	selfTx, err := NewTx("test", addr, []byte("Hello World"), 0, "", getPrivKey())
	require.NoError(t, err)

	if !reflect.DeepEqual(selfTx.GetReceiver().Bytes(), addr.Bytes()) {
		t.Fatal()
	}
}

func TestMarshaAndUnmarshalWithSig(t *testing.T) {
	tx, err := NewTx("test", common.ZeroAddress, []byte("Hello World"), 0, "", getPrivKey())
	require.NoError(t, err)
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
			ChannelID: "test",
			Nonce:     1,
			Recipient: common.ZeroAddress.Bytes(),
			Payload:   []byte("Hello World"),
			Version:   1,
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

func BenchmarkNewTx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTx("test", common.ZeroAddress, []byte("Hello World"), 10, "To be or not to be, this is a question", getPrivKey())
	}
}

func BenchmarkMarshal(b *testing.B) {
	tx, _ := NewTx("test", common.ZeroAddress, []byte("Hello World"), 10, "To be or not to be, this is a question", getPrivKey())
	for i := 0; i < b.N; i++ {
		tx.Bytes()
	}
}

func BenchmarkVerify(b *testing.B) {
	tx, _ := NewTx("test", common.ZeroAddress, []byte("Hello World"), 10, "To be or not to be, this is a question", getPrivKey())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tx.Verify()
		}
	})
}

func getPrivKey() crypto.PrivateKey {
	privKey, _ := crypto.NewPrivateKey(rawPrivKey, crypto.KeyAlgoSecp256k1)
	return privKey
}
