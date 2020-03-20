#include <stdio.h>
#include <string.h>
#include "util.h"
#include <openssl/evp.h>

#include <openssl/ec.h>
#include <openssl/x509.h>


int main() {
    EC_KEY    *eckey;
    eckey = EC_KEY_new_by_curve_name(NID_sm2);
    // EC_KEY_set_asn1_flag(eckey, OPENSSL_EC_NAMED_CURVE);

    if (eckey == NULL) {
        goto err;
    }
    // generate priv and pub key
    if (EC_KEY_generate_key(eckey) != 1) {
        goto err;
    }
    
    // encoding public key hex
    EC_POINT *ec_point = EC_KEY_get0_public_key(eckey);
    if (ec_point == NULL) {
        goto err;
    }
    EC_GROUP *ec_group = EC_KEY_get0_group(eckey);
    char * pubhex = EC_POINT_point2hex(ec_group, ec_point, POINT_CONVERSION_UNCOMPRESSED, NULL);
    printf("%lu HEX PUB:\n%s\n\n", strlen(pubhex), pubhex);

    // convert into evp_key
    EVP_PKEY *pkey = EVP_PKEY_new();
    EVP_PKEY_set1_EC_KEY(pkey, eckey);
    // EVP_PKEY_set_alias_type(pkey, EVP_PKEY_SM2);

    // encoding private key asn der
    int length = i2d_PrivateKey(pkey, NULL);
    if (length <= 0) {
        goto err;
    }
    unsigned char * p = OPENSSL_malloc(length);
    length = i2d_PrivateKey(pkey, &p);
    print_it("i2d_PrivateKey", p, length);
    printf("base58(i2d_PrivateKey):\n%s\n\n", Base64Encode(p, length, true));


    BIO *bp_priv_der = BIO_new(BIO_s_mem());
    i2d_PrivateKey_bio(bp_priv_der, pkey);
    BUF_MEM *privmem;
    BIO_get_mem_ptr(bp_priv_der, &privmem);
    BIO_set_close(bp_priv_der, BIO_NOCLOSE); /* So BIO_free() leaves BUF_MEM alone */

    print_it("i2d_PrivateKey_bio", privmem->data, privmem -> length);
    printf("base58(i2d_PrivateKey_bio):\n%s\n\n", Base64Encode(privmem->data, privmem -> length, true));

    // encode private key
    unsigned char * priv_str = NULL;
    int priv_length = i2d_ECPrivateKey(eckey, NULL);
    priv_str = OPENSSL_malloc(priv_length);
    priv_length = i2d_ECPrivateKey(eckey, &priv_str);
    print_it("i2d_ECPrivateKey", priv_str, priv_length);

    // encoding public key
    BIO *bp_public_der = BIO_new(BIO_s_mem());
    i2d_PUBKEY_bio(bp_public_der, pkey);
    
    BUF_MEM *bptr;
    BIO_get_mem_ptr(bp_public_der, &bptr);
    BIO_set_close(bp_public_der, BIO_NOCLOSE); /* So BIO_free() leaves BUF_MEM alone */

    print_it("i2d_PUBKEY_bio", bptr->data, bptr -> length);
    printf("base58(i2d_PUBKEY_bio):\n%s\n", Base64Encode(bptr->data, bptr -> length, true));

    FILE *fp = fopen("private.pem", "wb");
    if (fp == NULL) {
        goto err;
    }
    
    if (!PEM_write_PrivateKey(fp, pkey, NULL, NULL, 0, 0, NULL)) {
        /* Error */
        goto err;
    }
    FILE *fp2 = fopen("public.pem", "wb");
    if (fp == NULL) {
        goto err;
    }
    if (!PEM_write_PUBKEY(fp2, pkey)) {
        /* Error */
        goto err;
    }
    FILE *fp3 = fopen("privateec.pem", "wb");
    if (fp == NULL) {
        goto err;
    }
    if (!PEM_write_ECPrivateKey(fp3, eckey, NULL, NULL, 0, 0, NULL)) {
        /* Error */
        goto err;
    }
    BIO *pkc_bio = BIO_new(BIO_s_mem());
    // i2d_PrivateKey_bio(pkc_bio, pkey);
     
    if(!i2d_PKCS8PrivateKey_bio(pkc_bio, pkey, NULL, NULL, 0, 0, NULL)) {
        goto err;
    }
    BUF_MEM *pkcmem;
    BIO_get_mem_ptr(pkc_bio, &pkcmem);
    BIO_set_close(pkc_bio, BIO_NOCLOSE);
    print_it("i2d_PKCS8PrivateKey_bio", pkcmem->data, pkcmem -> length);
    printf("base58(i2d_PKCS8PrivateKey_bio):\n%s\n", Base64Encode(pkcmem->data, pkcmem -> length, true));
    FILE *fp4 = fopen("test.pem", "wb");
    BIO *pem_bio = BIO_new_fp(fp4, BIO_NOCLOSE);
    PEM_write_bio_PKCS8PrivateKey(pem_bio, pkey, NULL, NULL, 0, 0, NULL);
    
    return 0;
err:
    geterror();
    return -1;
}
