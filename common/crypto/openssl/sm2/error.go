// +build cgo

package sm2

/*
#include <openssl/buffer.h>
#include <openssl/err.h>
#include <openssl/bio.h>
#include "util.h"
*/
import "C"
import (
	"errors"
	"io/ioutil"
)

// GetErrors get errors from openssl
func GetErrors() error {
	bio := C.BIO_new(C.BIO_s_mem())
	if bio == nil {
		return errors.New("GetErrors function failure: init BIO failed")
	}
	defer C.BIO_free(bio)
	C.ERR_print_errors(bio)
	data, err := ioutil.ReadAll(asAnyBio(bio))
	if err != nil {
		return errors.New("GetErrors Failed: " + err.Error())
	}
	return errors.New(string(data))
}

// OpenError get openssl error string
func OpenError() string {
	return GetErrors().Error()
}

// GetErrors get errors from openssl
// func GetErrors() error {
// 	var errs []string
// 	for {
// 		err := C.ERR_get_error()
// 		if err == 0 {
// 			break
// 		}
// 		errs = append(errs, fmt.Sprintf("%s:%s:%s",
// 			C.GoString(C.ERR_lib_error_string(err)),
// 			C.GoString(C.ERR_func_error_string(err)),
// 			C.GoString(C.ERR_reason_error_string(err))))
// 	}
// 	return fmt.Errorf("SSL errors: %s", strings.Join(errs, "\n"))
// }
