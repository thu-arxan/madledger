#include <stdio.h>
#include "util.h"
#include <openssl/bio.h>

long __BIO_get_mem_data(BIO *b, char **pp) {
	return BIO_get_mem_data(b, pp);
}

