#!/bin/bash
# This is our first script.
echo 'Test raft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test3 -f getNum -r 0x5f5490e4845d0f44fa0f71c3e6c55d89f4aa8222
echo 'set number = 12 ...'
client tx call -a ./MyTest.abi -n test3 -f setNum -i 12 -r 0x5f5490e4845d0f44fa0f71c3e6c55d89f4aa8222
echo 'add number 9 ...'
client tx call -a ./MyTest.abi -n test3 -f add -i 9 -r 0x5f5490e4845d0f44fa0f71c3e6c55d89f4aa8222
echo 'subtract number 5 ...'
client tx call -a ./MyTest.abi -n test3 -f sub -i 5 -r 0x5f5490e4845d0f44fa0f71c3e6c55d89f4aa8222
