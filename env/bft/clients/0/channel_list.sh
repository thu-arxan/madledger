#!/bin/bash
# This is our first script.
echo 'Test bft consensus...'
for ((i=1; i<=4; i++))
do
  echo 'try '$i' times to test channel list'
  client channel list
  #client channel list
  #echo $i
done

