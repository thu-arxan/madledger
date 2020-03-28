package sm2

/*
#cgo linux CFLAGS: -I/usr/local/openssl/include
#cgo linux LDFLAGS: -L/usr/local/openssl/lib -lcrypto -lssl
#cgo darwin CFLAGS: -I/usr/local/openssl/include
#cgo darwin LDFLAGS: -L/usr/local/openssl/lib -lcrypto -lssl

#include <openssl/bio.h>
#include <openssl/crypto.h>
#include <openssl/evp.h>
#include <openssl/ec.h>
#include <openssl/buffer.h>

*/
import "C"
