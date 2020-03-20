package sm2

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPemWriteLoad(t *testing.T) {
	priv, err := GenerateKey()
	if err != nil {
		t.Error(err)
	}
	privPEM, err := WritePrivateKeytoMem(priv, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(privPEM))

	ret, err := WritePrivateKeyToPem("priv.pem", priv, nil)
	if err != nil || !ret {
		t.Error(err)
	}

	content, err := ioutil.ReadFile("priv.pem")
	if err != nil {
		t.Error(err)
	}
	if string(content) != string(privPEM) {
		t.Error("not match")
	}

	pubPem, err := WritePublicKeytoMem(priv.PublicKey, nil)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(pubPem))

	ret, err = WritePublicKeyToPem("pub.pem", priv.PublicKey, nil)
	if err != nil || !ret {
		t.Error(err)
	}

	content, err = ioutil.ReadFile("pub.pem")
	if err != nil {
		t.Error(err)
	}
	if string(content) != string(pubPem) {
		t.Error("not match")
	}

	priv2, err := ReadPrivateKeyFromPem("priv.pem", nil)
	if err != nil {
		t.Error(err)
	}
	privPEM2, err := WritePrivateKeytoMem(priv2, nil)
	if err != nil {
		t.Error(err)
	}
	if hex.EncodeToString(privPEM2) != hex.EncodeToString(privPEM) {
		t.Error("not match")
	}

	pubPEM2, err := WritePublicKeytoMem(priv2.PublicKey, nil)
	if err != nil {
		t.Error(err)
	}
	if hex.EncodeToString(pubPem) != hex.EncodeToString(pubPEM2) {
		t.Error("Not Match")
	}

	pub3, err := ReadPublicKeyFromPem("pub.pem", nil)
	if err != nil {
		t.Error(err)
	}
	pubPem3, err := WritePublicKeytoMem(pub3, nil)
	if err != nil {
		t.Error(err)
	}
	if hex.EncodeToString(pubPem) != hex.EncodeToString(pubPem3) {
		t.Error("Not Match")
	}
}
