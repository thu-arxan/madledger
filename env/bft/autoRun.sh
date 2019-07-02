#!/bin/bash
#run init.sh
#pwd
. init.sh
sleep 2
#run orderers
echo 'run orderers'
cd ./orderers/0/
#pwd
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
sleep 1
cd ../1/
#pwd
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
sleep 1
cd ../2/
#pwd
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
sleep 1
cd ../3/
#pwd
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
sleep 1
#run peers
echo 'run peers'
cd ../../peers/0/
#pwd
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
sleep 1
cd ../1/
#pwd
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
sleep 1
cd ../2/
#pwd
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
sleep 1
cd ../3/
#pwd
gnome-terminal -e 'bash -c ". start.sh; exec bash"'
sleep 1
#list
echo 'list channels'
cd ../../clients/0/
gnome-terminal -e 'bash -c "client channel list; exec bash"'


