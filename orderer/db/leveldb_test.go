package db

import (
	"errors"
	"fmt"
	cc "madledger/blockchain/config"
	"madledger/util"
	"os"
	"testing"
)

var (
	dir = ".leveldb"
	db  DB
)

func TestInit(t *testing.T) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewLevelDB(t *testing.T) {
	var err error
	db, err = NewLevelDB(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListChannel(t *testing.T) {
	channels := db.ListChannel()
	if len(channels) != 0 {
		t.Fatal()
	}
}

func TestUpdateChannel(t *testing.T) {
	err := db.UpdateChannel("_config", cc.Profile{
		Public: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	var channels []string
	channels = db.ListChannel()
	if len(channels) != 1 {
		t.Fatal(fmt.Errorf("Should contain one channel rather than %d channels", len(channels)))
	}
	if channels[0] != "_config" {
		t.Fatal(fmt.Errorf("Should contain channel _config rather than %s", channels[0]))
	}
	// add _global
	err = db.UpdateChannel("_global", cc.Profile{
		Public: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	channels = db.ListChannel()
	if len(channels) != 2 {
		t.Fatal(fmt.Errorf("Should contain two channels rather than %d channels", len(channels)))
	}
	if !util.Contain(channels, "_global") {
		t.Fatal(errors.New("Channel _global is not contained"))
	}
	// add user channel
	err = db.UpdateChannel("test", cc.Profile{
		Public: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	channels = db.ListChannel()
	if len(channels) != 3 {
		t.Fatal(fmt.Errorf("Should contain three channels rather than %d channels", len(channels)))
	}
	if !util.Contain(channels, "test") {
		t.Fatal(errors.New("Channel test is not contained"))
	}
	// todo: maybe illegal channel id
}
