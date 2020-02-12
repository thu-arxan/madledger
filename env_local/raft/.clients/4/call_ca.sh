#!/bin/bash
# This is our first script.
echo 'Test raft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test4 -f getNum -r 0xf6d5f19750a5976b01247cd65744809f252672d2
echo 'set number = 12 ...'
client tx call -a ./MyTest.abi -n test4 -f setNum -i 12 -r 0xf6d5f19750a5976b01247cd65744809f252672d2
echo 'add number 9 ...'
client tx call -a ./MyTest.abi -n test4 -f add -i 9 -r 0xf6d5f19750a5976b01247cd65744809f252672d2
echo 'subtract number 5 ...'
client tx call -a ./MyTest.abi -n test4 -f sub -i 5 -r 0xf6d5f19750a5976b01247cd65744809f252672d2
