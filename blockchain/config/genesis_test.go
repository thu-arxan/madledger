package config

import (
	"testing"
)

func TestCreateGenesisBlock(t *testing.T) {
	admins, err := CreateAdmins()
	if err != nil {
		t.Fatal(err)
	}
	_, err = CreateGenesisBlock(admins)
	if err != nil {
		t.Fatal(err)
	}
}
