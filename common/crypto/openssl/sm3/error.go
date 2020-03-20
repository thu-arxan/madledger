// +build cgo

package sm3

/*
#include <openssl/buffer.h>
#include <openssl/err.h>
#include <openssl/bio.h>
#include "util.h"
*/
import "C"
import "errors"

// GetErrors get errors from openssl
func GetErrors() error {
	bio := C.BIO_new(C.BIO_s_mem())
	if bio == nil {
		return errors.New("GetErrors function failure: init BIO failed")
	}
	defer C.BIO_free(bio)
	C.ERR_print_errors(bio)
	var p *C.char
	len := C.__BIO_get_mem_data(bio, &p)
	if len <= 0 {
		return errors.New("GetErrors function failure: get mem data failed")
	}
	return errors.New(C.GoString(p))
}
