# delete all .git in vendor
find $GOPATH/src/madledger/vendor/ -name .git |xargs rm -rf