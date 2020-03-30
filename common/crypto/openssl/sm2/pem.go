// +build cgo

package sm2

/*
#include "util.h"
#include <stdlib.h>
#include <openssl/evp.h>
#include <openssl/ec.h>
#include <openssl/x509.h>
*/
import "C"
import (
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

// ReadPrivateKeyFromMem parse private key from pem data, ignore pwd
func ReadPrivateKeyFromMem(data []byte, pwd []byte) (*PrivateKey, error) {

	if len(data) == 0 {
		return nil, errors.New("empty pem block")
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode priv key")
	}
	return ParseSm2PrivateKey(block.Bytes)
}

// ReadPrivateKeyFromPem parse private key from pem file, ignore pwd
func ReadPrivateKeyFromPem(FileName string, pwd []byte) (*PrivateKey, error) {
	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	return ReadPrivateKeyFromMem(data, pwd)
}

// WritePrivateKeytoMem encoding private key into pem data, ignore pwd
func WritePrivateKeytoMem(key *PrivateKey, pwd []byte) ([]byte, error) {
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: key.Key,
	}
	return pem.EncodeToMemory(block), nil
}

// WritePrivateKeyToPem encoding private key into pem data, ignore pwd
func WritePrivateKeyToPem(FileName string, key *PrivateKey, pwd []byte) (bool, error) {
	if key == nil || key.Key == nil {
		return false, errors.New("Nil private key")
	}
	file, err := os.Create(FileName)
	defer file.Close()
	if err != nil {
		return false, err
	}
	data, err := WritePrivateKeytoMem(key, pwd)
	if err != nil {
		return false, err
	}
	_, err = file.Write(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ReadPublicKeyFromMem load public key from der-encoded block
func ReadPublicKeyFromMem(data []byte, _ []byte) (*PublicKey, error) {
	if len(data) == 0 {
		return nil, errors.New("empty pem block")
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode public key")
	}
	return ParseSm2PublicKey(block.Bytes)
}

// ReadPublicKeyFromPem load public key from pem
func ReadPublicKeyFromPem(FileName string, pwd []byte) (*PublicKey, error) {
	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	return ReadPublicKeyFromMem(data, pwd)
}

// WritePublicKeytoMem encoding pk into pem block
func WritePublicKeytoMem(key *PublicKey, _ []byte) ([]byte, error) {
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: key.Key,
	}
	return pem.EncodeToMemory(block), nil
}

// WritePublicKeyToPem encoding pk into pem file
func WritePublicKeyToPem(FileName string, key *PublicKey, _ []byte) (bool, error) {
	if key == nil || key.Key == nil {
		return false, errors.New("Nil public key")
	}
	file, err := os.Create(FileName)
	defer file.Close()
	if err != nil {
		return false, err
	}
	data, err := WritePublicKeytoMem(key, nil)
	if err != nil {
		return false, err
	}
	_, err = file.Write(data)
	if err != nil {
		return false, err
	}
	if key == nil || key.Key == nil {
		return false, errors.New("Unexpected nil")
	}
	return true, nil
}
