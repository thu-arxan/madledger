package evm

import "fmt"

// ErrorCode is uint32 to represent the Error code
type ErrorCode uint32

// Defines different kind of errors
const (
	UnknownAddress ErrorCode = iota
	InsufficientBalance
	InvalidJumpDest
	InsufficientGas
	MemoryOutOfBounds
	CodeOutOfBounds
	InputOutOfBounds
	ReturnDataOutOfBounds
	CallStackOverflow
	CallStackUnderflow
	DataStackOverflow
	DataStackUnderflow
	InvalidContract
	NativeContractCodeCopy
	ExecutionAborted
	ExecutionReverted
	PermissionDenied
	NativeFunction
	EventPublish
	InvalidString
	EventMapping
	InvalidAddress
	DuplicateAddress
	InsufficientFunds
	Overpayment
	ZeroPayment
	InvalidSequence
	NotImplementation
	InvalidByteLength
	UnknownError
)

func (c ErrorCode) Error() string {
	return fmt.Sprintf("Error %d: %s", c, c.String())
}

func (c ErrorCode) String() string {
	switch c {
	case UnknownAddress:
		return "Unknown address"
	case InsufficientBalance:
		return "Insufficient balance"
	case InvalidJumpDest:
		return "Invalid jump dest"
	case InsufficientGas:
		return "Insufficient gas"
	case MemoryOutOfBounds:
		return "Memory out of bounds"
	case CodeOutOfBounds:
		return "Code out of bounds"
	case InputOutOfBounds:
		return "Input out of bounds"
	case ReturnDataOutOfBounds:
		return "Return data out of bounds"
	case CallStackOverflow:
		return "Call stack overflow"
	case CallStackUnderflow:
		return "Call stack underflow"
	case DataStackOverflow:
		return "Data stack overflow"
	case DataStackUnderflow:
		return "Data stack underflow"
	case InvalidContract:
		return "Invalid contract"
	case PermissionDenied:
		return "Permission denied"
	case NativeContractCodeCopy:
		return "Tried to copy native contract code"
	case ExecutionAborted:
		return "Execution aborted"
	case ExecutionReverted:
		return "Execution reverted"
	case NativeFunction:
		return "Native function error"
	case EventPublish:
		return "Event publish error"
	case InvalidString:
		return "Invalid string"
	case EventMapping:
		return "Event mapping error"
	case InvalidAddress:
		return "Invalid address"
	case DuplicateAddress:
		return "Duplicate address"
	case InsufficientFunds:
		return "Insufficient funds"
	case Overpayment:
		return "Overpayment"
	case ZeroPayment:
		return "Zero payment error"
	case InvalidSequence:
		return "Invalid sequence number"
	case NotImplementation:
		return "Not implementation yet"
	case InvalidByteLength:
		return "Invalid length of bytes"
	default:
		return "Unknown error"
	}
}

// Error is the error of EVM
type Error struct {
	code ErrorCode
}

// NewError is the constructor of Error
func NewError(code ErrorCode) error {
	return &Error{
		code: code,
	}
}

func (e *Error) Error() string {
	return e.code.String()
}
