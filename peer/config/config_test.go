package config

import (
	"fmt"
	"madledger/common/util"
	"os"
	"reflect"
	"testing"
)

func TestGetServerConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	if err != nil {
		t.Fatal(err)
	}
	serverCfg, err := cfg.GetServerConfig()
	if err != nil {
		t.Fatal(err)
	}
	if serverCfg.Port != 23456 {
		t.Fatal(fmt.Errorf("The port is %d", serverCfg.Port))
	}
	if serverCfg.Address != "localhost" {
		t.Fatal(fmt.Errorf("The address is %s", serverCfg.Address))
	}
	if serverCfg.Debug != true {
		t.Fatal(fmt.Errorf("The Debug is %t", serverCfg.Debug))
	}
	// then change the value of cfg
	// check address
	cfg.Address = ""
	_, err = cfg.GetServerConfig()
	if err.Error() != "The address can not be empty" {
		t.Fatal(err)
	}
	// check port
	cfg.Address = "localhost"
	cfg.Port = -1
	_, err = cfg.GetServerConfig()
	if err.Error() != "The port can not be -1" {
		t.Fatal(err)
	}
	cfg.Port = -1
	_, err = cfg.GetServerConfig()
	if err == nil || err.Error() != "The port can not be -1" {
		t.Fatal(err)
	}
	cfg.Port = 1023
	_, err = cfg.GetServerConfig()
	if err == nil || err.Error() != "The port can not be 1023" {
		t.Fatal(err)
	}
	cfg.Port = 1024
	_, err = cfg.GetServerConfig()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOrdererConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	if err != nil {
		t.Fatal(err)
	}
	ordererCfg, err := cfg.GetOrdererConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ordererCfg.Address, []string{"localhost:12345"}) {
		t.Fatal(fmt.Errorf("Address is %s", ordererCfg.Address))
	}
}

func getTestConfigFilePath() string {
	gopath := os.Getenv("GOPATH")
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/peer/config/.config.yaml", gopath)
	return cfgFilePath
}
