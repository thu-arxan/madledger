// +build cgo

package sm2

/*
#include "util.h"
*/
import "C"

func init() {
	ret := C.init_openssl()
	if ret != 0 {
		panic("Failed to init openssl")
	}
}
