package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBlock(t *testing.T) {
	var block = NewBlock("", 0, nil, nil)
	require.EqualValues(t, GenesisBlockPrevHash, block.Header.PrevBlock)
	require.Len(t, block.GetMerkleRoot(), 32)
}
