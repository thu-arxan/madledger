#!/bin/bash
set -e

# start orderer
cd orderers/0 && orderer start -c orderer.yaml > orderer.log 2>&1 &

# start peer
for((i = 0; i <= 2; i++))
do
    cd peers/"$i" && peer start -c peer.yaml > peer."$i".log 2>&1 &
done
