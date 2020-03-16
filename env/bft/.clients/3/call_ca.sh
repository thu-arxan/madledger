#!/bin/bash
# Copyright (c) 2020 THU-Arxan
# Madledger is licensed under Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#          http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
# EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
# MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
# See the Mulan PSL v2 for more details.

# This is our first script.
echo 'Test bft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test3 -f getNum -r 0x59a42b16a01bc674566dc51c821384b361317576
echo 'set number = 12 ...'
client tx call -a ./MyTest.abi -n test3 -f setNum -i 12 -r 0x59a42b16a01bc674566dc51c821384b361317576
echo 'add number 9 ...'
client tx call -a ./MyTest.abi -n test3 -f add -i 9 -r 0x59a42b16a01bc674566dc51c821384b361317576
echo 'subtract number 5 ...'
client tx call -a ./MyTest.abi -n test3 -f sub -i 5 -r 0x59a42b16a01bc674566dc51c821384b361317576
