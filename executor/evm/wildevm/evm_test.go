package wildevm

import (
	"encoding/hex"
	"fmt"
	"madledger/common"
	"madledger/common/util"
	"madledger/executor/evm/wildevm/simulate"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ripemd160"
)

func newContext() *Context {
	return &Context{
		ChannelID: "test",
		Number:    0,
		BlockHash: common.ZeroHash,
		BlockTime: 0,
		GasLimit:  0,
		CoinBase:  common.ZeroWord256,
		Diffculty: 0,
	}
}

func newAccount(seed ...byte) common.Account {
	hasher := ripemd160.New()
	hasher.Write(seed)

	addr := common.BytesToAddress(hasher.Sum(nil))

	return common.NewDefaultAccount(addr)
}

func TestVmCall(t *testing.T) {
	db := simulate.NewStateDB()
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	// Create accounts
	account1 := newAccount(1)
	code := `6080604052600a600060005090905534801561001b5760006000fd5b50610021565b61026d806100306000396000f3fe60806040523480156100115760006000fd5b506004361061005c5760003560e01c80631003e2d21461006257806327ee58a6146100a5578063370158ea146100e857806360fe47b1146101395780636d4ce63c146101805761005c565b60006000fd5b61008f600480360360208110156100795760006000fd5b810190808035906020019092919050505061019e565b6040518082815260200191505060405180910390f35b6100d2600480360360208110156100bc5760006000fd5b81019080803590602001909291905050506101c6565b6040518082815260200191505060405180910390f35b6100f06101ee565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b610166600480360360208110156101505760006000fd5b8101908080359060200190929190505050610209565b604051808215151515815260200191505060405180910390f35b610188610225565b6040518082815260200191505060405180910390f35b6000816000600082828250540192505081909090555060006000505490506101c1565b919050565b6000816000600082828250540392505081909090555060006000505490506101e9565b919050565b600060003360006000505481915091509150610205565b9091565b600081600060005081909090555060019050610220565b919050565b60006000600050549050610234565b9056fea26469706673582212206d1f7e72f2d26816fe48ff60de6fa42d7b6fb40fc1603494b8c63221cd3c2c2364736f6c63430006000033`
	byteCode, _ := util.HexToBytes(code)
	output, address, err := vm.Create(account1, byteCode, []byte{}, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", output)
	contract, _ := db.GetAccount(address)
	input, _ := util.HexToBytes("6d4ce63c")
	output, err = vm.Call(account1, contract, output, input, 0)
	require.NoError(t, err)
	fmt.Printf("%x\n", output)
}

// Runs a basic loop
func TestVM(t *testing.T) {
	db := simulate.NewStateDB()
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	// Create accounts
	account1 := newAccount(1)
	account2 := newAccount(1, 0, 1)

	// var gas uint64 = 100000

	bytecode := MustSplice(PUSH1, 0x00, PUSH1, 0x20, MSTORE, JUMPDEST, PUSH2, 0x0F, 0x0F, PUSH1, 0x20, MLOAD,
		SLT, ISZERO, PUSH1, 0x1D, JUMPI, PUSH1, 0x01, PUSH1, 0x20, MLOAD, ADD, PUSH1, 0x20,
		MSTORE, PUSH1, 0x05, JUMP, JUMPDEST)

	start := time.Now()
	output, err := vm.Call(account1, account2, bytecode, []byte{}, 0)
	fmt.Printf("Output: %v Error: %v\n", output, err)
	fmt.Println("Call took:", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSHL(t *testing.T) {
	db := simulate.NewStateDB()
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)
	account1 := newAccount(1)
	account2 := newAccount(1, 0, 1)

	//Shift left 0
	bytecode := MustSplice(PUSH1, 0x01, PUSH1, 0x00, SHL, return1())
	output, err := vm.Call(account1, account2, bytecode, []byte{}, 0)
	value := []uint8([]byte{0x1})
	expected := common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift left 0
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0x00, SHL, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift left 1
	bytecode = MustSplice(PUSH1, 0x01, PUSH1, 0x01, SHL, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x2})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift left 1
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0x01, SHL, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift left 1
	bytecode = MustSplice(PUSH32, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0x01, SHL, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift left 255
	bytecode = MustSplice(PUSH1, 0x01, PUSH1, 0xFF, SHL, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x80})
	expected = common.RightPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift left 255
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0xFF, SHL, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x80})
	expected = common.RightPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift left 256 (overflow)
	bytecode = MustSplice(PUSH1, 0x01, PUSH2, 0x01, 0x00, SHL, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift left 256 (overflow)
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH2, 0x01, 0x00, SHL,
		return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift left 257 (overflow)
	bytecode = MustSplice(PUSH1, 0x01, PUSH2, 0x01, 0x01, SHL, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

}

func TestSHR(t *testing.T) {
	db := simulate.NewStateDB()
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)
	account1 := newAccount(1)
	account2 := newAccount(1, 0, 1)

	//Shift right 0
	bytecode := MustSplice(PUSH1, 0x01, PUSH1, 0x00, SHR, return1())
	output, err := vm.Call(account1, account2, bytecode, []byte{}, 0)
	value := []uint8([]byte{0x1})
	expected := common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift right 0
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0x00, SHR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift right 1
	bytecode = MustSplice(PUSH1, 0x01, PUSH1, 0x01, SHR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift right 1
	bytecode = MustSplice(PUSH32, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH1, 0x01, SHR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x40})
	expected = common.RightPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift right 1
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0x01, SHR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift right 255
	bytecode = MustSplice(PUSH32, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH1, 0xFF, SHR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x1})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift right 255
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0xFF, SHR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x1})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift right 256 (underflow)
	bytecode = MustSplice(PUSH32, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH2, 0x01, 0x00, SHR,
		return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift right 256 (underflow)
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH2, 0x01, 0x00, SHR,
		return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift right 257 (underflow)
	bytecode = MustSplice(PUSH32, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH2, 0x01, 0x01, SHR,
		return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

}

func TestSAR(t *testing.T) {
	db := simulate.NewStateDB()
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)
	account1 := newAccount(1)
	account2 := newAccount(1, 0, 1)

	//Shift arith right 0
	bytecode := MustSplice(PUSH1, 0x01, PUSH1, 0x00, SAR, return1())
	output, err := vm.Call(account1, account2, bytecode, []byte{}, 0)
	value := []uint8([]byte{0x1})
	expected := common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative arith shift right 0
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0x00, SAR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift arith right 1
	bytecode = MustSplice(PUSH1, 0x01, PUSH1, 0x01, SAR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift arith right 1
	bytecode = MustSplice(PUSH32, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH1, 0x01, SAR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0xc0})
	expected = common.RightPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift arith right 1
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0x01, SAR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift arith right 255
	bytecode = MustSplice(PUSH32, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH1, 0xFF, SAR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift arith right 255
	bytecode = MustSplice(PUSH32, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0xFF, SAR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift arith right 255
	bytecode = MustSplice(PUSH32, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH1, 0xFF, SAR, return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.RightPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift arith right 256 (reset)
	bytecode = MustSplice(PUSH32, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH2, 0x01, 0x00, SAR,
		return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Alternative shift arith right 256 (reset)
	bytecode = MustSplice(PUSH32, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, PUSH2, 0x01, 0x00, SAR,
		return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	value = []uint8([]byte{0x00})
	expected = common.LeftPadBytes(value, 32)
	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

	//Shift arith right 257 (reset)
	bytecode = MustSplice(PUSH32, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH2, 0x01, 0x01, SAR,
		return1())
	output, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
	expected = []uint8([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	assert.Equal(t, expected, output)

	t.Logf("Result: %v == %v\n", output, expected)

	if err != nil {
		t.Fatal(err)
	}

}

//Test attempt to jump to bad destination (position 16)
func TestJumpErr(t *testing.T) {
	db := simulate.NewStateDB()
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	// Create accounts
	account1 := newAccount(1)
	account2 := newAccount(2)

	bytecode := MustSplice(PUSH1, 0x10, JUMP)

	var err error
	ch := make(chan struct{})
	go func() {
		_, err = vm.Call(account1, account2, bytecode, []byte{}, 0)
		ch <- struct{}{}
	}()
	tick := time.NewTicker(time.Second * 2)
	select {
	case <-tick.C:
		t.Fatal("VM ended up in an infinite loop from bad jump dest (it took too long!)")
	case <-ch:
		if err == nil {
			t.Fatal("Expected invalid jump dest err")
		}
	}
}

// Tests the code for a subcurrency contract compiled by serpent
func TestSubcurrency(t *testing.T) {
	db := simulate.NewStateDB()
	// Create accounts
	account1 := newAccount(1, 2, 3)
	account2 := newAccount(3, 2, 1)
	db.SetAccount(account1)
	db.SetAccount(account2)
	fmt.Printf("%s\n", account1.GetAddress().String())
	fmt.Printf("%s\n", account2.GetAddress().String())

	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	bytecode := MustSplice(PUSH3, 0x0F, 0x42, 0x40, CALLER, SSTORE, PUSH29, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, PUSH1,
		0x00, CALLDATALOAD, DIV, PUSH4, 0x15, 0xCF, 0x26, 0x84, DUP2, EQ, ISZERO, PUSH2,
		0x00, 0x46, JUMPI, PUSH1, 0x04, CALLDATALOAD, PUSH1, 0x40, MSTORE, PUSH1, 0x40,
		MLOAD, SLOAD, PUSH1, 0x60, MSTORE, PUSH1, 0x20, PUSH1, 0x60, RETURN, JUMPDEST,
		PUSH4, 0x69, 0x32, 0x00, 0xCE, DUP2, EQ, ISZERO, PUSH2, 0x00, 0x87, JUMPI, PUSH1,
		0x04, CALLDATALOAD, PUSH1, 0x80, MSTORE, PUSH1, 0x24, CALLDATALOAD, PUSH1, 0xA0,
		MSTORE, CALLER, SLOAD, PUSH1, 0xC0, MSTORE, CALLER, PUSH1, 0xE0, MSTORE, PUSH1,
		0xA0, MLOAD, PUSH1, 0xC0, MLOAD, SLT, ISZERO, ISZERO, PUSH2, 0x00, 0x86, JUMPI,
		PUSH1, 0xA0, MLOAD, PUSH1, 0xC0, MLOAD, SUB, PUSH1, 0xE0, MLOAD, SSTORE, PUSH1,
		0xA0, MLOAD, PUSH1, 0x80, MLOAD, SLOAD, ADD, PUSH1, 0x80, MLOAD, SSTORE, JUMPDEST,
		JUMPDEST, POP, JUMPDEST, PUSH1, 0x00, PUSH1, 0x00, RETURN)

	data, _ := hex.DecodeString("693200CE0000000000000000000000004B4363CDE27C2EB05E66357DB05BC5C88F850C1A0000000000000000000000000000000000000000000000000000000000000005")
	output, err := vm.Call(account1, account2, bytecode, data, 0)
	fmt.Printf("Output: %v Error: %v\n", output, err)
	if err != nil {
		t.Fatal(err)
	}
}

//This test case is taken from EIP-140 (https://github.com/ethereum/EIPs/blob/master/EIPS/eip-140.md);
//it is meant to test the implementation of the REVERT opcode
func TestRevert(t *testing.T) {
	db := simulate.NewStateDB()
	// Create accounts
	account1 := newAccount(1)
	account2 := newAccount(1, 0, 1)
	key, value := []byte{0x00}, []byte{0x00}
	db.SetAccount(account1)
	db.SetStorage(account1.GetAddress(), common.LeftPadWord256(key), common.LeftPadWord256(value))
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	bytecode := MustSplice(PUSH13, 0x72, 0x65, 0x76, 0x65, 0x72, 0x74, 0x65, 0x64, 0x20, 0x64, 0x61, 0x74, 0x61,
		PUSH1, 0x00, SSTORE, PUSH32, 0x72, 0x65, 0x76, 0x65, 0x72, 0x74, 0x20, 0x6D, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		PUSH1, 0x00, MSTORE, PUSH1, 0x0E, PUSH1, 0x00, REVERT)

	output, cErr := vm.Call(account1, account2, bytecode, []byte{}, 0)
	assert.Error(t, cErr, "Expected execution reverted error")

	storageVal, err := db.GetStorage(account1.GetAddress(), common.LeftPadWord256(key))
	assert.Equal(t, common.LeftPadWord256(value), storageVal)

	t.Logf("Output: %v Error: %v\n", output, err)
}

// Test sending tokens from a contract to another account
// func TestSendCall(t *testing.T) {
// 	db := simulate.NewStateDB()
// 	// Create accounts
// 	account1 := newAccount(1)
// 	account2 := newAccount(2)
// 	account3 := newAccount(3)

// 	db.SetAccount(account1)
// 	db.SetAccount(account2)
// 	db.SetAccount(account3)
// 	vm := NewEVM(newContext(), common.ZeroAddress, db)

// 	// account1 will call account2 which will trigger CALL opcode to account3
// 	addr := account3.GetAddress()
// 	contractCode := callContractCode(addr)

// 	//----------------------------------------------
// 	// account2 has insufficient balance, should fail
// 	_, err := runVMWaitError(vm, &account1, &account2, addr, contractCode, 1000)
// 	assert.Error(t, err, "Expected insufficient balance error")

// 	//----------------------------------------------
// 	// give account2 sufficient balance, should pass
// 	account2 = newAccount(2)
// 	err = account2.AddBalance(100000)
// 	require.NoError(t, err)
// 	_, err = runVMWaitError(vm, &account1, &account2, addr, contractCode, 1000)
// 	assert.NoError(t, err, "Should have sufficient balance")

// 	//----------------------------------------------
// 	// insufficient gas, should fail
// 	account2 = newAccount(2)
// 	err = account2.AddBalance(100000)
// 	require.NoError(t, err)
// 	_, err = runVMWaitError(vm, &account1, &account2, addr, contractCode, 100)
// 	assert.NoError(t, err, "Expected insufficient gas error")
// }

// // This test was introduced to cover an issues exposed in our handling of the
// // gas limit passed from caller to callee on various forms of CALL.
// // The idea of this test is to implement a simple DelegateCall in EVM code
// // We first run the DELEGATECALL with _just_ enough gas expecting a simple return,
// // and then run it with 1 gas unit less, expecting a failure
// func TestDelegateCallGas(t *testing.T) {
// 	cache := state.NewCache(newAppState())
// 	vm := NewVM(newParams(), crypto.ZeroAddress, nil, logger)

// 	inOff := 0
// 	inSize := 0 // no call data
// 	retOff := 0
// 	retSize := 32
// 	calleeReturnValue := int64(20)

// 	// DELEGATECALL(retSize, refOffset, inSize, inOffset, addr, gasLimit)
// 	// 6 pops
// 	delegateCallCost := GasStackOp * 6
// 	// 1 push
// 	gasCost := GasStackOp
// 	// 2 pops, 1 push
// 	subCost := GasStackOp * 3
// 	pushCost := GasStackOp

// 	costBetweenGasAndDelegateCall := gasCost + subCost + delegateCallCost + pushCost

// 	// Do a simple operation using 1 gas unit
// 	calleeAccount, calleeAddress := makeAccountWithCode(cache, "callee",
// 		MustSplice(PUSH1, calleeReturnValue, return1()))

// 	// Here we split up the caller code so we can make a DELEGATE call with
// 	// different amounts of gas. The value we sandwich in the middle is the amount
// 	// we subtract from the available gas (that the caller has available), so:
// 	// code := MustSplice(callerCodePrefix, <amount to subtract from GAS> , callerCodeSuffix)
// 	// gives us the code to make the call
// 	callerCodePrefix := MustSplice(PUSH1, retSize, PUSH1, retOff, PUSH1, inSize,
// 		PUSH1, inOff, PUSH20, calleeAddress, PUSH1)
// 	callerCodeSuffix := MustSplice(GAS, SUB, DELEGATECALL, returnWord())

// 	// Perform a delegate call
// 	callerAccount, _ := makeAccountWithCode(cache, "caller",
// 		MustSplice(callerCodePrefix,
// 			// Give just enough gas to make the DELEGATECALL
// 			costBetweenGasAndDelegateCall,
// 			callerCodeSuffix))

// 	// Should pass
// 	output, err := runVMWaitError(cache, vm, callerAccount, calleeAccount, calleeAddress,
// 		callerAccount.Code(), 100)
// 	assert.NoError(t, err, "Should have sufficient funds for call")
// 	assert.Equal(t, Int64ToWord256(calleeReturnValue).Bytes(), output)

// 	callerAccount.SetCode(MustSplice(callerCodePrefix,
// 		// Shouldn't be enough gas to make call
// 		costBetweenGasAndDelegateCall-1,
// 		callerCodeSuffix))

// 	// Should fail
// 	_, err = runVMWaitError(cache, vm, callerAccount, calleeAccount, calleeAddress,
// 		callerAccount.Code(), 100)
// 	assert.Error(t, err, "Should have insufficient gas for call")
// }

func TestMemoryBounds(t *testing.T) {
	db := simulate.NewStateDB()
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)
	vm.memoryProvider = func() Memory {
		return NewDynamicMemory(1024, 2048)
	}
	caller, _ := makeAccountWithCode(db, "caller", nil)
	callee, _ := makeAccountWithCode(db, "callee", nil)
	// gas := uint64(100000)
	// This attempts to store a value at the memory boundary and return it
	word := common.OneWord256
	output, err := vm.call(caller, callee,
		MustSplice(pushWord(word), storeAtEnd(), MLOAD, storeAtEnd(), returnAfterStore()),
		nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, word.Bytes(), output)

	// Same with number
	word = common.Int64ToWord256(232234234432)
	output, err = vm.call(caller, callee,
		MustSplice(pushWord(word), storeAtEnd(), MLOAD, storeAtEnd(), returnAfterStore()),
		nil, 0)
	// assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, word.Bytes(), output)

	// Now test a series of boundary stores
	code := pushWord(word)
	for i := 0; i < 10; i++ {
		code = MustSplice(code, storeAtEnd(), MLOAD)
	}
	output, err = vm.call(caller, callee, MustSplice(code, storeAtEnd(), returnAfterStore()),
		nil, 0)
	// assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, word.Bytes(), output)

	// Same as above but we should breach the upper memory limit set in memoryProvider
	code = pushWord(word)
	for i := 0; i < 100; i++ {
		code = MustSplice(code, storeAtEnd(), MLOAD)
	}
	output, err = vm.call(caller, callee, MustSplice(code, storeAtEnd(), returnAfterStore()),
		nil, 0)
	assert.Error(t, err, "Should hit memory out of bounds")
}

func TestMsgSender(t *testing.T) {
	db := simulate.NewStateDB()
	account1 := newAccount(1, 2, 3)
	account2 := newAccount(3, 2, 1)
	db.SetAccount(account1)
	db.SetAccount(account2)

	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	// var gas uint64 = 100000

	/*
			pragma solidity ^0.4.0;

			contract SimpleStorage {
		                function get() public constant returns (address) {
		        	        return msg.sender;
		    	        }
			}
	*/

	// This bytecode is compiled from Solidity contract above using remix.ethereum.org online compiler
	code, err := hex.DecodeString("6060604052341561000f57600080fd5b60ca8061001d6000396000f30060606040526004361060" +
		"3f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680636d4ce63c14604457" +
		"5b600080fd5b3415604e57600080fd5b60546096565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ff" +
		"ffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000339050905600a165627a" +
		"7a72305820b9ebf49535372094ae88f56d9ad18f2a79c146c8f56e7ef33b9402924045071e0029")
	require.NoError(t, err)

	// Run the contract initialisation code to obtain the contract code that would be mounted at account2
	contractCode, errCode := vm.Call(account1, account2, code, code, 0)
	// require.NoError(t, err)
	if errCode != nil {
		t.Fatal(errCode)
	}

	// Not needed for this test (since contract code is passed as argument to vm), but this is what an execution
	// framework must do
	account2.SetCode(contractCode)

	// Input is the function hash of `get()`
	input, err := hex.DecodeString("6d4ce63c")

	output, errCode := vm.Call(account1, account2, contractCode, input, 0)
	if errCode != nil {
		t.Fatal(errCode)
	}

	assert.Equal(t, account1.GetAddress().Word256().Bytes(), output)

}

func TestInvalid(t *testing.T) {
	db := simulate.NewStateDB()
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	// Create accounts
	account1 := newAccount(1)
	account2 := newAccount(1, 0, 1)

	// var gas uint64 = 100000

	bytecode := MustSplice(PUSH32, 0x72, 0x65, 0x76, 0x65, 0x72, 0x74, 0x20, 0x6D, 0x65, 0x73, 0x73, 0x61,
		0x67, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, PUSH1, 0x00, MSTORE, PUSH1, 0x0E, PUSH1, 0x00, INVALID)

	_, err := vm.Call(account1, account2, bytecode, []byte{}, 0)
	// assert.Equal(t, errors.ErrorCodeExecutionAborted, err.ErrorCode())
	if err.Error() != "Execution aborted" {
		t.Fatal(err)
	}
	// t.Logf("Output: %v Error: %v\n", output, err)
}

func TestReturnDataSize(t *testing.T) {
	db := simulate.NewStateDB()
	account1 := newAccount(1)
	// account2
	accountName := "account2addresstests"
	callcode := MustSplice(PUSH32, 0x72, 0x65, 0x76, 0x65, 0x72, 0x74, 0x20, 0x6D, 0x65, 0x73, 0x73, 0x61,
		0x67, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, PUSH1, 0x00, MSTORE, PUSH1, 0x0E, PUSH1, 0x00, RETURN)
	account2, _ := makeAccountWithCode(db, accountName, callcode)
	db.SetAccount(account2)
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	gas1, gas2 := byte(0x1), byte(0x1)
	value := byte(0x69)
	inOff, inSize := byte(0x0), byte(0x0) // no call data
	retOff, retSize := byte(0x0), byte(0x0E)

	bytecode := MustSplice(PUSH1, retSize, PUSH1, retOff, PUSH1, inSize, PUSH1, inOff, PUSH1, value, PUSH20,
		0x61, 0x63, 0x63, 0x6F, 0x75, 0x6E, 0x74, 0x32, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65,
		0x73, 0x74, 0x73, PUSH2, gas1, gas2, CALL, RETURNDATASIZE, PUSH1, 0x00, MSTORE, PUSH1, 0x20, PUSH1, 0x00, RETURN)

	expected := common.LeftPadBytes([]byte{0x0E}, 32)

	output, err := vm.Call(account1, account2, bytecode, []byte{}, 0)

	assert.Equal(t, expected, output)

	t.Logf("Output: %v Error: %v\n", output, err)

	if err != nil {
		t.Fatal(err)
	}
}

func TestReturnDataCopy(t *testing.T) {
	db := simulate.NewStateDB()
	accountName := "account2addresstests"

	callcode := MustSplice(PUSH32, 0x72, 0x65, 0x76, 0x65, 0x72, 0x74, 0x20, 0x6D, 0x65, 0x73, 0x73, 0x61,
		0x67, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, PUSH1, 0x00, MSTORE, PUSH1, 0x0E, PUSH1, 0x00, RETURN)

	// Create accounts
	account1 := newAccount(1)
	account2, _ := makeAccountWithCode(db, accountName, callcode)
	db.SetAccount(account2)
	wb := db.NewWriteBatch()
	vm := NewEVM(newContext(), common.ZeroAddress, db, wb)

	gas1, gas2 := byte(0x1), byte(0x1)
	value := byte(0x69)
	inOff, inSize := byte(0x0), byte(0x0) // no call data
	retOff, retSize := byte(0x0), byte(0x0E)

	bytecode := MustSplice(PUSH1, retSize, PUSH1, retOff, PUSH1, inSize, PUSH1, inOff, PUSH1, value, PUSH20,
		0x61, 0x63, 0x63, 0x6F, 0x75, 0x6E, 0x74, 0x32, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65,
		0x73, 0x74, 0x73, PUSH2, gas1, gas2, CALL, RETURNDATASIZE, PUSH1, 0x00, PUSH1, 0x00, RETURNDATACOPY,
		RETURNDATASIZE, PUSH1, 0x00, RETURN)

	expected := []byte{0x72, 0x65, 0x76, 0x65, 0x72, 0x74, 0x20, 0x6D, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65}

	output, err := vm.Call(account1, account2, bytecode, []byte{}, 0)

	assert.Equal(t, expected, output)

	t.Logf("Output: %v Error: %v\n", output, err)

	if err != nil {
		t.Fatal(err)
	}
}

// These code segment helpers exercise the MSTORE MLOAD MSTORE cycle to test
// both of the memory operations. Each MSTORE is done on the memory boundary
// (at MSIZE) which Solidity uses to find guaranteed unallocated memory.

// storeAtEnd expects the value to be stored to be on top of the stack, it then
// stores that value at the current memory boundary
func storeAtEnd() []byte {
	// Pull in MSIZE (to carry forward to MLOAD), swap in value to store, store it at MSIZE
	return MustSplice(MSIZE, SWAP1, DUP2, MSTORE)
}

func returnAfterStore() []byte {
	return MustSplice(PUSH1, 32, DUP2, RETURN)
}

// Store the top element of the stack (which is a 32-byte word) in memory
// and return it. Useful for a simple return value.
func return1() []byte {
	return MustSplice(PUSH1, 0, MSTORE, returnWord())
}

func returnWord() []byte {
	// PUSH1 => return size, PUSH1 => return offset, RETURN
	return MustSplice(PUSH1, 32, PUSH1, 0, RETURN)
}

func makeAccountWithCode(db StateDB, name string, code []byte) (common.Account, common.Address) {
	address, _ := common.AddressFromBytes([]byte(name))
	account := common.NewDefaultAccount(address)
	account.AddBalance(9999999)
	account.SetCode(code)
	db.SetAccount(account)
	return account, account.GetAddress()
}

// Subscribes to an AccCall, runs the vm, returns the output any direct exception
// and then waits for any exceptions transmitted by Data in the AccCall
// event (in the case of no direct error from call we will block waiting for
// at least 1 AccCall event)
func runVMWaitError(vm *EVM, caller, callee *common.Account, subscribeAddr common.Address,
	contractCode []byte, gas uint64) ([]byte, error) {
	// txe := new(exec.TxExecution)
	output, err := runVM(vm, caller, callee, subscribeAddr, contractCode, gas)
	if err != nil {
		return output, err
	}
	// if len(txe.Events) > 0 {
	// 	ex := txe.Events[0].Header.Exception
	// 	if ex != nil {
	// 		return output, ex
	// 	}
	// }
	return output, nil
}

// Subscribes to an AccCall, runs the vm, returns the output and any direct
// exception
func runVM(vm *EVM, caller, callee *common.Account,
	subscribeAddr common.Address, contractCode []byte, gas uint64) ([]byte, error) {

	fmt.Printf("subscribe to %s\n", subscribeAddr)

	// vm.SetEventSink(sink)
	start := time.Now()
	output, err := vm.Call(*caller, *callee, contractCode, []byte{}, 0)
	fmt.Printf("Output: %v Error: %v\n", output, err)
	fmt.Println("Call took:", time.Since(start))
	return output, err
}

// this is code to call another contract (hardcoded as addr)
func callContractCode(addr common.Address) []byte {
	gas1, gas2 := byte(0x1), byte(0x1)
	value := byte(0x69)
	inOff, inSize := byte(0x0), byte(0x0) // no call data
	retOff, retSize := byte(0x0), byte(0x20)
	// this is the code we want to run (send funds to an account and return)
	return MustSplice(PUSH1, retSize, PUSH1, retOff, PUSH1, inSize, PUSH1,
		inOff, PUSH1, value, PUSH20, addr, PUSH2, gas1, gas2, CALL, PUSH1, retSize,
		PUSH1, retOff, RETURN)
}

// func pushInt64(i int64) []byte {
// 	return pushWord(Int64ToWord256(i))
// }

// Produce bytecode for a PUSH<N>, b_1, ..., b_N where the N is number of bytes
// contained in the unpadded word
func pushWord(word common.Word256) []byte {
	leadingZeros := byte(0)
	for leadingZeros < 32 {
		if word[leadingZeros] == 0 {
			leadingZeros++
		} else {
			return MustSplice(byte(PUSH32)-leadingZeros, word[leadingZeros:])
		}
	}
	fmt.Printf("push word get %b\n", MustSplice(PUSH1, 0))
	return MustSplice(PUSH1, 0)
}

func TestPushWord(t *testing.T) {
	word := common.Int64ToWord256(int64(2133213213))
	assert.Equal(t, MustSplice(PUSH4, 0x7F, 0x26, 0x40, 0x1D), pushWord(word))
	word[0] = 1
	assert.Equal(t, MustSplice(PUSH32,
		1, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0x7F, 0x26, 0x40, 0x1D), pushWord(word))
	assert.Equal(t, MustSplice(PUSH1, 0), pushWord(common.Word256{}))
	assert.Equal(t, MustSplice(PUSH1, 1), pushWord(common.Int64ToWord256(1)))
}

// Kind of indirect test of Splice, but here to avoid import cycles
// func TestBytecode(t *testing.T) {
// 	assert.Equal(t,
// 		MustSplice(1, 2, 3, 4, 5, 6),
// 		MustSplice(1, 2, 3, MustSplice(4, 5, 6)))
// 	assert.Equal(t,
// 		MustSplice(1, 2, 3, 4, 5, 6, 7, 8),
// 		MustSplice(1, 2, 3, MustSplice(4, MustSplice(5), 6), 7, 8))
// 	assert.Equal(t,
// 		MustSplice(PUSH1, 2),
// 		MustSplice(byte(PUSH1), 0x02))
// 	assert.Equal(t,
// 		[]byte{},
// 		MustSplice(MustSplice(MustSplice())))

// 	contractAccount := &common.Account{Address: common.AddressFromWord256(common.Int64ToWord256(102))}
// 	addr := contractAccount.Address
// 	gas1, gas2 := byte(0x1), byte(0x1)
// 	value := byte(0x69)
// 	inOff, inSize := byte(0x0), byte(0x0) // no call data
// 	retOff, retSize := byte(0x0), byte(0x20)
// 	contractCodeBytecode := MustSplice(PUSH1, retSize, PUSH1, retOff, PUSH1, inSize, PUSH1,
// 		inOff, PUSH1, value, PUSH20, addr, PUSH2, gas1, gas2, CALL, PUSH1, retSize,
// 		PUSH1, retOff, RETURN)
// 	contractCode := []byte{0x60, retSize, 0x60, retOff, 0x60, inSize, 0x60, inOff, 0x60, value, 0x73}
// 	contractCode = append(contractCode, addr[:]...)
// 	contractCode = append(contractCode, []byte{0x61, gas1, gas2, 0xf1, 0x60, 0x20, 0x60, 0x0, 0xf3}...)
// 	assert.Equal(t, contractCode, contractCodeBytecode)
// }

// func TestConcat(t *testing.T) {
// 	assert.Equal(t,
// 		[]byte{0x01, 0x02, 0x03, 0x04},
// 		Concat([]byte{0x01, 0x02}, []byte{0x03, 0x04}))
// }

func TestSubslice(t *testing.T) {
	const size = 10
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(i)
	}
	for n := int64(0); n < size; n++ {
		data = data[:n]
		for offset := int64(-size); offset < size; offset++ {
			for length := int64(-size); length < size; length++ {
				_, ok := subslice(data, offset, length)
				if offset < 0 || length < 0 || n < offset {
					assert.False(t, ok)
				} else {
					assert.True(t, ok)
				}
			}
		}
	}
}

type ByteSlicable interface {
	Bytes() []byte
}

// Splice or panic
func MustSplice(bytelikes ...interface{}) []byte {
	spliced, err := Splice(bytelikes...)
	if err != nil {
		panic(err)
	}
	return spliced
}

// Convenience function to allow us to mix bytes, ints, and OpCodes that
// represent bytes in an EVM assembly code to make assembly more readable.
// Also allows us to splice together assembly
// fragments because any []byte arguments are flattened in the result.
func Splice(bytelikes ...interface{}) ([]byte, error) {
	bytes := make([]byte, 0, len(bytelikes))
	for _, bytelike := range bytelikes {
		bs, err := byteSlicify(bytelike)
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, bs...)
	}
	return bytes, nil
}

// Convert anything byte or byte slice like to a byte slice
func byteSlicify(bytelike interface{}) ([]byte, error) {
	switch b := bytelike.(type) {
	case byte:
		return []byte{b}, nil
	case OpCode:
		return []byte{byte(b)}, nil
	case int:
		if int(byte(b)) != b {
			return nil, fmt.Errorf("the int %v does not fit inside a byte", b)
		}
		return []byte{byte(b)}, nil
	case int64:
		if int64(byte(b)) != b {
			return nil, fmt.Errorf("the int64 %v does not fit inside a byte", b)
		}
		return []byte{byte(b)}, nil
	case uint64:
		if uint64(byte(b)) != b {
			return nil, fmt.Errorf("the uint64 %v does not fit inside a byte", b)
		}
		return []byte{byte(b)}, nil
	case string:
		return []byte(b), nil
	case ByteSlicable:
		return b.Bytes(), nil
	case []byte:
		return b, nil
	default:
		return nil, fmt.Errorf("could not convert %s to a byte or sequence of bytes", bytelike)
	}
}
