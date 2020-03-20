package channel

import (
	"errors"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	"madledger/peer/db"
	"reflect"
)

type Cache struct {
	db      db.DB
	wb      db.WriteBatch
	storage map[common.Address]common.Account
	adminPK crypto.PublicKey
	token   map[string][]byte
}

func NewCache(db db.DB) Cache {
	return Cache{
		db:      db,
		wb:      db.NewWriteBatch(),
		storage: make(map[common.Address]common.Account),
		token:   make(map[string][]byte),
	}
}

func (cache *Cache) GetToken(key []byte) ([]byte, error) {
	if _, ok := cache.token[string(key)]; ok {
		return cache.token[string(key)], nil
	}
	return cache.db.Get(key)
}

func (cache *Cache) SetToken(key, value []byte) {
	cache.token[string(key)] = value
}

// IsAssetAdmin decides whether a pk is the admin public key of _asset
func (cache *Cache) IsAssetAdmin(pk crypto.PublicKey) bool {
	if pk == nil {
		return false
	}
	if cache.adminPK != nil {
		return reflect.DeepEqual(pk, cache.adminPK)
	}
	pkBytes := cache.db.GetAssetAdminPKBytes()
	if pkBytes == nil {
		return false
	}
	cache.adminPK, _ = crypto.NewPublicKey(pkBytes)
	return reflect.DeepEqual(pk, cache.adminPK)
}

func (cache *Cache) GetOrCreateAccount(address common.Address) (common.Account, error) {
	if account, ok := cache.storage[address]; ok {
		return account, nil
	}
	account, err := cache.db.GetOrCreateAccount(address)
	return account, err
}

// UpdateAccounts update account info
func (cache *Cache) UpdateAccounts(accounts ...common.Account) error {
	for _, acc := range accounts {
		cache.storage[acc.GetAddress()] = acc
	}
	return cache.wb.UpdateAccounts(accounts...)
}

// SetAssetAdmin only works when it is first called
func (cache *Cache) SetAssetAdmin(pk crypto.PublicKey) error {
	if cache.adminPK != nil {
		return errors.New("_asset admin exists")
	}
	pkBytes := cache.db.GetAssetAdminPKBytes()
	if pkBytes != nil {
		return errors.New("_asset admin exists")
	}
	cache.adminPK = pk
	return cache.wb.SetAssetAdmin(pk)
}

//SetTxStatus store tx execution information to db
func (cache *Cache) SetTxStatus(tx *core.Tx, status *db.TxStatus) error {
	return cache.wb.SetTxStatus(tx, status)
}

// Sync writes updated data in cache to db
func (cache *Cache) Sync() error {
	for key, value := range cache.token {
		cache.wb.Put([]byte(key), value)
	}
	return cache.wb.Sync()
}
