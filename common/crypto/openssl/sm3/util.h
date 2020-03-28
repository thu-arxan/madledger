#ifndef HEADER_UTIL_H
#define HEADER_UTIL_H
#include <openssl/bio.h>
long __BIO_get_mem_data(BIO *b, char **pp);
#endif