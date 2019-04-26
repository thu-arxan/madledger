package db

import (
	"encoding/json"
	"errors"
	"madledger/common"
	"madledger/common/util"
	"madledger/core/types"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/syndtr/goleveldb/leveldb"
)

/*
* Here defines some key rules.
* 1. Account: key = []bytes("account:") + address.Bytes()
* 2. Storage: key = address.Bytes()
 */

// LevelDB is the implementation of DB on leveldb
type LevelDB struct {
	// the dir of data
	dir     string
	connect *leveldb.DB
}

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "db"})
)

// NewLevelDB is the constructor of LevelDB
func NewLevelDB(dir string) (DB, error) {
	db := new(LevelDB)
	db.dir = dir
	connect, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	db.connect = connect
	return db, nil
}

// AccountExist is the implementation of the interface
func (db *LevelDB) AccountExist(address common.Address) bool {
	var key = util.BytesCombine([]byte("account:"), address.Bytes())
	_, err := db.connect.Get(key, nil)
	if err != nil {
		return false
	}
	return true
}

// GetAccount returns an account of an address
func (db *LevelDB) GetAccount(address common.Address) (common.Account, error) {
	var key = util.BytesCombine([]byte("account:"), address.Bytes())
	value, err := db.connect.Get(key, nil)
	if err != nil {
		return common.NewDefaultAccount(address), nil
	}
	var account common.DefaultAccount
	err = json.Unmarshal(value, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// SetAccount updates an account or add an account
func (db *LevelDB) SetAccount(account common.Account) error {
	var key = util.BytesCombine([]byte("account:"), account.GetAddress().Bytes())
	value, err := account.Bytes()
	if err != nil {
		return err
	}
	err = db.connect.Put(key, value, nil)
	return err
}

// RemoveAccount removes an account if exist
func (db *LevelDB) RemoveAccount(address common.Address) error {
	var key = util.BytesCombine([]byte("account:"), address.Bytes())
	if ok, _ := db.connect.Has(key, nil); !ok {
		return nil
	}

	return db.connect.Delete(key, nil)
}

// GetStorage returns the key of an address if exist, else returns an error
func (db *LevelDB) GetStorage(address common.Address, key common.Word256) (common.Word256, error) {
	// return common.ZeroWord256, nil
	storageKey := util.BytesCombine(address.Bytes(), key.Bytes())
	value, err := db.connect.Get(storageKey, nil)
	if err != nil {
		return common.ZeroWord256, err
	}
	return common.BytesToWord256(value)
}

// SetStorage sets the value of a key belongs to an address
func (db *LevelDB) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	storageKey := util.BytesCombine(address.Bytes(), key.Bytes())
	db.connect.Put(storageKey, value.Bytes(), nil)
	return nil
}

// GetTxStatus is the implementation of interface
func (db *LevelDB) GetTxStatus(channelID, txID string) (*TxStatus, error) {
	var key = util.BytesCombine([]byte(channelID), []byte(txID))
	if ok, _ := db.connect.Has(key, nil); !ok {
		return nil, errors.New("Not exist")
	}
	value, err := db.connect.Get(key, nil)
	if err != nil {
		return nil, err
	}
	var status TxStatus
	err = json.Unmarshal(value, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetTxStatusAsync is the implementation of interface
func (db *LevelDB) GetTxStatusAsync(channelID, txID string) (*TxStatus, error) {
	var key = util.BytesCombine([]byte(channelID), []byte(txID))
	for {
		if ok, _ := db.connect.Has(key, nil); ok {
			value, err := db.connect.Get(key, nil)
			if err != nil {
				return nil, err
			}
			var status TxStatus
			err = json.Unmarshal(value, &status)
			if err != nil {
				return nil, err
			}
			return &status, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// SetTxStatus is the implementation of interface
func (db *LevelDB) SetTxStatus(tx *types.Tx, status *TxStatus) error {
	value, err := json.Marshal(status)
	if err != nil {
		return err
	}
	var key = util.BytesCombine([]byte(tx.Data.ChannelID), []byte(tx.ID))
	err = db.connect.Put(key, value, nil)
	if err != nil {
		return err
	}
	sender, err := tx.GetSender()
	if err != nil {
		return err
	}
	db.addHistory(sender.Bytes(), tx.Data.ChannelID, tx.ID)
	return nil
}

// BelongChannel is the implementation of interface
func (db *LevelDB) BelongChannel(channelID string) bool {
	channels := db.GetChannels()
	if util.Contain(channels, channelID) {
		return true
	}
	return false
}

// AddChannel is the implementation of interface
func (db *LevelDB) AddChannel(channelID string) {
	channels := db.GetChannels()
	if !util.Contain(channels, channelID) {
		channels = append(channels, channelID)
	}
	db.setChannels(channels)
}

// DeleteChannel is the implementation of interface
func (db *LevelDB) DeleteChannel(channelID string) {
	oldChannels := db.GetChannels()
	var newChannels []string
	for i := range oldChannels {
		if channelID != oldChannels[i] {
			newChannels = append(newChannels, oldChannels[i])
		}
	}
	db.setChannels(newChannels)
}

// GetChannels is the implementation of interface
func (db *LevelDB) GetChannels() []string {
	var channels []string
	var key = []byte("channels")
	if ok, _ := db.connect.Has(key, nil); !ok {
		return channels
	}
	value, err := db.connect.Get(key, nil)
	if err != nil {
		return channels
	}
	json.Unmarshal(value, &channels)
	return channels
}

// ListTxHistory is the implementation of interface
func (db *LevelDB) ListTxHistory(address []byte) map[string][]string {
	var txs = make(map[string][]string)
	if ok, _ := db.connect.Has(address, nil); ok {
		value, _ := db.connect.Get(address, nil)
		json.Unmarshal(value, &txs)
	}

	/*for channel, tx := range txs {
		for _, id := range tx {
			log.Infof("db/ListTxHistory: Channel %s, tx.ID %s", channel, id)
		}
	}*/

	return txs
}

func (db *LevelDB) addHistory(address []byte, channelID, txID string) {
	var txs = make(map[string][]string)
	if ok, _ := db.connect.Has(address, nil); !ok {
		txs[channelID] = []string{txID}
		value, _ := json.Marshal(txs)
		db.connect.Put(address, value, nil)
		log.Infoln("account ", address, " Channel ", channelID, "add ", txID, "to db")
	} else {
		value, err := db.connect.Get(address, nil)
		if err == nil {
			json.Unmarshal(value, &txs)
			if !util.Contain(txs, channelID) {
				txs[channelID] = []string{txID}
			} else {
				// todo: temporary fix the bug that duplicate txs, but this is not enough
				if !util.Contain(txs[channelID], txID) {
					txs[channelID] = append(txs[channelID], txID)
				}
			}
			value, _ := json.Marshal(txs)
			db.connect.Put(address, value, nil)
			log.Infoln("account ", address, " Channel ", channelID, "add ", txID, "to db")
		}
	}
}

func (db *LevelDB) setChannels(channels []string) {
	var key = []byte("channels")
	value, _ := json.Marshal(channels)
	db.connect.Put(key, value, nil)
}
