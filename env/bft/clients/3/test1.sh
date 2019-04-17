#!/bin/bash
# This is our first script.
echo 'Test bft consensus...'
for ((i=17; i<=87; i+=10))
do
  echo 'create channel test'$i
  client channel create -n test$i
  #client channel list
  #echo $i
done

