# Copyright (c) 2020 THU-Arxan
# Madledger is licensed under Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#          http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
# EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
# MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
# See the Mulan PSL v2 for more details.

# test the order of blocks
rm -rf $GOPATH/src/madledger/peer/channel/log.txt
rm -rf $GOPATH/src/madledger/peer/channel/.data
go test -v madledger/peer/channel -count=1 -s 1 > $GOPATH/src/madledger/peer/channel/log.txt
go test madledger/peer/channel -s 2 -count=1 -cover