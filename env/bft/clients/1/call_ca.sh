#!/bin/bash
# This is our first script.
echo 'Test bft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test1 -f getNum -r 0x1b66001e01d3c8d3893187fee59e3bea1d9bdd9b
echo 'set number = 5 ...'
client tx call -a ./MyTest.abi -n test1 -f setNum -i 5 -r 0x1b66001e01d3c8d3893187fee59e3bea1d9bdd9b
echo 'add number 7 ...'
client tx call -a ./MyTest.abi -n test1 -f add -i 7 -r 0x1b66001e01d3c8d3893187fee59e3bea1d9bdd9b
echo 'subtract number 4 ...'
client tx call -a ./MyTest.abi -n test1 -f sub -i 4 -r 0x1b66001e01d3c8d3893187fee59e3bea1d9bdd9b
