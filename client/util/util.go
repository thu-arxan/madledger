package util

import (
	"io/ioutil"
	"madledger/common/crypto"
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
	hash := cutil.Hex(crypto.Hash(privKeyBytes))
	filePath, err := cutil.MakeFileAbs(hash, path)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filePath, []byte(privKeyHex), 0600)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// LoadPrivateKey load private key from file
func LoadPrivateKey(file string) (crypto.PrivateKey, error) {
	return crypto.LoadPrivateKeyFromFile(file)
}
