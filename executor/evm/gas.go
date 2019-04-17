package evm

/*
* This file defines the cost of each kind of subsets of instructions.
 */

type gasCost uint64

// selfdestruct is not included, because it has two values {24000, 5000}
const (
	// gasZero includes {STOP, RETURN, REVERT}
	gasZero gasCost = 0
	// gasBase includes {ADDRESS, ORIGIN ,CALLER, CALLVALUE, CALLDATASIZE, CODESIZE, GASPRICE,
	// COINBASE, TIMESTAMP, NUMBER, DIFFICULTY, GASLIMIT, RETURNDATASIZE, POP, PC, MSIZE, GAS}
	gasBase gasCost = 2
	// gasVerylow includes {ADD, SUB, NOT, LT, GT, SLT, SGT, EQ, ISZERO, AND, OR, XOR, BYTE,
	// CALLDATALOAD, MLOAD, MSTORE, MSTORE8, PUSH*, DUP*, SWAP*}
	gasVerylow gasCost = 3
	// gasLow includes {MUL, DIV,SDIV, MOD, SMOD, SIGNEXTEND}
	gasLow gasCost = 5
	// gasMid includes {ADDMOD, MULMOD, JUMP}
	gasMid gasCost = 8
	// gasHigh includes {JUMPI}
	gasHigh gasCost = 10
	// gasExtcode includes {EXTCODESIZE}
	gasExtcode       gasCost = 700
	gasBalance       gasCost = 400
	gasSload         gasCost = 200
	gasJumpdest      gasCost = 1
	gasSset          gasCost = 20000
	gasSreset        gasCost = 5000
	gasSclear        gasCost = 15000
	gasCreate        gasCost = 32000
	gasCodedeposit   gasCost = 200
	gasCall          gasCost = 700
	gasCallvalue     gasCost = 9000
	gasCallstipend   gasCost = 2300
	gasNewaccount    gasCost = 25000
	gasExp           gasCost = 10
	gasExpbyte       gasCost = 50
	gasMemory        gasCost = 3
	gasTxcreate      gasCost = 32000
	gasTxdatazero    gasCost = 4
	gasTxdatanonzero gasCost = 68
	gasTransaction   gasCost = 21000
	gasLog           gasCost = 375
	gasLogdata       gasCost = 8
	gasLogtopic      gasCost = 375
	gasSha3          gasCost = 30
	gasSha3word      gasCost = 6
	gasCopy          gasCost = 3
	gasBlockhash     gasCost = 20
	gasQuaddivisor   gasCost = 100
)

// GasCosts is the map contains the cost of some instructions
// However, not all instructions are included, because some instruction
// will cost different gas in different situtation.
var GasCosts = map[OpCode]gasCost{
	// gasZero
	STOP:   gasZero,
	RETURN: gasZero,
	REVERT: gasZero,
	// gasBase
	ADDRESS:        gasBase,
	ORIGIN:         gasBase,
	CALLER:         gasBase,
	CALLVALUE:      gasBase,
	CALLDATASIZE:   gasBase,
	CODESIZE:       gasBase,
	GASPRICE:       gasBase,
	COINBASE:       gasBase,
	TIMESTAMP:      gasBase,
	NUMBER:         gasBase,
	DIFFICULTY:     gasBase,
	GASLIMIT:       gasBase,
	RETURNDATASIZE: gasBase,
	POP:            gasBase,
	PC:             gasBase,
	MSIZE:          gasBase,
	GAS:            gasBase,
	// gasVerylow
	ADD:          gasVerylow,
	SUB:          gasVerylow,
	NOT:          gasVerylow,
	LT:           gasVerylow,
	GT:           gasVerylow,
	SLT:          gasVerylow,
	SGT:          gasVerylow,
	EQ:           gasVerylow,
	ISZERO:       gasVerylow,
	AND:          gasVerylow,
	OR:           gasVerylow,
	XOR:          gasVerylow,
	BYTE:         gasVerylow,
	CALLDATALOAD: gasVerylow,
	MLOAD:        gasVerylow,
	MSTORE:       gasVerylow,
	MSTORE8:      gasVerylow,
	PUSH1:        gasVerylow,
	PUSH2:        gasVerylow,
	PUSH3:        gasVerylow,
	PUSH4:        gasVerylow,
	PUSH5:        gasVerylow,
	PUSH6:        gasVerylow,
	PUSH7:        gasVerylow,
	PUSH8:        gasVerylow,
	PUSH9:        gasVerylow,
	PUSH10:       gasVerylow,
	PUSH11:       gasVerylow,
	PUSH12:       gasVerylow,
	PUSH13:       gasVerylow,
	PUSH14:       gasVerylow,
	PUSH15:       gasVerylow,
	PUSH16:       gasVerylow,
	PUSH17:       gasVerylow,
	PUSH18:       gasVerylow,
	PUSH19:       gasVerylow,
	PUSH20:       gasVerylow,
	PUSH21:       gasVerylow,
	PUSH22:       gasVerylow,
	PUSH23:       gasVerylow,
	PUSH24:       gasVerylow,
	PUSH25:       gasVerylow,
	PUSH26:       gasVerylow,
	PUSH27:       gasVerylow,
	PUSH28:       gasVerylow,
	PUSH29:       gasVerylow,
	PUSH30:       gasVerylow,
	PUSH31:       gasVerylow,
	PUSH32:       gasVerylow,
	DUP1:         gasVerylow,
	DUP2:         gasVerylow,
	DUP3:         gasVerylow,
	DUP4:         gasVerylow,
	DUP5:         gasVerylow,
	DUP6:         gasVerylow,
	DUP7:         gasVerylow,
	DUP8:         gasVerylow,
	DUP9:         gasVerylow,
	DUP10:        gasVerylow,
	DUP11:        gasVerylow,
	DUP12:        gasVerylow,
	DUP13:        gasVerylow,
	DUP14:        gasVerylow,
	DUP15:        gasVerylow,
	DUP16:        gasVerylow,
	SWAP1:        gasVerylow,
	SWAP2:        gasVerylow,
	SWAP3:        gasVerylow,
	SWAP4:        gasVerylow,
	SWAP5:        gasVerylow,
	SWAP6:        gasVerylow,
	SWAP7:        gasVerylow,
	SWAP8:        gasVerylow,
	SWAP9:        gasVerylow,
	SWAP10:       gasVerylow,
	SWAP11:       gasVerylow,
	SWAP12:       gasVerylow,
	SWAP13:       gasVerylow,
	SWAP14:       gasVerylow,
	SWAP15:       gasVerylow,
	SWAP16:       gasVerylow,
	// gasLow
	MUL:        gasLow,
	DIV:        gasLow,
	SDIV:       gasLow,
	MOD:        gasLow,
	SMOD:       gasLow,
	SIGNEXTEND: gasLow,
	// gasMid
	ADDMOD: gasMid,
	MULMOD: gasMid,
	JUMP:   gasMid,
	// gasHigh
	JUMPI: gasHigh,
	// gasExtcode
	EXTCODESIZE: gasExtcode,
}
