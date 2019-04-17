package lib

import (
	"fmt"
	"madledger/client/config"
	"testing"
)

func  TestGetPeerClients(t *testing.T) {
	cfg, err := config.LoadConfig("/home/hadoop/GOPATH/src/madledger/client/config/.config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	peerClients, err := getPeerClients(cfg)
	if err != nil {
		t.Fatal(err)
	}
	for _, client := range peerClients {
		fmt.Println(client)
	}
}

func TestGetOPrdererClient(t *testing.T)  {
	cfg, err := config.LoadConfig("/home/hadoop/GOPATH/src/madledger/client/config/.config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ordererClient, err := getOrdererClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ordererClient)
}

func TestGetOrdererClients(t *testing.T){
	cfg, err := config.LoadConfig("/home/hadoop/GOPATH/src/madledger/client/config/.config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ordererClients, err := getOrdererClients(cfg)
	if err != nil {
		t.Fatal(err)
	}
	for _, client := range ordererClients {
		fmt.Println(client)
	}
}
