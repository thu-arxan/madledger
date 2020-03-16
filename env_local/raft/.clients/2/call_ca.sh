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
client tx call -a ./MyTest.abi -n test2 -f getNum -r 0xfd2b9589386ddeb55c85c5a27f8d039c3b585dba
echo 'set number = 6 ...'
client tx call -a ./MyTest.abi -n test2 -f setNum -i 6 -r 0xfd2b9589386ddeb55c85c5a27f8d039c3b585dba
echo 'add number 8 ...'
client tx call -a ./MyTest.abi -n test2 -f add -i 8 -r 0xfd2b9589386ddeb55c85c5a27f8d039c3b585dba
echo 'subtract number 3 ...'
client tx call -a ./MyTest.abi -n test2 -f sub -i 3 -r 0xfd2b9589386ddeb55c85c5a27f8d039c3b585dba
