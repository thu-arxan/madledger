#!/bin/bash
# This is our first script.
echo 'Test bft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test2 -f getNum -r 0x6a23b446720edeeeea79035b3919d59442fa67da
echo 'set number = 6 ...'
client tx call -a ./MyTest.abi -n test2 -f setNum -i 6 -r 0x6a23b446720edeeeea79035b3919d59442fa67da
echo 'add number 8 ...'
client tx call -a ./MyTest.abi -n test2 -f add -i 8 -r 0x6a23b446720edeeeea79035b3919d59442fa67da
echo 'subtract number 3 ...'
client tx call -a ./MyTest.abi -n test2 -f sub -i 3 -r 0x6a23b446720edeeeea79035b3919d59442fa67da
