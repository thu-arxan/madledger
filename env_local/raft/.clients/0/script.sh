# madledger/samples/clients/0

set -v

ADDRESS=0x3faab990ebfd0183c1402e6ec201ed0ddf0a5d81 # address show in `client account list`
CHANNEL=myccc

client account list
client asset issue -a $ADDRESS -v 10001
client channel create -n $CHANNEL -g 0 -m 1000000000
client tx create -b ./Balance.bin -n $CHANNEL
client tx call -a ./Balance.abi -n myccc -f add -i 5 -r 0xb0ff6840b5f9e0f7abd5beb5836ea296b1d0b790