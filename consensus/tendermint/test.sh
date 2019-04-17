rm -rf .test
mkdir .test
mkdir .test/env
cp -r ../../env/bft/.orderers/ /home/liuyihua/gopath/src/madledger/consensus/tendermint/.test/env
mv .test/env/.orderers .test/env/orderers
go test -v madledger/consensus/tendermint -count=1