//  Copyright 2020 The THU-Arxan Authors
//  This file is part of the evm library.
//
//  The evm library is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  The evm library is distributed in the hope that it will be useful,/
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with the evm library. If not, see <http://www.gnu.org/licenses/>.
//

package tests

import (
	"github.com/thu-arxan/evm"
	"github.com/thu-arxan/evm/db"
	"github.com/thu-arxan/evm/example"
	"github.com/thu-arxan/evm/util"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	hashAbi     = "sols/Hash_sol_Hash.abi"
	hashBin     = "sols/Hash_sol_Hash.bin"
	hashCode    []byte
	hashAddress evm.Address
)

func TestHashSol(t *testing.T) {
	var err error
	binBytes, err := util.ReadBinFile(hashBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var exceptCode = `60806040523480156100115760006000fd5b506004361061003b5760003560e01c8063a0fe620214610041578063fadf2722146101205761003b565b60006000fd5b610102600480360360208110156100585760006000fd5b81019080803590602001906401000000008111156100765760006000fd5b8201836020820111156100895760006000fd5b803590602001918460018302840111640100000000831117156100ac5760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509090919290909192905050506101ff565b60405180826000191660001916815260200191505060405180910390f35b6101e1600480360360208110156101375760006000fd5b81019080803590602001906401000000008111156101555760006000fd5b8201836020820111156101685760006000fd5b8035906020019184600183028401116401000000008311171561018b5760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505090909192909091929050505061033c565b60405180826000191660001916815260200191505060405180910390f35b600060006002836040516020018080602001828103825283818151815260200191508051906020019080838360005b8381101561024a5780820151818401525b60208101905061022e565b50505050905090810190601f1680156102775780820380516001836020036101000a031916815260200191505b50925050506040516020818303038152906040526040518082805190602001908083835b6020831015156102c157805182525b60208201915060208101905060208303925061029b565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610304573d600060003e3d6000fd5b5050506040513d602081101561031a5760006000fd5b810190808051906020019092919050505090508091505061033756505b919050565b60006000826040516020018080602001828103825283818151815260200191508051906020019080838360005b838110156103855780820151818401525b602081019050610369565b50505050905090810190601f1680156103b25780820380516001836020036101000a031916815260200191505b5092505050604051602081830303815290604052805190602001209050809150506103d956505b91905056fea2646970667358221220bf2039bb6ae8a0093d2a34a2985062c668b5754efe465cbac6dec472f12b207a64736f6c63430006000033`
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	hashCode, hashAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, exceptCode, 209063)
	// then call
	callWithPayload(t, memoryDB, bc, origin, hashAddress, mustPack(hashAbi, "SHA256", "hello"), 2528, 0)
	// payload, err = parsePayload(hashAbi, "KECCAK256", "hello")
	// require.NoError(t, err)
	callWithPayload(t, memoryDB, bc, origin, hashAddress, mustPack(hashAbi, "KECCAK256", "hello"), 1170, 0)
}
