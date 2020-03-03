#!/bin/bash
# This is our first script.
echo 'Test raft consensus by creating channels'
for ((i=1; i<=81; i+=10))
do
  echo 'create channel test'$i
  client channel create -n test$i
  #echo $i
done

client channel list
