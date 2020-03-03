#!/bin/bash
# This is our first script.
echo 'Test raft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test1 -f getNum -r 0x2fb2b97af9950eddb511cd149ddadaed6273c9c0
echo 'set number = 5 ...'
client tx call -a ./MyTest.abi -n test1 -f setNum -i 5 -r 0x2fb2b97af9950eddb511cd149ddadaed6273c9c0
echo 'add number 7 ...'
client tx call -a ./MyTest.abi -n test1 -f add -i 7 -r 0x2fb2b97af9950eddb511cd149ddadaed6273c9c0
echo 'subtract number 4 ...'
client tx call -a ./MyTest.abi -n test1 -f sub -i 4 -r 0x2fb2b97af9950eddb511cd149ddadaed6273c9c0
