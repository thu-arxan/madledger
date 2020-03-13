#!/bin/bash
# This is our first script.
echo 'Test raft consensus by creating channels'
client asset issue -v 1000000000000 -s
for ((i=0; i<=80; i+=10))
do
  echo 'create channel test'$i
  client channel create -n test$i
  client asset transfer -v 100000000 -n test$i
  client channel token -v 1000000000 -n test$i
done

client channel list
