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
echo 'Test raft consensus by call contract...'
echo 'get number ...'
client tx call -a ./MyTest.abi -n test3 -f getNum -r 0x5f5490e4845d0f44fa0f71c3e6c55d89f4aa8222
echo 'set number = 12 ...'
client tx call -a ./MyTest.abi -n test3 -f setNum -i 12 -r 0x5f5490e4845d0f44fa0f71c3e6c55d89f4aa8222
echo 'add number 9 ...'
client tx call -a ./MyTest.abi -n test3 -f add -i 9 -r 0x5f5490e4845d0f44fa0f71c3e6c55d89f4aa8222
echo 'subtract number 5 ...'
client tx call -a ./MyTest.abi -n test3 -f sub -i 5 -r 0x5f5490e4845d0f44fa0f71c3e6c55d89f4aa8222
