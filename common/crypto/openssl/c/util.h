#ifndef HEADER_UTIL_H
#define HEADER_UTIL_H
#include <stdio.h>
#include <stdbool.h>
#include <openssl/ec.h>
void geterror();
void print_it(const char* label, const unsigned char* buff, size_t len);
int ECDSA_SIG_Encode(const ECDSA_SIG *sig, unsigned char **bytes);
ECDSA_SIG * ECDSA_SIG_Decode(const unsigned char **pp, long len);
char * Base64Encode(const char * input, int length, bool with_new_line);
char * Base64Decode(char * input, int length, bool with_new_line);
#endif