package config

import (
	"testing"
)

func TestCreateGenesisBlock(t *testing.T) {
	_, err := CreateGenesisBlock()
	if err != nil {
		t.Fatal(err)
	}
}
