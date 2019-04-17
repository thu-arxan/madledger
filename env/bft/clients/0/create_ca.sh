#!/bin/bash
# This is our first script.
echo 'Test bft consensus by creating contract...'
for ((i=1; i<=8; i++))
do
  echo 'create contract '$i
  client tx create -b /home/hadoop/solidity/MyTest.bin -n test0
  #client channel list
  #echo $i
done

