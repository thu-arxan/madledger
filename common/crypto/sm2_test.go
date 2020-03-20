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
	"madledger/common/crypto/hash"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSM2(t *testing.T) {
	priv, err := GenerateSM2PrivateKey()
	require.NoError(t, err)
	pub := priv.PubKey()
	digest := hash.SM3([]byte("hello world"))
	sig, err := priv.Sign(digest)
	require.True(t, sig.Verify(digest, pub))
}

func TestSM2Priv(t *testing.T) {
	priv, err := GenerateSM2PrivateKey()
	require.NoError(t, err)
	pubKeyBytes, err := priv.PubKey().Bytes()
	require.NoError(t, err)
	bs, err := priv.Bytes()
	require.NoError(t, err)
	newPriv, err := toSM2PrivateKey(bs)
	require.NoError(t, err)
	newPubKeyBytes, err := newPriv.PubKey().Bytes()
	require.NoError(t, err)
	require.Equal(t, pubKeyBytes, newPubKeyBytes)
}

func TestSM2Pub(t *testing.T) {
	priv, err := GenerateSM2PrivateKey()
	require.NoError(t, err)
	pubKeyBytes, err := priv.PubKey().Bytes()
	require.NoError(t, err)
	digest := hash.SM3([]byte("hello world"))
	sig, err := priv.Sign(digest)
	pub, err := newSM2PublicKey(pubKeyBytes)
	require.NoError(t, err)
	require.True(t, sig.Verify(digest, pub))
}
