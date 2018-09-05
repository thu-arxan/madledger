package crypto

import (
	"encoding/hex"
	"fmt"
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

func TestNewPrivateKey(t *testing.T) {
	_, err := NewPrivateKey(rawPrivKey)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignVerify(t *testing.T) {
	privKey, _ := NewPrivateKey(rawPrivKey)
	hash := Hash([]byte("abc"))
	sig, err := privKey.Sign(hash)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := privKey.PubKey()
	if !sig.Verify(hash, pubKey) {
		t.Fatal()
	}
	if sig.Verify(Hash([]byte("ab")), pubKey) {
		t.Fatal()
	}
}

func TestAddress(t *testing.T) {
	privKey, _ := NewPrivateKey(rawPrivKey)
	pubKey := privKey.PubKey()
	addr, err := pubKey.Address()
	if err != nil {
		t.Fatal(err)
	}
	if addr.String() != "0x970e8128ab834e8eac17ab8e3812f010678cf791" {
		t.Fatal(fmt.Errorf("The address is %s", addr.String()))
	}
}

func BenchmarkNewPrivateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPrivateKey(rawPrivKey)
	}
}

func BenchmarkSign(b *testing.B) {
	privKey, _ := NewPrivateKey(rawPrivKey)
	hash := Hash([]byte("abc"))
	for i := 0; i < b.N; i++ {
		privKey.Sign(hash)
	}
}

func BenchmarkGetPubKeyFromPrivKey(b *testing.B) {
	privKey, _ := NewPrivateKey(rawPrivKey)
	for i := 0; i < b.N; i++ {
		privKey.PubKey()
	}
}

func BenchmarkSignVerify(b *testing.B) {
	privKey, _ := NewPrivateKey(rawPrivKey)
	hash := Hash([]byte("abc"))
	sig, _ := privKey.Sign(hash)
	pubKey := privKey.PubKey()
	for i := 0; i < b.N; i++ {
		sig.Verify(hash, pubKey)
	}
}