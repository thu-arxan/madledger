package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"madledger/common/util"
	"madledger/core"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/orderer/db"
	"reflect"
)

type Cache struct {
	db db.DB
	wb db.WriteBatch
	storage map[common.Address]common.Account
	adminPK crypto.PublicKey
}

func NewCache(db db.DB) Cache {
	return Cache{
		db: db,
		wb: db.NewWriteBatch(),
		storage: make(map[common.Address]common.Account),
	}
}

// IsAssetAdmin decides whether a pk is the admin public key of _asset
func(cache *Cache) IsAssetAdmin(pk crypto.PublicKey) bool {
	if pk == nil {
		return false
	}
	if cache.adminPK != nil {
		return reflect.DeepEqual(pk, cache.adminPK)
	}
	pkBytes, err := cache.db.Get(getAssetAdminKey())
	if err != nil {
		return false
	}
	cache.adminPK, _ = crypto.NewPublicKey(pkBytes)
	return pk == cache.adminPK
}

func (cache *Cache) GetOrCreateAccount(address common.Address) (common.Account, error) {
	account := common.Account{}
	if account, ok := cache.storage[address]; ok {
		return account, nil
	}
	value, err := cache.db.GetIgnoreNotFound(getAccountKey(address))
	if err != nil {
		return account, err
	}
	err = json.Unmarshal(value, &account)
	if  err != nil {
		return account, err
	}
	cache.storage[address] = account
	return account, nil
}

// UpdateAccounts update account info
func(cache *Cache) UpdateAccounts(accounts ...common.Account) error {
	for _, acc := range accounts {
		cache.storage[acc.GetAddress()] = acc
	}
	return nil
}

// SetAssetAdmin only works when it is first called
func(cache *Cache) SetAssetAdmin(pk crypto.PublicKey) error {
	if cache.adminPK != nil {
		return errors.New("_asset admin exists")
	}
	pkBytes, err := cache.db.Get(getAssetAdminKey())
	if err != nil {
		return err
	}
	if pkBytes != nil {
		return errors.New("_asset admin exists")
	}
	cache.adminPK = pk
	cache.wb.Put(getAssetAdminKey(), pkBytes)
	return nil
}

//SetTxStatus store tx execution information to db
func(cache *Cache) SetTxStatus(tx *core.Tx, status *db.TxStatus) error {
	value, err := json.Marshal(status)
	if err != nil {
		return err
	}
	var key = util.BytesCombine([]byte(tx.Data.ChannelID), []byte(tx.ID))
	cache.wb.Put(key, value)
	return nil
}

// Sync writes updated data in cache to db
func(cache *Cache) Sync() error {
	for addr, acc := range cache.storage {
		key := getAccountKey(addr)
		value, err := json.Marshal(acc)
		if err != nil {
			return err
		}
		cache.wb.Put(key, value)
	}
	return cache.wb.Sync()
}

func getAccountKey(address common.Address) []byte {
	return []byte(fmt.Sprintf("%s@%s", core.ASSETCHANNELID, address.String()))
}

func getAssetAdminKey() []byte {
	return []byte("_asset_admin")
}
