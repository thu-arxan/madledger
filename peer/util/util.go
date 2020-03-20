// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package util

import (
	"io/ioutil"
	"madledger/common/crypto"
	"madledger/common/crypto/hash"
	cutil "madledger/common/util"
)

// GeneratePrivateKey try to generate a private key below the path
func GeneratePrivateKey(path string) (string, error) {
	privKey, err := crypto.GeneratePrivateKey()
	if err != nil {
		return "", err
	}
	privKeyBytes, _ := privKey.Bytes()
	privKeyHex := cutil.Hex(privKeyBytes)
	var digest string
	switch privKey.Algo() {
	case crypto.KeyAlgoSecp256k1:
		digest = cutil.Hex(hash.SHA256(privKeyBytes))
	default:
		digest = cutil.Hex(hash.SM3(privKeyBytes))
	}
	filePath, err := cutil.MakeFileAbs(digest, path)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filePath, []byte(privKeyHex), 0600)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
