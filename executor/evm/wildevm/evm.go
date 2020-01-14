package wildevm

import (
	"bytes"
	"fmt"
	"madledger/common"
	"madledger/common/crypto/sha3"
	"madledger/peer/db"
	"math/big"
)

const (
	// DefaultStackCapacity define the default capacity of stack
	DefaultStackCapacity = 1024
)

// EVM provide a environment to run ethereum compile codes
type EVM struct {
	context        *Context
	db             StateDB
	origin         common.Address
	stackDepth     uint64
	returnData     []byte
	cache          *Cache
	memoryProvider func() Memory
	wb             db.WriteBatch
}

// NewEVM is the constructor of EVM
func NewEVM(context *Context, origin common.Address, db StateDB, wb db.WriteBatch) *EVM {
	return &EVM{
		context:        context,
		db:             db,
		origin:         origin,
		stackDepth:     0,
		cache:          NewCache(db),
		memoryProvider: DefaultDynamicMemoryProvider,
		wb:             wb,
	}
}

// Create create a contract.
// If there exist a contract on the address then a error occurs.
func (evm *EVM) Create(caller common.Account, code, input []byte, value uint64) ([]byte, common.Address, error) {
	contract, err := evm.createAccount(caller, code)
	if err != nil {
		return nil, common.ZeroAddress, err
	}

	// Run the contract bytes and return the runtime bytes
	output, err := evm.Call(caller, contract, code, input, value)
	if err != nil {
		return nil, common.ZeroAddress, err
	}
	contract.SetCode(output)
	//err = evm.cache.SetAccount(contract)
	err = evm.wb.SetAccount(contract)
	if err != nil {
		return nil, common.ZeroAddress, err
	}
	evm.cache.Sync(evm.wb)

	return output, contract.GetAddress(), nil
}

// Call run code on a evm.
// Remember, the function will not add the nonce of caller.
func (evm *EVM) Call(caller, callee common.Account, code, input []byte, value uint64) ([]byte, error) {
	if err := transfer(caller, callee, value); err != nil {
		return nil, err
	}

	// Here also run some codes
	if len(code) > 0 {
		evm.stackDepth++
		output, err := evm.call(caller, callee, code, input, value)
		evm.stackDepth--
		if err != nil {
			return nil, err
		}
		evm.cache.Sync(evm.wb)
		return output, nil
	}
	return nil, nil
}

// DelegateCall is executed by the DELEGATECALL opcode, introduced as off Ethereum Homestead.
// The intent of delegate call is to run the code of the callee in the storage context of the caller;
// while preserving the original caller to the previous callee.
// Different to the normal CALL or CALLCODE, the value does not need to be transferred to the callee.
func (evm *EVM) DelegateCall(caller, callee common.Account, code, input []byte, value uint64) ([]byte, error) {
	if len(code) > 0 {
		evm.stackDepth++
		output, err := evm.call(caller, callee, code, input, value)
		evm.stackDepth--
		if err != nil {
			return nil, err
		}
		return output, nil
	}

	return nil, nil
}

// log is used for debug
func (evm *EVM) log(op OpCode, pc int64) {
	fmt.Printf("Op is %s, pc is %d\n", op, pc)
}

// Just like Call() but does not transfer 'value' or modify the callDepth.
func (evm *EVM) call(caller, callee common.Account, code, input []byte, value uint64) ([]byte, error) {
	var (
		pc     int64
		stack  = NewStack(DefaultStackCapacity)
		memory = evm.memoryProvider()
		cache  = evm.cache
		i      = 0
	)

	for {
		var op = codeGetOp(code, pc)
		// fmt.Printf(">>> Opcode[%d] is %s, pc is %d\n", i, op, pc)
		// stack.Print(10)
		i++
		switch op {
		case ADD: // 0x01
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			sum := new(big.Int).Add(x, y)
			if err = stack.PushBigInt(sum); err != nil {
				return nil, err
			}

		case MUL: // 0x02
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			prod := new(big.Int).Mul(x, y)
			if err = stack.PushBigInt(prod); err != nil {
				return nil, err
			}

		case SUB: // 0x03
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			diff := new(big.Int).Sub(x, y)
			if err = stack.PushBigInt(diff); err != nil {
				return nil, err
			}

		case DIV: // 0x04
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			if y.Sign() == 0 {
				if err = stack.Push(common.ZeroWord256); err != nil {
					return nil, err
				}
			} else {
				div := new(big.Int).Div(x, y)
				if err = stack.PushBigInt(div); err != nil {
					return nil, err
				}
			}

		case SDIV: // 0x05
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigSigned(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			if y.Sign() == 0 {
				if err = stack.Push(common.ZeroWord256); err != nil {
					return nil, err
				}
			} else {
				div := new(big.Int).Div(x, y)
				if err = stack.PushBigInt(div); err != nil {
					return nil, err
				}
			}

		case MOD: // 0x06
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			if y.Sign() == 0 {
				if err = stack.Push(common.ZeroWord256); err != nil {
					return nil, err
				}
			} else {
				mod := new(big.Int).Mod(x, y)
				if err = stack.PushBigInt(mod); err != nil {
					return nil, err
				}
			}

		case SMOD: // 0x07
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigSigned(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			if y.Sign() == 0 {
				if err = stack.Push(common.ZeroWord256); err != nil {
					return nil, err
				}
			} else {
				mod := new(big.Int).Mod(x, y)
				if err = stack.PushBigInt(mod); err != nil {
					return nil, err
				}
			}

		case ADDMOD: // 0x08
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(3); err != nil {
				return nil, err
			}
			x, y, z := values[0], values[1], values[2]
			if z.Sign() == 0 {
				if err = stack.Push(common.ZeroWord256); err != nil {
					return nil, err
				}
			} else {
				add := new(big.Int).Add(x, y)
				mod := add.Mod(add, z)
				if err = stack.PushBigInt(mod); err != nil {
					return nil, err
				}
			}

		case MULMOD: // 0x09
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(3); err != nil {
				return nil, err
			}
			x, y, z := values[0], values[1], values[2]
			if z.Sign() == 0 {
				if err = stack.Push(common.ZeroWord256); err != nil {
					return nil, err
				}
			} else {
				mul := new(big.Int).Mul(x, y)
				mod := mul.Mod(mul, z)
				if err = stack.PushBigInt(mod); err != nil {
					return nil, err
				}
			}

		case EXP: // 0x0A
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			pow := new(big.Int).Exp(x, y, nil)
			if err = stack.PushBigInt(pow); err != nil {
				return nil, err
			}

		case SIGNEXTEND: // 0x0B
			var back uint64
			var err error
			back, err = stack.PopU64()
			if err != nil {
				return nil, err
			}
			if back < common.Word256Length-1 {
				var x *big.Int
				if x, err = stack.PopBigInt(); err != nil {
					return nil, err
				}
				if err = stack.PushBigInt(common.SignExtend(back, x)); err != nil {
					return nil, err
				}
			}

		case LT: // 0x10
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			var word = common.ZeroWord256
			if x.Cmp(y) < 0 {
				word = common.OneWord256
			}
			if err = stack.Push(word); err != nil {
				return nil, err
			}

		case GT: // 0x11
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			var word = common.ZeroWord256
			if x.Cmp(y) > 0 {
				word = common.OneWord256
			}
			if err = stack.Push(word); err != nil {
				return nil, err
			}

		case SLT: // 0x12
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigSigned(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			var word = common.ZeroWord256
			if x.Cmp(y) < 0 {
				word = common.OneWord256
			}
			if err = stack.Push(word); err != nil {
				return nil, err
			}

		case SGT: // 0x13
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigSigned(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			var word = common.ZeroWord256
			if x.Cmp(y) > 0 {
				word = common.OneWord256
			}
			if err = stack.Push(word); err != nil {
				return nil, err
			}

		case EQ: // 0x14
			var values []common.Word256
			var err error
			if values, err = stack.Pops(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			var word = common.ZeroWord256
			if bytes.Equal(x[:], y[:]) {
				word = common.OneWord256
			}
			if err = stack.Push(word); err != nil {
				return nil, err
			}

		case ISZERO: // 0x15
			x, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			var word = common.ZeroWord256
			if x.IsZero() {
				word = common.OneWord256
			}
			if err := stack.Push(word); err != nil {
				return nil, err
			}

		case AND: // 0x16
			var values []common.Word256
			var err error
			if values, err = stack.Pops(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			z := [32]byte{}
			for i := 0; i < 32; i++ {
				z[i] = x[i] & y[i]
			}
			if err = stack.Push(z); err != nil {
				return nil, err
			}

		case OR: // 0x17
			var values []common.Word256
			var err error
			if values, err = stack.Pops(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			z := [32]byte{}
			for i := 0; i < 32; i++ {
				z[i] = x[i] | y[i]
			}
			if err = stack.Push(z); err != nil {
				return nil, err
			}

		case XOR: // 0x18
			var values []common.Word256
			var err error
			if values, err = stack.Pops(2); err != nil {
				return nil, err
			}
			x, y := values[0], values[1]
			z := [32]byte{}
			for i := 0; i < 32; i++ {
				z[i] = x[i] ^ y[i]
			}
			if err = stack.Push(z); err != nil {
				return nil, err
			}

		case NOT: // 0x19
			x, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			z := [32]byte{}
			for i := 0; i < 32; i++ {
				z[i] = ^x[i]
			}
			if err := stack.Push(z); err != nil {
				return nil, err
			}

		case BYTE: // 0x1A
			idx, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			val, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			res := byte(0)
			if idx < 32 {
				res = val[idx]
			}
			if err := stack.Push64(int64(res)); err != nil {
				return nil, err
			}

		case SHL: //0x1B
			var err error
			var values []*big.Int
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			shift, x := values[0], values[1]

			if shift.Cmp(common.Big256) >= 0 {
				reset := big.NewInt(0)
				err = stack.PushBigInt(reset)
			} else {
				shiftedValue := x.Lsh(x, uint(shift.Uint64()))
				err = stack.PushBigInt(shiftedValue)
			}
			if err != nil {
				return nil, err
			}

		case SHR: //0x1C
			var err error
			var values []*big.Int
			if values, err = stack.PopsBigInt(2); err != nil {
				return nil, err
			}
			shift, x := values[0], values[1]

			if shift.Cmp(common.Big256) >= 0 {
				reset := big.NewInt(0)
				err = stack.PushBigInt(reset)
			} else {
				shiftedValue := x.Rsh(x, uint(shift.Uint64()))
				err = stack.PushBigInt(shiftedValue)
			}
			if err != nil {
				return nil, err
			}

		case SAR: //0x1D
			shift, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			x, err := stack.PopBigIntSigned()
			if err != nil {
				return nil, err
			}

			if shift.Cmp(common.Big256) >= 0 {
				reset := big.NewInt(0)
				if x.Sign() < 0 {
					reset.SetInt64(-1)
				}
				err = stack.PushBigInt(reset)
			} else {
				shiftedValue := x.Rsh(x, uint(shift.Uint64()))
				err = stack.PushBigInt(shiftedValue)
			}
			if err != nil {
				return nil, err
			}

		case SHA3: // 0x20
			var err error
			var values []*big.Int
			values, err = stack.PopsBigInt(2)
			offset, size := values[0], values[1]
			data, memErr := memory.Read(offset, size)
			if memErr != nil {
				return nil, err
			}
			data = sha3.Sha3(data)
			if err = stack.PushBytes(data); err != nil {
				return nil, err
			}

		case ADDRESS: // 0x30
			if err := stack.Push(callee.GetAddress().Word256()); err != nil {
				return nil, err
			}

		case BALANCE: // 0x31
			addr, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			acc, err := cache.GetAccount(common.AddressFromWord256(addr))
			if err != nil {
				return nil, err
			}
			if acc == nil {
				return nil, NewError(UnknownAddress)
			}
			balance := acc.GetBalance()
			if err = stack.PushU64(balance); err != nil {
				return nil, err
			}

		case ORIGIN: // 0x32
			// Origin is the origin of sender, the reason why the origin is needed is because
			// if a contract a call a contract b, the the caller that the b can see is a.
			// So this is the reason why the origin is needed.
			if err := stack.Push(evm.origin.Word256()); err != nil {
				return nil, err
			}

		case CALLER: // 0x33
			if err := stack.Push(caller.GetAddress().Word256()); err != nil {
				return nil, err
			}

		case CALLVALUE: // 0x34
			if err := stack.PushU64(value); err != nil {
				return nil, err
			}

		case CALLDATALOAD: // 0x35
			offset, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			data, ok := subslice(input, offset, 32)
			if !ok {
				return nil, NewError(InputOutOfBounds)
			}
			res := common.LeftPadWord256(data)
			if err = stack.Push(res); err != nil {
				return nil, err
			}

		case CALLDATASIZE: // 0x36
			if err := stack.Push64(int64(len(input))); err != nil {
				return nil, err
			}

		case CALLDATACOPY: // 0x37
			memOff, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			inputOff, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			length, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			data, ok := subslice(input, inputOff, length)
			if !ok {
				return nil, NewError(InputOutOfBounds)
			}
			err = memory.Write(memOff, data)
			if err != nil {
				return nil, err
			}

		case CODESIZE: // 0x38
			l := int64(len(code))
			if err := stack.Push64(l); err != nil {
				return nil, err
			}

		case CODECOPY: // 0x39
			memOff, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			codeOff, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			length, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			data, ok := subslice(code, codeOff, length)
			if !ok {
				return nil, NewError(CodeOutOfBounds)
			}
			err = memory.Write(memOff, data)
			if err != nil {
				return nil, err
			}

		case GASPRICE: // 0x3A
			// Now the price of gas is zero
			if err := stack.Push(common.ZeroWord256); err != nil {
				return nil, err
			}

		case EXTCODESIZE: // 0x3B
			addr, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			// if the account is not exist, so this should be an error?
			acc, err := cache.GetAccount(common.AddressFromWord256(addr))
			if err != nil {
				return nil, err
			}
			if acc == nil {
				// if _, ok := registeredNativeContracts[addr]; !ok {
				// 	return nil, firstErr(err, ErrorCodeUnknownAddress)
				// }
				err = stack.Push(common.ZeroWord256)
			} else {
				code := acc.GetCode()
				l := int64(len(code))
				err = stack.Push64(l)
			}
			if err != nil {
				return nil, err
			}

		case EXTCODECOPY: // 0x3C
			addr, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			acc, err := cache.GetAccount(common.AddressFromWord256(addr))
			if err != nil {
				return nil, err
			}
			// if acc == nil {
			// 	if _, ok := registeredNativeContracts[addr]; ok {
			// 		return nil, firstErr(err, ErrorCodeNativeContractCodeCopy)
			// 	}
			// 	return nil, firstErr(err, ErrorCodeUnknownAddress)
			// }
			code := acc.GetCode()
			memOff, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			codeOff, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			length, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			data, ok := subslice(code, codeOff, length)
			if !ok {
				return nil, NewError(CodeOutOfBounds)
			}
			err = memory.Write(memOff, data)
			if err != nil {
				return nil, err
			}

		case RETURNDATASIZE: // 0x3D
			if err := stack.Push64(int64(len(evm.returnData))); err != nil {
				return nil, err
			}

		case RETURNDATACOPY: // 0x3E
			var values []*big.Int
			var err error
			if values, err = stack.PopsBigInt(3); err != nil {
				return nil, err
			}
			memOff, outputOff, length := values[0], values[1], values[2]

			end := new(big.Int).Add(outputOff, length)

			if end.BitLen() > 64 || uint64(len(evm.returnData)) < end.Uint64() {
				return nil, NewError(ReturnDataOutOfBounds)
			}

			data := evm.returnData
			err = memory.Write(memOff, data)
			if err != nil {
				return nil, err
			}

		case BLOCKHASH: // 0x40
			if err := stack.Push(evm.context.BlockHash.Word256()); err != nil {
				return nil, err
			}

		case COINBASE: // 0x41
			if err := stack.Push(evm.context.CoinBase); err != nil {
				return nil, err
			}

		case TIMESTAMP: // 0x42
			if err := stack.Push64(int64(evm.context.BlockTime)); err != nil {
				return nil, err
			}

		case NUMBER: // 0x43
			if err := stack.PushU64(evm.context.Number); err != nil {
				return nil, err
			}

		case DIFFICULTY:
			if err := stack.PushU64(evm.context.Diffculty); err != nil {
				return nil, err
			}

		case GASLIMIT: // 0x45
			if err := stack.PushU64(evm.context.GasLimit); err != nil {
				return nil, err
			}

		case POP: // 0x50
			if _, err := stack.Pop(); err != nil {
				return nil, err
			}

		case MLOAD: // 0x51
			offset, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			data, err := memory.Read(offset, common.BigWord256Length)
			if err != nil {
				return nil, err
			}
			if err = stack.Push(common.LeftPadWord256(data)); err != nil {
				return nil, err
			}

		case MSTORE: // 0x52
			offset, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			data, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			err = memory.Write(offset, data.Bytes())
			if err != nil {
				return nil, err
			}

		case MSTORE8: // 0x53
			offset, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			val64, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			val := byte(val64 & 0xFF)
			err = memory.Write(offset, []byte{val})
			if err != nil {
				return nil, err
			}

		case SLOAD: // 0x54
			loc, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			data, err := cache.GetStorage(callee.GetAddress(), loc)
			if err != nil {
				// fmt.Println(errSto)
				return nil, err
			}
			if err = stack.Push(data); err != nil {
				return nil, err
			}

		case SSTORE: // 0x55
			var err error
			var values []common.Word256
			if values, err = stack.Pops(2); err != nil {
				return nil, err
			}
			loc, data := values[0], values[1]
			if err = cache.SetStorage(callee.GetAddress(), loc, data); err != nil {
				return nil, err
			}

		case JUMP: // 0x56
			to, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			err = evm.jump(code, to, &pc)
			if err != nil {
				return nil, err
			}
			continue

		case JUMPI: // 0x57
			pos, err := stack.Pop64()
			if err != nil {
				return nil, err
			}
			cond, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			if !cond.IsZero() {
				err = evm.jump(code, pos, &pc)
				if err != nil {
					return nil, err
				}
				continue
			}

		case PC: // 0x58
			if err := stack.Push64(pc); err != nil {
				return nil, err
			}

		case MSIZE: // 0x59
			// Note: Solidity will write to this offset expecting to find guaranteed
			// free memory to be allocated for it if a subsequent MSTORE is made to
			// this offset.
			capacity := memory.Capacity()
			if err := stack.PushBigInt(capacity); err != nil {
				return nil, err
			}

		case GAS: // 0x5A
			// stack.PushU64(*gas)
			if err := stack.PushU64(0); err != nil {
				return nil, err
			}

		case JUMPDEST: // 0x5B
			// Do nothing

		case PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8, PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16, PUSH17, PUSH18, PUSH19, PUSH20, PUSH21, PUSH22, PUSH23, PUSH24, PUSH25, PUSH26, PUSH27, PUSH28, PUSH29, PUSH30, PUSH31, PUSH32:
			a := int64(op - PUSH1 + 1)
			codeSegment, ok := subslice(code, pc+1, a)
			if !ok {
				return nil, NewError(CodeOutOfBounds)
			}
			res := common.LeftPadWord256(codeSegment)
			if err := stack.Push(res); err != nil {
				return nil, err
			}
			pc += a
			//

		case DUP1, DUP2, DUP3, DUP4, DUP5, DUP6, DUP7, DUP8, DUP9, DUP10, DUP11, DUP12, DUP13, DUP14, DUP15, DUP16:
			n := int(op - DUP1 + 1)
			if err := stack.Dup(n); err != nil {
				return nil, err
			}

		case SWAP1, SWAP2, SWAP3, SWAP4, SWAP5, SWAP6, SWAP7, SWAP8, SWAP9, SWAP10, SWAP11, SWAP12, SWAP13, SWAP14, SWAP15, SWAP16:
			n := int(op - SWAP1 + 2)
			if err := stack.Swap(n); err != nil {
				return nil, err
			}

		case LOG0, LOG1, LOG2, LOG3, LOG4:
			n := int(op - LOG0)
			topics := make([]common.Word256, n)
			offset, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			size, err := stack.PopBigInt()
			if err != nil {
				return nil, err
			}
			for i := 0; i < n; i++ {
				value, err := stack.Pop()
				if err != nil {
					return nil, nil
				}
				topics[i] = value
			}
			// data, memErr
			_, err = memory.Read(offset, size)
			if err != nil {
				return nil, err
			}
			// vm.eventSink.Log(&exec.LogEvent{
			// 	Address: callee.GetAddress(),
			// 	Topics:  topics,
			// 	Data:    data,
			// })

		case CREATE: // 0xF0
			evm.returnData = nil

			// if !HasPermission(callState, callee, permission.CreateContract) {
			// 	return nil, PermissionDenied{
			// 		Address: callee.GetAddress(),
			// 		Perm:    permission.CreateContract,
			// 	}
			// }
			contractValue, err := stack.PopU64()
			if err != nil {
				return nil, err
			}
			var values []*big.Int
			values, err = stack.PopsBigInt(2)
			offset, size := values[0], values[1]
			input, err := memory.Read(offset, size)
			if err != nil {
				return nil, err
			}

			newAccount, err := evm.createAccount(callee, input)
			if err != nil {
				return nil, err
			}

			// // Run the input to get the contract code.
			// // NOTE: no need to copy 'input' as per Call contract.
			ret, err := evm.Call(callee, newAccount, input, input, contractValue)
			if err != nil {
				stack.Push(common.ZeroWord256)
				// Note we both set the return buffer and return the result normally
				evm.returnData = ret
				// if callErr == ErrorCodeExecutionReverted {
				// 	return ret, callErr
				// }
				return ret, err
			}
			newAccount.SetCode(ret) // Set the code (ret need not be copied as per Call contract)
			if err = stack.Push(newAccount.GetAddress().Word256()); err != nil {
				return nil, err
			}

		case CALL, CALLCODE, DELEGATECALL: // 0xF1, 0xF2, 0xF4
			evm.returnData = nil
			_, err := stack.PopU64()
			if err != nil {
				return nil, err
			}
			addr, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			// NOTE: for DELEGATECALL value is preserved from the original
			// caller, as such it is not stored on stack as an argument
			// for DELEGATECALL and should not be popped.  Instead previous
			// caller value is used.  for CALL and CALLCODE value is stored
			// on stack and needs to be overwritten from the given value.
			if op != DELEGATECALL {
				value, err = stack.PopU64()
				if err != nil {
					return nil, err
				}
			}
			var values []*big.Int
			values, err = stack.PopsBigInt(3)
			inOffset, inSize, retOffset := values[0], values[1], values[2]
			retSize, err := stack.Pop64()
			if err != nil {
				return nil, err
			}

			// Get the arguments from the memory
			args, err := memory.Read(inOffset, inSize)
			if err != nil {
				return nil, err
			}

			// // Ensure that gasLimit is reasonable
			// if *gas < gasLimit {
			// 	// EIP150 - the 63/64 rule - rather than CodedError we pass this specified fraction of the total available gas
			// 	gasLimit = *gas - *gas/64
			// }
			// NOTE: we will return any used gas later.
			// *gas -= gasLimit

			// Begin execution
			var ret []byte
			var callErr error

			// if IsRegisteredNativeContract(addr) {
			// 	// Native contract
			// 	ret, callErr = ExecuteNativeContract(addr, callState, callee, args, &gasLimit, logger)
			// 	// for now we fire the Call event. maybe later we'll fire more particulars
			// 	// NOTE: these fire call go_events and not particular go_events for eg name reg or permissions
			// 	vm.fireCallEvent(&callErr, &ret, callee.GetAddress(), crypto.AddressFromWord256(addr), args, value, &gasLimit)
			// } else {
			// EVM contract
			// if useGasNegative(gas, GasGetAccount, &callErr) {
			// 	return nil, callErr
			// }
			acc, err := cache.GetAccount(common.AddressFromWord256(addr))
			if err != nil {
				return nil, err
			}
			// since CALL is used also for sending funds,
			// acc may not exist yet. This is an CodedError for
			// CALLCODE, but not for CALL, though I don't think
			// ethereum actually cares
			if op == CALLCODE {
				if acc == nil {
					return nil, NewError(UnknownAddress)
				}
				ret, callErr = evm.Call(callee, callee, acc.GetCode(), args, value)
			} else if op == DELEGATECALL {
				if acc == nil {
					return nil, NewError(UnknownAddress)
				}
				ret, callErr = evm.DelegateCall(caller, callee, acc.GetCode(), args, value)
			} else {
				// nil account means we're sending funds to a new account
				// if acc == nil {
				// 	if !HasPermission(callState, caller, permission.CreateAccount) {
				// 		return nil, PermissionDenied{
				// 			Address: callee.GetAddress(),
				// 			Perm:    permission.CreateAccount,
				// 		}
				// 	}
				// 	acc = acm.ConcreteAccount{Address: crypto.AddressFromWord256(addr)}.MutableAccount()
				// }
				// add account to the tx cache
				cache.SetAccount(acc)
				ret, callErr = evm.Call(callee, acc, acc.GetCode(), args, value)
			}
			// }
			evm.returnData = ret
			// In case any calls deeper in the stack (particularly SNatives) has altered either of two accounts to which
			// we hold a reference, we need to freshen our state for subsequent iterations of this call frame's EVM loop
			// var getErr error
			// caller, getErr = cache.GetAccount(caller.GetAddress())
			// if getErr != nil {
			// 	// fmt.Println(caller)
			// 	return nil, firstErr(err, ErrorCodeUnknownAddress)
			// }
			// callee, getErr = cache.GetAccount(callee.GetAddress())
			// if getErr != nil {
			// 	fmt.Println(2)
			// 	return nil, firstErr(err, ErrorCodeUnknownAddress)
			// }

			// Push result
			if callErr != nil {
				// So we can return nested CodedError if the top level return is an CodedError
				// vm.nestedCallErrors = append(vm.nestedCallErrors, NestedCall{
				// 	NestedError: callErr,
				// 	StackDepth:  vm.stackDepth,
				// 	Caller:      caller.GetAddress(),
				// 	Callee:      callee.GetAddress(),
				// })
				stack.Push(common.ZeroWord256)

				if callErr.Error() == "Execution reverted" {
					memory.Write(retOffset, common.RightPadBytes(ret, int(retSize)))
				}
			} else {
				stack.Push(common.OneWord256)

				// Should probably only be necessary when there is no return value and
				// ret is empty, but since EVM expects retSize to be respected this will
				// defensively pad or truncate the portion of ret to be returned.
				err = memory.Write(retOffset, common.RightPadBytes(ret, int(retSize)))
				if err != nil {
					return nil, err
				}
			}

			// Handle remaining gas.
			// *gas += gasLimit

		case RETURN: // 0xF3
			var err error
			var values []*big.Int
			values, err = stack.PopsBigInt(2)
			if err != nil {
				return nil, err
			}
			offset, size := values[0], values[1]
			output, err := memory.Read(offset, size)
			if err != nil {
				return nil, err
			}
			return output, nil

		case REVERT: // 0xFD
			var err error
			var values []*big.Int
			values, err = stack.PopsBigInt(2)
			offset, size := values[0], values[1]
			output, err := memory.Read(offset, size)
			if err != nil {
				return nil, err
			}

			return output, NewError(ExecutionReverted)

		case INVALID: //0xFE
			return nil, NewError(ExecutionAborted)

		case SELFDESTRUCT: // 0xFF
			addr, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			// if useGasNegative(gas, GasGetAccount, &err) {
			// 	return nil, err
			// }
			receiver, err := cache.GetAccount(common.AddressFromWord256(addr))
			if err != nil {
				return nil, err
			}
			if receiver == nil {
				// var gasErr ErrorCode
				// if useGasNegative(gas, GasCreateAccount, &gasErr) {
				// 	return nil, firstErr(err, gasErr)
				// }
				// if !HasPermission(callState, callee, permission.CreateContract) {
				// 	return nil, firstErr(err, PermissionDenied{
				// 		Address: callee.GetAddress(),
				// 		Perm:    permission.CreateContract,
				// 	})
				// }
				receiver, err = evm.createAccount(callee, receiver.GetCode())
				if err != nil {
					return nil, err
				}

			}

			err = receiver.AddBalance(callee.GetBalance())
			if err != nil {
				return nil, err
			}
			cache.SetAccount(receiver)
			cache.RemoveAccount(callee.GetAddress())
			// vm.Debugf(" => (%X) %v\n", addr[:4], callee.GetBalance())
			// fallthrough

		case STOP: // 0x00
			return nil, nil

		case STATICCALL, CREATE2:
			return nil, NewError(NotImplementation)

		default:
			return nil, NewError(UnknownError)
		}
		pc++
	}
}

// createAccount will cacaluate the address of the new contract account.
// First, check if the address exists before.
// Then, create a default account.
func (evm *EVM) createAccount(account common.Account, code []byte) (common.Account, error) {
	// fmt.Printf("callee address is %s\n", caller.GetAddress().String())
	var cache = evm.cache
	addr := common.NewContractAddress(evm.context.ChannelID, account.GetAddress(), code)
	if cache.AccountExist(addr) {
		return nil, NewError(DuplicateAddress)
	}
	newAccount := common.NewDefaultAccount(addr)
	err := cache.SetAccount(newAccount)
	if err != nil {
		return nil, err
	}
	return newAccount, nil
}

func (evm *EVM) jump(code []byte, to int64, pc *int64) (err error) {
	dest := codeGetOp(code, to)
	if dest != JUMPDEST {
		return NewError(InvalidJumpDest)
	}
	*pc = to
	return nil
}

func transfer(from, to common.Account, amount uint64) error {
	if from.GetBalance() < amount {
		return NewError(InsufficientBalance)
	}
	from.SubBalance(amount)
	err := to.AddBalance(amount)
	if err != nil {
		return NewError(UnknownError)
	}
	return nil
}

func codeGetOp(code []byte, n int64) OpCode {
	if int64(len(code)) <= n {
		return OpCode(0) // stop
	}
	return OpCode(code[n])
}

func subslice(data []byte, offset, length int64) (ret []byte, ok bool) {
	size := int64(len(data))
	if size < offset || offset < 0 || length < 0 {
		return nil, false
	} else if size < offset+length {
		ret, ok = data[offset:], true
		ret = common.RightPadBytes(ret, 32)
	} else {
		ret, ok = data[offset:offset+length], true
	}
	return
}
