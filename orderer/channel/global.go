package channel

import "madledger/core/types"

// AddGlobalBlock add a global block
// Note: It should not add block file again.
// TODO: update something in the db
func (manager *Manager) AddGlobalBlock(block *types.Block) error {
	return nil
}
