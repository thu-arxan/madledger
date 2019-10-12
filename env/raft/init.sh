# This file will init orderers, clients and peers
rm -rf orderers
cp -R .orderers orderers
rm -rf clients
cp -R .clients clients
rm -rf peers
cp -R .peers peers
rm -rf autoRun.sh
cp .autoRun.sh autoRun.sh
