package types

import (
	"encoding/hex"
	"fmt"
	"madledger/common"
	"madledger/common/crypto"
	"testing"
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
	tx, err = NewTx("test", common.ZeroAddress, []byte(""), getPrivKey())
	if err != nil {
		t.Fatal(err)
	}

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
}

func TestGetSender(t *testing.T) {
	sender, err := tx.GetSender()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sender.Bytes())
	t.Fatal()
}

func getPrivKey() crypto.PrivateKey {
	privKey, _ := crypto.NewPrivateKey(rawPrivKey)
	return privKey
}
