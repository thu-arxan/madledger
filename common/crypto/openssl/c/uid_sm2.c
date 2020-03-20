#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <openssl/bio.h>
#include <openssl/crypto.h>
#include <openssl/err.h>
#include <openssl/evp.h>
#include <openssl/ec.h>
#include <openssl/x509.h>
#include <openssl/pem.h>

void printerror() {
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

void print_it(const char* label,unsigned char * buff, size_t len)
{
    if(!buff || !len)
        return;
    
    if(label)
        printf("%s %zu:\n", label, len);
    
    for(size_t i=0; i < len; ++i)
        printf("%02x", buff[i]);
   
    printf("\n\n");
}

int testSm2() {
    int ret = 0;
    
    const char *privKey = 
    "-----BEGIN PRIVATE KEY-----\n"
    "MIGHAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBG0wawIBAQQg5g84gDqEk4c7lkl7\n"
    "BhHp7YU8E+cKsDK9uI+m7PEP7qihRANCAAS5DobGl+NnY7LnXvyCWbhNHUfa4nJo\n"
    "4780eAHRVbxyh92LxAaH9Vqio2vz2+kRK1B8rE+rV4Jrer1CuMnjIjTI\n"
    "-----END PRIVATE KEY-----\n";

    const char *pubKey =
    "-----BEGIN PUBLIC KEY-----\n"
    "MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEuQ6GxpfjZ2Oy5178glm4TR1H2uJy\n"
    "aOO/NHgB0VW8cofdi8QGh/VaoqNr89vpEStQfKxPq1eCa3q9QrjJ4yI0yA==\n"
    "-----END PUBLIC KEY-----\n";

    const char *msg = "hello";
    const char *uid = "xiaoxigua,dashagua";

    BIO *privBio = NULL;
    BIO *pubBio = NULL;

    EVP_PKEY *priv = NULL;
    EVP_PKEY *pub = NULL;

    EVP_MD_CTX *mctx = NULL;
    EVP_MD_CTX *mctx_verify = NULL;

    EVP_PKEY_CTX *skctx = NULL;
    EVP_PKEY_CTX *vkctx = NULL;

    unsigned char *sig = NULL;
    size_t sig_len = 0;

    privBio = BIO_new_mem_buf(privKey, strlen(privKey));
    if (privBio == NULL) {
        goto err;
    }
    priv = PEM_read_bio_PrivateKey(privBio, NULL, NULL, NULL);
    if (priv == NULL) {
        goto err;
    }
    ret = EVP_PKEY_set_alias_type(priv, EVP_PKEY_SM2);
    if (ret != 1) {
        goto err;
    }
    mctx = EVP_MD_CTX_new();
    if (mctx == NULL) {
        goto err;
    }
    skctx = EVP_PKEY_CTX_new(priv, NULL);
    if (skctx == NULL) {
        goto err;
    }
    ret = EVP_PKEY_CTX_set1_id(skctx, (const uint8_t *)uid, strlen(uid));
    if (ret <= 0) {
        goto err;
    }
    EVP_MD_CTX_set_pkey_ctx(mctx, skctx);

    /* Sign msg */
    ret = EVP_DigestSignInit(mctx, NULL, EVP_sm3(), NULL, priv);
    if (ret != 1) {
        goto err;
    }
    ret = EVP_DigestSignUpdate(mctx, msg, strlen(msg));
    if (ret != 1) {
        goto err;
    }
    /* Determine the size of the signature. */
    ret = EVP_DigestSignFinal(mctx, NULL, &sig_len);
    if (ret != 1) {
        goto err;
    }
    printf("Sig len %zu\n", sig_len);

    sig = OPENSSL_malloc(sig_len);
    if (sig == NULL) {
        goto err;
    }
    ret = EVP_DigestSignFinal(mctx, sig, &sig_len);
    if (ret != 1) {
        goto err;
    }

    print_it("sig", sig, sig_len);

    /* Ensure that the signature round-trips. */
    pubBio = BIO_new_mem_buf(pubKey, strlen(pubKey));
    if (pubBio == NULL) {
        goto err;
    }
    pub = PEM_read_bio_PUBKEY(pubBio, NULL, NULL, NULL);
    if (pub == NULL) {
        goto err;
    }
    ret = EVP_PKEY_set_alias_type(pub, EVP_PKEY_SM2);
    if (ret != 1) {
        goto err;
    }
    mctx_verify = EVP_MD_CTX_new();
    if (mctx == NULL) {
        goto err;
    }
    vkctx = EVP_PKEY_CTX_new(pub, NULL);
    if (vkctx == NULL) {
        goto err;
    }
    ret = EVP_PKEY_CTX_set1_id(vkctx, (const uint8_t *)uid, strlen(uid));
    if (ret <= 0) {
        goto err;
    }
    EVP_MD_CTX_set_pkey_ctx(mctx_verify, vkctx);

    /* Verify sig */
    ret = EVP_DigestVerifyInit(mctx_verify, NULL, EVP_sm3(), NULL, pub);
    if (ret != 1) {
        goto err;
    }
    ret = EVP_DigestVerifyUpdate(mctx_verify, msg, strlen(msg));
    if (ret != 1) {
        goto err;
    }
    ret = EVP_DigestVerifyFinal(mctx_verify, sig, sig_len);
    if (ret != 1) {
        goto err;
    }
    sig[10] = sig[10] + 1;
    ret = EVP_DigestVerifyFinal(mctx_verify, sig, sig_len);
    if (ret == 1) {
        goto err;
    }
    printf("verify succeed\n");
    ret = 1;
    goto done;
err:
    printf("failed\n");
    printerror();
    ret = -1;
    goto done;

done:
    BIO_free(privBio);
    BIO_free(pubBio);
    EVP_PKEY_CTX_free(skctx);
    EVP_PKEY_CTX_free(vkctx);
    EVP_MD_CTX_free(mctx);
    EVP_MD_CTX_free(mctx_verify);
    EVP_PKEY_free(priv);
    EVP_PKEY_free(pub);
    OPENSSL_free(sig);

    return ret;
}

int main() {
    int ret = testSm2();
    if (ret != 1) {
        printf("Error\n");
    }
    return 0;
}