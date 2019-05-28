package protos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetChainAddress(t *testing.T) {
	eraftAddress := "localhost:12345"
	require.Equal(t, "localhost:12343", ERaftToChain(eraftAddress))

	eraftAddress = "localhost: 12345"
	require.Equal(t, "localhost:12343", ERaftToChain(eraftAddress))

	eraftAddress = "127.0.0.1: 12345"
	require.Equal(t, "127.0.0.1:12343", ERaftToChain(eraftAddress))

	eraftAddress = "127.0.0.1://12345"
	require.Equal(t, "", ERaftToChain(eraftAddress))
}
