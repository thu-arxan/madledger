#!/bin/bash
# This is our first script.
echo 'Test bft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test5 -f getNum -r 0x917dc1ec18b5dfe46f311f462f26150aa261126d
echo 'set number = 4 ...'
client tx call -a ./MyTest.abi -n test5 -f setNum -i 4 -r 0x917dc1ec18b5dfe46f311f462f26150aa261126d
echo 'add number 6 ...'
client tx call -a ./MyTest.abi -n test5 -f add -i 6 -r 0x917dc1ec18b5dfe46f311f462f26150aa261126d
echo 'subtract number 2 ...'
client tx call -a ./MyTest.abi -n test5 -f sub -i 2 -r 0x917dc1ec18b5dfe46f311f462f26150aa261126d
