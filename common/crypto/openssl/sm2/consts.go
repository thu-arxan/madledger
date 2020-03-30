package sm2

const (
	// DefaultUID The default user id as specified in GM/T 0009-2012
	DefaultUID = "1234567812345678"
)
const (
	pubkeyCompressed   byte = 0x2  // y_bit + x coord
	pubkeyASN1         byte = 0x30 // asn1
	pubkeyUncompressed byte = 0x4  // x coord + y coord
)

// These constants define the lengths of serialized public keys.
const (
	pubKeyBytesLenCompressed   = 33
	pubKeyBytesLenASN1         = 91
	pubKeyBytesLenUncompressed = 65
)

// These constants define the lengths of serialized signature. 70-72
const (
	minSigLen = 64
	maxSigLen = 72
)
