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

set -e

# start orderer
cd orderers/0 && orderer start -c orderer.yaml > orderer.log 2>&1 &

# start peer
for((i = 0; i <= 2; i++))
do
    cd peers/"$i" && peer start -c peer.yaml > peer."$i".log 2>&1 &
done
