#!/bin/bash
#run init.sh
gnome-terminal -e 'bash -c ". init.sh"'

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



