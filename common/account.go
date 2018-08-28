package common

// Account is the account of madledger
type Account interface {
	// Will be removed as soon as possible
	Address() Address
	Balance() uint64
	Code() []byte
	// These functions will be remained.
	GetAddress() Address
	GetBalance() uint64
	AddBalance(balance uint64) error
	SubBalance(balance uint64) error
	GetCode() []byte
	SetCode(code []byte)
	GetNonce() uint64
	SetNonce(nonce uint64)
}

// DefaultAccount is the default implementation of Account
type DefaultAccount struct {
	address Address
	balance uint64
	code    []byte
	nonce   uint64
}

// NewDefaultAccount is the constructor of DefaultAccount
func NewDefaultAccount(addr Address) *DefaultAccount {
	return &DefaultAccount{
		address: addr,
		balance: 0,
		code:    []byte{},
		nonce:   0,
	}
}

// Address is the implementation of Account
func (a *DefaultAccount) Address() Address {
	return a.address
}

// GetAddress is the implementation of Account
func (a *DefaultAccount) GetAddress() Address {
	return a.address
}

// Balance is the implementation of Account
func (a *DefaultAccount) Balance() uint64 {
	return a.balance
}

// GetBalance is the implementation of Account
func (a *DefaultAccount) GetBalance() uint64 {
	return a.balance
}

// AddBalance is the implementation of Account
// todo: should conside the attack
func (a *DefaultAccount) AddBalance(balance uint64) error {
	a.balance += balance
	return nil
}

// SubBalance is the implementation of Account
// todo: should conside the attack
func (a *DefaultAccount) SubBalance(balance uint64) error {
	a.balance -= balance
	return nil
}

// Code is the implementation of Account
func (a *DefaultAccount) Code() []byte {
	return a.code
}

// GetCode is the implementation of Account
func (a *DefaultAccount) GetCode() []byte {
	return a.code
}

// SetCode is the implementation of Account
func (a *DefaultAccount) SetCode(code []byte) {
	a.code = code
}

// GetNonce is the implementation of Account
func (a *DefaultAccount) GetNonce() uint64 {
	return a.nonce
}

// SetNonce is the implementation of Account
func (a *DefaultAccount) SetNonce(nonce uint64) {
	a.nonce = nonce
}
