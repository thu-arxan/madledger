package channel

import (
	"errors"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/orderer/db"
	"reflect"
)

// Cache used for AddAssetBlock
type Cache struct {
	db       db.DB
	wb       db.WriteBatch
	accounts map[common.Address]common.Account
	adminPK  crypto.PublicKey
	// useful kvs that get and set by []byte
	kvs map[string][]byte
}

// NewCache new a cache for AddAssetBlock
func NewCache(db db.DB) Cache {
	return Cache{
		db:       db,
		wb:       db.NewWriteBatch(),
		accounts: make(map[common.Address]common.Account),
		kvs:      make(map[string][]byte),
	}
}

// Get get []byte indexed by []byte from db
func (cache *Cache) Get(key []byte, couldBeEmpty bool) ([]byte, error) {
	if val, ok := cache.kvs[string(key)]; ok {
		return val, nil
	}
	val, err := cache.db.Get(key, couldBeEmpty)
	if err != nil {
		return nil, err
	}
	cache.kvs[string(key)] = val
	return val, nil
}

// IsAssetAdmin decides whether a pk is the admin public key of _asset
func (cache *Cache) IsAssetAdmin(pk crypto.PublicKey, pkAlgo crypto.Algorithm) bool {
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
	cache.adminPK, _ = crypto.NewPublicKey(pkBytes, pkAlgo)
	return reflect.DeepEqual(pk, cache.adminPK)
}

// GetOrCreateAccount returns default account if not exist
func (cache *Cache) GetOrCreateAccount(address common.Address) (common.Account, error) {
	if account, ok := cache.accounts[address]; ok {
		return account, nil
	}
	account, err := cache.db.GetOrCreateAccount(address)
	return account, err
}

// UpdateAccounts update account info
func (cache *Cache) UpdateAccounts(accs ...common.Account) error {
	for _, acc := range accs {
		cache.accounts[acc.GetAddress()] = acc
	}
	return cache.wb.UpdateAccounts(accs...)
}

// SetAssetAdmin only works when it is first called
func (cache *Cache) SetAssetAdmin(pk crypto.PublicKey, pkAlgo crypto.Algorithm) error {
	if cache.adminPK != nil {
		return errors.New("_asset admin exists")
	}
	pkBytes := cache.db.GetAssetAdminPKBytes()
	if pkBytes != nil {
		cache.adminPK, _ = crypto.NewPublicKey(pkBytes, pkAlgo)
		return errors.New("_asset admin exists")
	}
	cache.adminPK = pk
	return cache.wb.SetAssetAdmin(pk)
}

// Put store []byte indexed by []byte
func (cache *Cache) Put(key, value []byte) {
	cache.kvs[string(key)] = value
	cache.wb.Put(key, value)
}

// Sync writes updated data in cache to db
func (cache *Cache) Sync() error {
	return cache.wb.Sync()
}
