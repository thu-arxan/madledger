#!/bin/bash
# This is our first script.
echo 'Test solo consensus by call contract...'
echo 'get balance ...'
client tx call -a ./Balance/Balance.abi -n test0 -f get -r 0x8b3f0e6422f392defd6a1db282f0bbd778f3ff56
echo 'set balance = 4 ...'
client tx call -a ./Balance/Balance.abi -n test0 -f set -i 4 -r 0x8b3f0e6422f392defd6a1db282f0bbd778f3ff56
echo 'balance add 6 ...'
client tx call -a ./Balance/Balance.abi -n test0 -f add -i 6 -r 0x8b3f0e6422f392defd6a1db282f0bbd778f3ff56
echo 'balance subtract 2 ...'
client tx call -a ./Balance/Balance.abi -n test0 -f sub -i 2 -r 0x8b3f0e6422f392defd6a1db282f0bbd778f3ff56
echo 'get info ...'
client tx call -a ./Balance/Balance.abi -n test0 -f info -r 0x8b3f0e6422f392defd6a1db282f0bbd778f3ff56
