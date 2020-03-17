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
client tx call -a ./MyTest.abi -n test5 -f getNum -r 0x3a546aef4692784fdc92df8ca4a4a563fa67ef5e
echo 'set number = 4 ...'
client tx call -a ./MyTest.abi -n test5 -f setNum -i 4 -r 0x3a546aef4692784fdc92df8ca4a4a563fa67ef5e
echo 'add number 6 ...'
client tx call -a ./MyTest.abi -n test5 -f add -i 6 -r 0x3a546aef4692784fdc92df8ca4a4a563fa67ef5e
echo 'subtract number 2 ...'
client tx call -a ./MyTest.abi -n test5 -f sub -i 2 -r 0x3a546aef4692784fdc92df8ca4a4a563fa67ef5e
