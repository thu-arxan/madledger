#!/bin/bash
# This is our first script.
echo 'Test raft consensus by creating channels'
for ((i=2; i<=82; i+=10))
do
  echo 'create channel test'$i
  client channel create -n test$i
  #echo $i
done

client channel list
