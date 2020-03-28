// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

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
		bs := privKey.Bytes()
		newPrivKey, err := NewPrivateKey(bs, algos[i])
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
		_, err = NewPublicKey(bs, algos[i])
		require.NoError(t, err)
	}
}

func TestLoadPrivateKeyFromFile(t *testing.T) {
	var algos = []Algorithm{KeyAlgoSM2, KeyAlgoSecp256k1}
	var keyPath = "priv.key"
	for i := range algos {
		privKey, err := GeneratePrivateKey(algos[i])
		require.NoError(t, err)
		bs := privKey.Bytes()
		if algos[i] == KeyAlgoSecp256k1 {
			ioutil.WriteFile(keyPath, []byte(fmt.Sprintf("%x", bs)), os.ModePerm)
		} else {
			ioutil.WriteFile(keyPath, bs, os.ModePerm)
		}
		newPrivKey, err := LoadPrivateKeyFromFile(keyPath)
		require.NoError(t, err)
		require.Equal(t, algos[i], newPrivKey.Algo())
	}
	os.Remove(keyPath)
}
