package server

import (
	"madledger/core/types"
	"reflect"
	"testing"
)

func TestConvertBlock(t *testing.T) {
	typesBlock := types.NewBlock("test", 1, nil, nil)
	pbBlock, err := ConvertBlockFromTypesToPb(typesBlock)
	if err != nil {
		t.Fatal(err)
	}
	newTypesBlock, err := ConvertBlockFromPbToTypes(pbBlock)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(typesBlock, newTypesBlock) {
		t.Fatal()
	}
}
