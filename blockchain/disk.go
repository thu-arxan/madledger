package blockchain

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"madledger/common/util"
	"madledger/core/types"
	"os"
	"strconv"
)

// load loads a channel and return the block excepted
// if the envirnoment is not completed, then complete the env
// todo: the damaged circumstance is not consided completed
func load(dir string) (uint64, error) {
	var err error
	genesisBlockPath, _ := util.MakeFileAbs("0.json", dir)
	cachePath, _ := util.MakeFileAbs(".cache", dir)
	if util.FileExists(dir) {
		gbExist := util.FileExists(genesisBlockPath)
		cacheExist := util.FileExists(cachePath)
		if gbExist && cacheExist {
			num, err := loadCache(cachePath)
			if err != nil {
				return 0, err
			}
			return num + 1, err
		}
		if gbExist && !cacheExist { // only exists gb, then the env is destoryed
			return 0, fmt.Errorf("Cache file of %s is damaged", dir)
		}
		if cacheExist && !gbExist {
			return 0, fmt.Errorf("Blocks files of %s are damaged", dir)
		}
		return 0, nil
	}
	err = initEnv(dir)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// initEnv init the env that a channel needs
func initEnv(dir string) error {
	if !util.FileExists(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadCache(path string) (uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var data string
	_, err = fmt.Fscanln(file, &data)
	if err != nil {
		return 0, err
	}
	num, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("The checkpoint file %s is damaged", path)
	}
	return num, nil
}

func (manager *Manager) loadBlock(num uint64) (*types.Block, error) {
	path := manager.getJSONStorePath(num)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var block types.Block
	// However, this can not use json.Unmarshal, because the pk is unable to unmarshal

	err = json.Unmarshal(data, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (manager *Manager) storeBlock(block *types.Block) error {
	path := manager.getJSONStorePath(block.Header.Number)
	if util.FileExists(path) {
		return errors.New("The block file is aleardy exists")
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)

	err = enc.Encode(block)
	if err != nil {
		return err
	}
	return nil
}

func (manager *Manager) updateCache(num uint64) error {
	cachePath, _ := util.MakeFileAbs(".cache", manager.dir)
	file, err := os.OpenFile(cachePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(fmt.Sprintf("%d", num))
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) getJSONStorePath(num uint64) string {
	path, _ := util.MakeFileAbs(fmt.Sprintf("%d.json", num), manager.dir)
	return path
}
