#!/bin/bash
# This is our first script.
echo 'Test bft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test4 -f getNum -r 0xecfca247877013d1202f417f77daf8572cf5020c
echo 'set number = 12 ...'
client tx call -a ./MyTest.abi -n test4 -f setNum -i 12 -r 0xecfca247877013d1202f417f77daf8572cf5020c
echo 'add number 9 ...'
client tx call -a ./MyTest.abi -n test4 -f add -i 9 -r 0xecfca247877013d1202f417f77daf8572cf5020c
echo 'subtract number 5 ...'
client tx call -a ./MyTest.abi -n test4 -f sub -i 5 -r 0xecfca247877013d1202f417f77daf8572cf5020c
