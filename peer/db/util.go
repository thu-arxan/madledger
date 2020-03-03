package db

import (
	"errors"
	"madledger/common"
	"madledger/common/util"
)

// MarshalAccount provide a fast marshal implementaion of marshal account
func MarshalAccount(account *common.Account) []byte {
	var bytes = make([]byte, 0)
	bytes = util.BytesCombine(bytes, account.GetAddress().Bytes())
	bytes = util.BytesCombine(bytes, util.Uint64ToBytes(account.GetBalance()))
	bytes = util.BytesCombine(bytes, util.Uint64ToBytes(account.GetNonce()))
	bytes = util.BytesCombine(bytes, util.BoolToBytes(account.HasSuicide()))
	if len(account.GetCode()) != 0 {
		bytes = util.BytesCombine(bytes, account.GetCode())
	}
	return bytes
}

// UnmarshalAccount provide a fast unmarshal implementation of unmarshal account
func UnmarshalAccount(bytes []byte) (*common.Account, error) {
	var account = new(common.Account)
	if len(bytes) < 37 {
		return nil, errors.New("wrong length")
	}
	addr, err := common.AddressFromBytes(bytes[:20])
	if err != nil {
		return nil, err
	}
	balance, err := util.BytesToUint64(bytes[20:28])
	if err != nil {
		return nil, err
	}
	nonce, err := util.BytesToUint64(bytes[28:36])
	if err != nil {
		return nil, err
	}
	var suicide = false
	if bytes[36] == 1 {
		suicide = true
	}
	if len(bytes) > 37 {
		account.Code = bytes[37:]
	}
	account.Address = addr
	account.Balance = balance
	account.Nonce = nonce
	account.SuicideMark = suicide
	return account, nil
}
