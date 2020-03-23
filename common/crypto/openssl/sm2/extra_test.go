package sm2

import (
	"fmt"
	"testing"
)

func TestIdPrivate(t *testing.T) {
	priv, err := GenerateKey()
	if err != nil {
		t.Error(err)
	}
	_, err = MarshalSm2UnecryptedPrivateKey(priv)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(base64Encode(p1))
	pem, err := WritePrivateKeytoMem(priv, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pem))
	_, err = MarshalSm2UnecryptedPrivateKey(priv)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(base64Encode(p2))
}
