#!/bin/bash
set -e

for((i = 0; i <= 3; i++))
do
    cd orderers/"$i" && orderer start -c orderer.yaml > $GOPATH/src/madledger/samples/orderer."$i".log 2>&1 &
    cd peers/"$i" && peer start -c peer.yaml > $GOPATH/src/madledger/samples/peer."$i".log 2>&1 &
done
