#include <stdio.h>
#include "util.h"
#include <openssl/evp.h>
#include <openssl/ec.h>
#include <openssl/bio.h>
#include <openssl/err.h>

#include <openssl/buffer.h>

void geterror() {
    BIO * bio= BIO_new(BIO_s_mem());
    if (bio == NULL) {
        printf("Failed to get error\n");
        goto done;
    }
    ERR_print_errors(bio);
    char * err;
    BIO_get_mem_data(bio, &err);
    printf("err: %s\n", err);
done:
    BIO_free(bio);
}

void print_it(const char* label, const unsigned char* buff, size_t len)
{
    if(!buff || !len)
        return;
    
    if(label)
        printf("%s %zu:\n", label, len);
    
    for(size_t i=0; i < len; ++i)
        printf("%02x", buff[i]);
   
    printf("\n\n");
}

int ECDSA_SIG_Encode(const ECDSA_SIG *sig, unsigned char **bytes) {
    return i2d_ECDSA_SIG(sig, bytes);
}

ECDSA_SIG *ECDSA_SIG_Decode(const unsigned char **pp, long len) {
    return d2i_ECDSA_SIG(NULL, pp, len);
}

char * Base64Encode(const char * input, int length, bool with_new_line)
{
	BIO * bmem = NULL;
	BIO * b64 = NULL;
	BUF_MEM * bptr = NULL;
 
	b64 = BIO_new(BIO_f_base64());
	if(!with_new_line) {
		BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
	}
	bmem = BIO_new(BIO_s_mem());
	b64 = BIO_push(b64, bmem);
	BIO_write(b64, input, length);
	BIO_flush(b64);
	BIO_get_mem_ptr(b64, &bptr);
 
	char * buff = (char *)malloc(bptr->length + 1);
	memcpy(buff, bptr->data, bptr->length);
	buff[bptr->length] = 0;
 
	BIO_free_all(b64);
 
	return buff;
}
 
char * Base64Decode(char * input, int length, bool with_new_line)
{
	BIO * b64 = NULL;
	BIO * bmem = NULL;
	char * buffer = (char *)malloc(length);
	memset(buffer, 0, length);
 
	b64 = BIO_new(BIO_f_base64());
	if(!with_new_line) {
		BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
	}
	bmem = BIO_new_mem_buf(input, length);
	bmem = BIO_push(b64, bmem);
	BIO_read(bmem, buffer, length);
 
	BIO_free_all(bmem);
	return buffer;
}