#include <stdio.h>
#include <openssl/rsa.h>
#include <openssl/pem.h>
#include <openssl/evp.h>
#include "util.h"

int main() {
    int ret = 0;
    RSA *r = NULL;
    BIGNUM *bne = NULL;
    BIO *bp_public = NULL, *bp_private = NULL, *bp_public_der = NULL;
    int bits = 2048;
    unsigned long e = RSA_F4;

    // Generate the RSA key
    printf("Generating RSA key...\n");
    bne = BN_new();
    ret = BN_set_word(bne, e);
    if(ret != 1) {
        goto free_all;
    }
    r = RSA_new();
    ret = RSA_generate_key_ex(r, bits, bne, NULL);
    if(ret != 1) {
        goto free_all;
    }

    // Save the public key in PEM format
    printf("Writing key files...\n");
    bp_public = BIO_new_file("public.pem", "w+");
    ret = PEM_write_bio_RSAPublicKey(bp_public, r);
    if(ret != 1) {
        goto free_all;
    }

    // Save the private key in PEM format
    bp_private = BIO_new_file("private.pem", "w+");
    ret = PEM_write_bio_RSAPrivateKey(bp_private, r, NULL, NULL, 0, NULL, NULL);

    // Save in DER
    EVP_PKEY *evp = EVP_PKEY_new();
    ret = EVP_PKEY_assign_RSA(evp, r);
    if(ret != 1){
        printf("failure %i\n", ret);
    }
    // bp_public_der = BIO_new_file("public.key", "w+");
    bp_public_der = BIO_new(BIO_s_mem());
    ret = i2d_PUBKEY_bio(bp_public_der, evp);
    BUF_MEM *bptr;
    BIO_get_mem_ptr(bp_public_der, &bptr);
    BIO_set_close(bp_public_der, BIO_NOCLOSE); /* So BIO_free() leaves BUF_MEM alone */

    print_it("pub:", bptr->data, bptr -> length);

    // Free everything
    free_all:
    BIO_free_all(bp_public);
    BIO_free_all(bp_public_der);
    BIO_free_all(bp_private);
    RSA_free(r);
    BN_free(bne);
    printf("Done!\n");

    return 0;
}