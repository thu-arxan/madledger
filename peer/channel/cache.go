package channel

import (
	"encoding/binary"
	"errors"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core"
	"madledger/peer/db"
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

// GetGasParams returns channel specified gas info
// todo: ab
// could be empty?
func (cache *Cache) GetGasParams(channelID string) (uint64, uint64, error) {
	var maxGas, gasPrice uint64
	var maxGasBytes, gasPriceBytes []byte
	maxGasKey := util.BytesCombine([]byte(channelID), []byte("maxGas"))
	gasPriceKey := util.BytesCombine([]byte(channelID), []byte("gasPrice"))
	if _, ok := cache.kvs[string(maxGasKey)]; !ok {
		maxGasBytes, err := cache.Get(maxGasKey, true)
		if err != nil {
			return 0, 0, err
		}
		cache.kvs[string(maxGasKey)] = maxGasBytes
	}
	maxGasBytes = cache.kvs[string(maxGasKey)]
	maxGas = binary.BigEndian.Uint64(maxGasBytes)

	if _, ok := cache.kvs[string(gasPriceKey)]; !ok {
		gasPriceBytes, err := cache.Get(gasPriceKey, true)
		if err != nil {
			return 0, 0, err
		}
		cache.kvs[string(gasPriceKey)] = gasPriceBytes
	}
	gasPriceBytes = cache.kvs[string(gasPriceKey)]
	gasPrice = binary.BigEndian.Uint64(gasPriceBytes)
	return maxGas, gasPrice, nil
}

// GetToken return token sender has of channel
func (cache *Cache) GetToken(channelID string, sender common.Address) (uint64, error) {
	tokenKey := util.BytesCombine([]byte(channelID), []byte("token"), sender.Bytes())
	var tokenBytes []byte
	if _, ok := cache.kvs[string(tokenKey)]; !ok {
		tokenBytes, err := cache.Get(tokenKey, true)
		if err != nil {
			return 0, err
		}
		if tokenBytes == nil {
			tokenBytes = make([]byte, 8)
			binary.BigEndian.PutUint64(tokenBytes, 0)
		}
		cache.kvs[string(tokenKey)] = tokenBytes
	}
	tokenBytes = cache.kvs[string(tokenKey)]
	return binary.BigEndian.Uint64(tokenBytes), nil
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
	cache.adminPK, _ = crypto.NewPublicKey(pkBytes, crypto.KeyAlgoSecp256k1)
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

// SetTxStatus store tx execution information to db
func (cache *Cache) SetTxStatus(tx *core.Tx, status *db.TxStatus) error {
	return cache.wb.SetTxStatus(tx, status)
}

func (cache *Cache) SetToken(channelID string, sender common.Address, token []byte) {
	tokenKey := util.BytesCombine([]byte(channelID), []byte("token"), sender.Bytes())
	cache.kvs[string(tokenKey)] = token
	cache.wb.Put(tokenKey, token)
}

// PutBlock only used by addAssetBlock
// todo: why this is different from orderer
func (cache *Cache) PutBlock(block *core.Block) error {
	return cache.wb.PutBlock(block)
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
