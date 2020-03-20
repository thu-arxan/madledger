/*
 * Brokeness exhibited by libcrypto when using the DER functions to
 * encode/decode a Diffie Hellman private key using the PKCS8 functions.
 *
 * ---
 *
 * libcrypto's d2i_PKCS8PrivateKey_bio function is unable to read data that is
 * generated using it's own i2d_PKCS8PrivateKey_bio function, this also limits
 * functionality tremendously in that i2d_PrivateKey can not take any callbacks
 * or information for a password, and thus can't encrypt the end result, which
 * the PKCS8 functions are able to do.
 *
 * The biggest issue is that when trying to write a library that wraps this
 * functionality, we don't know ahead of time whether we are going to be
 * working with a DH key or not ... so if we try to d2i an memory BIO we are
 * passed, if we use the PKCS8 functions we may fail. At that point part of the
 * data has already been read from the BIO, and for read/write BIO's we can't
 * reset it back to where it was, so we can't even attempt to fall back on the
 * d2i_PrivateKey function.
 *
 * The biggest problem is that the PEM functions work without issues. So
 * PEM_write_bio_PKCS8PrivateKey can correctly be decoded by
 * PEM_read_bio_PrivateKey (there is no PKCS8 equivelant, since it handles it
 * behind the scenes).
 *
 * ---
 *
 * Compile using:
 *
 * clang `pkg-config libcrypto --libs --cflags` -Wall -Wextra d2i_pkcs8privatekey.c
 *
 * or for broken version:
 *
 * clang `pkg-config libcrypto --libs --cflags` -Wall -Wextra d2i_pkcs8privatekey.c -DBROKEN
 *
 * gcc -g -o main  convert.c -Wall -L/usr/local/openssl/lib -lcrypto -lssl -I/usr/local/openssl/include -DBROKEN
 * Run:
 *
 * ./a.out
 * echo $?
 *
 * When compiled with BROKEN, return code will be 8.
 * When compiled without BROKEN, return code will be 0.
 */

#include <stdio.h>

#include <openssl/evp.h>
#include <openssl/err.h>
#include <openssl/dh.h>
#include <openssl/pem.h>

int main() {
    OpenSSL_add_all_algorithms();
    ERR_load_crypto_strings();

    EVP_PKEY_CTX *ctx;
    EVP_PKEY *pkey = 0;

    /* Generate new DH parameters */
    ctx = EVP_PKEY_CTX_new_id(EVP_PKEY_DH, 0);

    if (!ctx) {
        ERR_print_errors_fp(stderr);
        return 1;
    }

    if (EVP_PKEY_paramgen_init(ctx) <= 0) {
        ERR_print_errors_fp(stderr);
        return 2;
    }

    // DONT USE 256 IN PRODUCTION
    if (EVP_PKEY_CTX_set_dh_paramgen_prime_len(ctx, 256) <= 0) {
        ERR_print_errors_fp(stderr);
        return 3;
    }

    if (EVP_PKEY_paramgen(ctx, &pkey) <= 0) {
       ERR_print_errors_fp(stderr);
       return 4;
    }

    EVP_PKEY_CTX_free(ctx);

    /* Print out the parameters in PEM format */
    BIO *out = BIO_new_fp (stdout, BIO_NOCLOSE);
    PEM_write_bio_Parameters(out, pkey);

    EVP_PKEY *dh_pkey = 0;
    ctx = 0;

    /* Generate new DH private/public key */
    ctx = EVP_PKEY_CTX_new(pkey, 0);

    if (ctx == 0) {
       ERR_print_errors_fp(stderr);
       return 5;
    }

    if (EVP_PKEY_keygen_init(ctx) <= 0) {
        ERR_print_errors_fp(stderr);
        return 6;
    }

    if (EVP_PKEY_keygen(ctx, &dh_pkey) <=0) {
        ERR_print_errors_fp(stderr);
    }

    /* Write the new DH private key in PEM format to stdout */
    if (PEM_write_bio_PKCS8PrivateKey(out, dh_pkey, 0, 0, 0, 0, 0) <= 0) {
        ERR_print_errors_fp(stderr);
    }

    /* Write the new DH public key in PEM format to stdout */
    if (PEM_write_bio_PUBKEY(out, dh_pkey) <= 0) {
        ERR_print_errors_fp(stderr);
    }
    
    BIO *mem = BIO_new(BIO_s_mem());
    EVP_PKEY *mem_pkey = 0;

#ifdef BROKEN
    /* Convert DH private key to DER format using PKCS8 */
    if (i2d_PKCS8PrivateKey_bio(mem, dh_pkey, 0, 0, 0, 0, 0) <= 0) {
        ERR_print_errors_fp(stderr);
        return 7;
    }

    /*
     * Attempt to load the DH private key from DER format using PKCS8
     *
     * This is where it fails.
     */
    if (d2i_PrivateKey_bio(mem, &mem_pkey) <= 0) {
        ERR_print_errors_fp(stderr);
        return 8;
    }
    // if (d2i_PKCS8PrivateKey_bio(mem, &mem_pkey, 0, 0) <= 0) {
    //     ERR_print_errors_fp(stderr);
    //     return 8;
    // }
#else
    /* Convert DH private key to DER format, no encryption/password possible */
    if (i2d_PrivateKey_bio(mem, dh_pkey) <= 0) {
        ERR_print_errors_fp(stderr);
        return 9;
    }

    /* Load DH private key from DER format, no encryption/password possible */
    if (d2i_PrivateKey_bio(mem, &mem_pkey) <= 0) {
        ERR_print_errors_fp(stderr);
        return 10;
    }
#endif
    
    /* Clean up :-) */
    EVP_PKEY_free(mem_pkey);
    BIO_free_all(mem);

    EVP_PKEY_free(pkey);
    EVP_PKEY_CTX_free(ctx);
    EVP_PKEY_free(dh_pkey);
    BIO_free_all(out);

    return 0;
}