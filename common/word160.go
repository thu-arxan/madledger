package common

const (
	// Word160Length define the length of Word160
	Word160Length = 20
	// Word256Word160Delta define the delta between the Word160 and Word256
	Word256Word160Delta = 12
)

// ZeroWord160 is the zero of Word160
var ZeroWord160 = Word160{}

// Word160 is the bytes which length is 20
type Word160 [Word160Length]byte

// Word256 return the Word256 of a Word160
func (w Word160) Word256() (word256 Word256) {
	copy(word256[Word256Word160Delta:], w[:])
	return
}

// Bytes return the bytes of Word160
func (w Word160) Bytes() []byte {
	return w[:]
}
