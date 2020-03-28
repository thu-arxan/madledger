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


#run orderers
echo 'run orderers'
cd ./orderers/0/
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
cd ../1/
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
cd ../2/
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
cd ../3/
gnome-terminal -e 'bash -c ". start.sh; exec bash"'

#run peers
echo 'run peers'
cd ../../peers/0/
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
cd ../1/
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
cd ../2/
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
cd ../3/
gnome-terminal -e 'bash -c ". start.sh; exec bash"'

#list
echo 'list channels'
cd ../../clients/0/
gnome-terminal -e 'bash -c "client channel list; exec bash"'



