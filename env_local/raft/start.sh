#!/bin/bash
set -e

for((i = 0; i <= 3; i++))
do
    cd orderers/"$i" && orderer start -c orderer.yaml > orderer."$i".log 2>&1 &
    cd peers/"$i" && peer start -c peer.yaml > peer."$i".log 2>&1 &
done
