// +build cgo

package sm3

/*
#include <stdio.h>
#include <string.h>
#include <openssl/evp.h>
*/
import "C"

import (
	"unsafe"
)

/*
__owur int EVP_Digest(const void *data, size_t count,
                          unsigned char *md, unsigned int *size,
                          const EVP_MD *type, ENGINE *impl);
*/

// Sm3Sum hash data with sm3
func Sm3Sum(data []byte) []byte {
	var hash [C.EVP_MAX_MD_SIZE]byte
	var length int
	cs := C.CString(string(data))
	defer C.free(unsafe.Pointer(cs))
	C.EVP_Digest(unsafe.Pointer(cs), C.size_t(len(data)), (*C.uchar)(unsafe.Pointer(&hash[0])), (*C.uint)(unsafe.Pointer(&length)), C.EVP_sm3(), nil)
	return hash[:length]
}

// Sha256 hash data with sha256
func Sha256(data []byte) []byte {
	var hash [C.EVP_MAX_MD_SIZE]byte
	var length int
	cs := C.CString(string(data))
	defer C.free(unsafe.Pointer(cs))
	C.EVP_Digest(unsafe.Pointer(cs), C.size_t(len(data)), (*C.uchar)(unsafe.Pointer(&hash[0])), (*C.uint)(unsafe.Pointer(&length)), C.EVP_sha256(), nil)
	return hash[:length]
}
