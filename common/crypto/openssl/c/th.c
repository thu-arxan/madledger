// gcc -g -o main  convert.c -Wall -L/usr/local/openssl/lib -lcrypto -lssl -I/usr/local/openssl/include
#include <stdio.h>

#include <openssl/evp.h>
#include <openssl/err.h>
#include <openssl/pem.h>

int testpem() {
    OpenSSL_add_all_algorithms();
    ERR_load_crypto_strings();

    EVP_PKEY *priv = NULL;
    EC_KEY *ecpriv = NULL;
    EVP_PKEY *mem_pkey = NULL;
    EVP_PKEY * p2 = NULL;
    BIO *out = NULL;
    BIO *mem = NULL;
    BIO *mem2 = NULL;

    
    int ret = 0;
    // generate sm2 priv Key
    ecpriv = EC_KEY_new_by_curve_name(NID_sm2);
    if (ecpriv == NULL) {
        ret = 2;
        goto err;
    }
    ret = EC_KEY_generate_key(ecpriv);
    if (ret != 1) {
        ret = 3;
        goto err;
    }

    // EVP_PKEY sm2 priv, pub
    if((priv = EVP_PKEY_new()) == NULL) {
        ret = 4;
        goto err;
    }
    if ((ret = EVP_PKEY_set1_EC_KEY(priv, ecpriv)) != 1) {
        ret = 5;
        goto err;
    }

    /* Print out the parameters in PEM format */
    out = BIO_new_fp (stdout, BIO_NOCLOSE);
    if (out == NULL) {
        ret = -1;
        goto err;
    }

    /* Write the new DH private key in PEM format to stdout */
    if (PEM_write_bio_PKCS8PrivateKey(out, priv, 0, 0, 0, 0, 0) <= 0) {
        ret = 6;
        goto err;
    }
    if (PEM_write_bio_PrivateKey(out, priv, 0, 0, 0, 0, 0) <= 0) {
        ret = 6;
        goto err;
    }

    mem2 = BIO_new(BIO_s_mem());
    if (mem2 == NULL) {
        ret = -1;
        goto err;
    }
    if (PEM_write_bio_PKCS8PrivateKey(mem2, priv, 0, 0, 0, 0, 0) <= 0) {
        ret = -2;
        goto err;
    }
    if( (p2 = PEM_read_bio_PrivateKey(mem2, NULL, NULL, NULL)) ==NULL ) {
        ret = -3;
        goto err;
    }

    /* Write the new DH public key in PEM format to stdout */
    // if (PEM_write_bio_PUBKEY(out, pub) <= 0) {
    //     ERR_print_errors_fp(stderr);
    // }
    
    mem = BIO_new(BIO_s_mem());
    if (mem == NULL) {
        ret = -1;
        goto err;
    }

    /* Convert DH private key to DER format, no encryption/password possible */
    if (i2d_PrivateKey_bio(mem, priv) <= 0) {
        ret = 7;
        goto err;
    }

    /* Load DH private key from DER format, no encryption/password possible */
    if (d2i_PrivateKey_bio(mem, &mem_pkey) <= 0) {
       ret = 8;
       goto err;
    }
    ret = 1;
    goto done;
err:
    ERR_print_errors_fp(stderr);
done:
    EC_KEY_free(ecpriv);
    EVP_PKEY_free(priv);
    EVP_PKEY_free(mem_pkey);
    EVP_PKEY_free(p2);
    BIO_free_all(mem);
    BIO_free(out);
    return ret;
}

int main() {
    int ret = testpem();
    if (ret != 1) {
        printf("Failed");
    }
    return 0;
}