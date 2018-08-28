package evm

import (
	"fmt"
	"madledger/common"

	"math/big"
)

// Stack is the stack
// It's not goroutine safe
type Stack struct {
	data []common.Word256
	ptr  int

	// gas is not used now
	gas *uint64
}

// NewStack is the constructor of Stack
func NewStack(capacity int) *Stack {
	return &Stack{
		data: make([]common.Word256, capacity),
		ptr:  0,
		// gas:  gas,
	}
}

func (st *Stack) useGas(gasToUse uint64) {
	// if *st.gas > gasToUse {
	// 	*st.gas -= gasToUse
	// } else {
	// 	st.setErr(ErrorCodeInsufficientGas)
	// }
}

// Push push a word
func (st *Stack) Push(d common.Word256) error {
	// st.useGas(GasStackOp)
	if st.ptr == cap(st.data) {
		return NewError(DataStackOverflow)
	}
	st.data[st.ptr] = d
	st.ptr++
	return nil
}

// PushBytes currently only called after sha3.Sha3
func (st *Stack) PushBytes(bz []byte) error {
	if len(bz) != 32 {
		return NewError(InvalidByteLength)
	}
	return st.Push(common.LeftPadWord256(bz))
}

// Push64 push an int64
func (st *Stack) Push64(i int64) error {
	return st.Push(common.Int64ToWord256(i))
}

// PushU64 push an uint64
func (st *Stack) PushU64(i uint64) error {
	return st.Push(common.Uint64ToWord256(i))
}

// PushBigInt pushes the bigInt as a common.Word256 encoding negative values in 32-byte twos complement and returns the encoded result
func (st *Stack) PushBigInt(bigInt *big.Int) error {
	word := common.LeftPadWord256(common.U256(bigInt).Bytes())
	return st.Push(word)
}

// Pop pops a word
func (st *Stack) Pop() (common.Word256, error) {
	// st.useGas(GasStackOp)
	if st.ptr == 0 {
		return common.ZeroWord256, NewError(DataStackUnderflow)
	}
	st.ptr--
	return st.data[st.ptr], nil
}

// Pops pops slice of  words which length is size
func (st *Stack) Pops(size int) ([]common.Word256, error) {
	var words []common.Word256
	for i := 0; i < size; i++ {
		word, err := st.Pop()
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

// PopBytes pops bytes
func (st *Stack) PopBytes() ([]byte, error) {
	word, err := st.Pop()
	if err != nil {
		return nil, err
	}
	return word.Bytes(), nil
}

// PopsBytes pops slice of bytes which length is size
func (st *Stack) PopsBytes(size int) ([][]byte, error) {
	var bs [][]byte
	for i := 0; i < size; i++ {
		b, err := st.PopBytes()
		if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}

// Pop64 pop int64
func (st *Stack) Pop64() (int64, error) {
	d, err := st.Pop()
	if err != nil {
		return 0, err
	}
	if d.Is64BitOverflow() {
		return 0, NewError(CallStackOverflow)
	}
	return common.Int64FromWord256(d), nil
}

// Pops64 pops slice of int64 which length is size
func (st *Stack) Pops64(size int) ([]int64, error) {
	var result []int64
	for i := 0; i < size; i++ {
		value, err := st.Pop64()
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

// PopU64 pop uint64
func (st *Stack) PopU64() (uint64, error) {
	d, err := st.Pop()
	if err != nil {
		return 0, err
	}
	if d.Is64BitOverflow() {
		return 0, NewError(CallStackOverflow)
	}
	return common.Uint64FromWord256(d), nil
}

// PopBigIntSigned pop a S256
func (st *Stack) PopBigIntSigned() (*big.Int, error) {
	value, err := st.PopBigInt()
	if err != nil {
		return nil, err
	}
	return common.S256(value), nil
}

// PopsBigSigned pops slice of big.Int which length is size
func (st *Stack) PopsBigSigned(size int) ([]*big.Int, error) {
	var result []*big.Int
	for i := 0; i < size; i++ {
		value, err := st.PopBigIntSigned()
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

// PopBigInt pop bit.Int
func (st *Stack) PopBigInt() (*big.Int, error) {
	d, err := st.Pop()
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(d[:]), nil
}

// PopsBigInt pops slice of big.Int which length is size
func (st *Stack) PopsBigInt(size int) ([]*big.Int, error) {
	var result []*big.Int
	for i := 0; i < size; i++ {
		value, err := st.PopBigInt()
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

// Len return the ptr
func (st *Stack) Len() int {
	return st.ptr
}

// Swap is the swap func
func (st *Stack) Swap(n int) error {
	// st.useGas(GasStackOp)
	if st.ptr < n {
		return NewError(DataStackUnderflow)
	}
	st.data[st.ptr-n], st.data[st.ptr-1] = st.data[st.ptr-1], st.data[st.ptr-n]
	return nil
}

// Dup is the dup func
func (st *Stack) Dup(n int) error {
	// st.useGas(GasStackOp)
	if st.ptr < n {
		return NewError(DataStackUnderflow)
	}
	return st.Push(st.data[st.ptr-n])
}

// Peek is not an opcode, costs no gas.
func (st *Stack) Peek() (common.Word256, error) {
	if st.ptr == 0 {
		return common.ZeroWord256, NewError(DataStackUnderflow)
	}
	return st.data[st.ptr-1], nil
}

// Print print the stack
func (st *Stack) Print(n int) {
	fmt.Println("### stack ###")
	if st.ptr > 0 {
		nn := n
		if st.ptr < n {
			nn = st.ptr
		}
		for j, i := 0, st.ptr-1; i > st.ptr-1-nn; i-- {
			fmt.Printf("%-3d  %X\n", j, st.data[i])
			j++
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
