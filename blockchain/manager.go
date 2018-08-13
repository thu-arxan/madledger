package blockchain

import (
	"madledger/util"
	"os"
)

// Manager manage the blockchain
type Manager struct {
	id  string
	dir string
}

// NewManager is the constructor of manager
func NewManager(id, dir string) (*Manager, error) {
	var m = Manager{
		id:  id,
		dir: dir,
	}
	if util.FileExists(dir) {
		return &m, nil
	}
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
