package crypto

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratePrivateKey(t *testing.T) {
	var algos = []Algorithm{KeyAlgoSM2, KeyAlgoSecp256k1}
	for i := range algos {
		privKey, err := GeneratePrivateKey(algos[i])
		require.NoError(t, err)
		require.Equal(t, algos[i], privKey.Algo())
	}
	// default sm2
	privKey, err := GeneratePrivateKey()
	require.NoError(t, err)
	require.Equal(t, KeyAlgoSM2, privKey.Algo())
}

func TestNewPrivateKey(t *testing.T) {
	var algos = []Algorithm{KeyAlgoSM2, KeyAlgoSecp256k1}
	for i := range algos {
		privKey, err := GeneratePrivateKey(algos[i])
		require.NoError(t, err)
		bs, err := privKey.Bytes()
		require.NoError(t, err)
		newPrivKey, err := NewPrivateKey(bs)
		require.NoError(t, err)
		require.Equal(t, algos[i], newPrivKey.Algo())
	}
}

func TestNewPublicKey(t *testing.T) {
	var algos = []Algorithm{KeyAlgoSM2, KeyAlgoSecp256k1}
	for i := range algos {
		privKey, err := GeneratePrivateKey(algos[i])
		require.NoError(t, err)
		bs, err := privKey.PubKey().Bytes()
		require.NoError(t, err)
		_, err = NewPublicKey(bs)
		require.NoError(t, err)
	}
}

func TestLoadPrivateKeyFromFile(t *testing.T) {
	var algos = []Algorithm{KeyAlgoSM2, KeyAlgoSecp256k1}
	var keyPath = "priv.key"
	for i := range algos {
		privKey, err := GeneratePrivateKey(algos[i])
		require.NoError(t, err)
		bs, err := privKey.Bytes()
		require.NoError(t, err)
		if algos[i] == KeyAlgoSecp256k1 {
			ioutil.WriteFile(keyPath, []byte(fmt.Sprintf("%x", bs)), os.ModePerm)
		} else {
			ioutil.WriteFile(keyPath, bs, os.ModePerm)
		}
		newPrivKey, err := LoadPrivateKeyFromFile(keyPath)
		require.NoError(t, err)
		require.Equal(t, algos[i], newPrivKey.Algo())
	}
}
