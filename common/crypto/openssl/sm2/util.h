#ifndef HEADER_UTIL_H
#define HEADER_UTIL_H
#include <stdio.h>
#include <stdlib.h>
#include <openssl/bio.h>
#include <openssl/pem.h>
#include <openssl/evp.h>
#include <openssl/ec.h>
#include <openssl/x509.h>
extern int go_init_locks();
extern void go_thread_locking_callback(int, int, const char *, int);
extern unsigned long go_thread_id_callback();
int init_openssl();
int _BIO_set_close(BIO *b, long c);
long _BIO_get_mem_data(BIO *b, char **pp);
EVP_PKEY *parsePrivateKey(char *data, int len);
EVP_PKEY *parsePubKey(char *data, int len);
int Sm2Sign(char *priv, int len, char *msg, size_t msglen, char *uid, size_t uidlen, unsigned char *sig, size_t *siglen);
int Sign(char *priv, int len, unsigned char *hash, size_t hashlen, unsigned char *sig, size_t *siglen);
int Sm2Verify(char *pub, int len, char *msg, size_t msglen, char *uid, size_t uidlen, unsigned char *sig, size_t siglen);
int Verify(char *pub, int len, unsigned char *hash, size_t hashlen, unsigned char *sig, size_t siglen);
int marshalPrivateKey(EVP_PKEY *privKey, char **priv, int *privlen);
int marshalPublicKey(EVP_PKEY *pubKey, char **pub, int *publen);
int GetPkFromPriv(char *priv, int privlen, char *pub, int *publen);
int X_GenerateKey(char *priv, int *privlen, char *pub, int *publen);
#endif