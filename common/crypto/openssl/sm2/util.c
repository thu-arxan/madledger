#include <stdio.h>
#include "util.h"
#include <openssl/bio.h>
#include <openssl/ssl.h>
#include <openssl/conf.h>
#include <openssl/engine.h>
#include <openssl/err.h>

int init_openssl()
{
    int rc = 0;
    // OPENSSL_config(NULL);
    // ENGINE_load_builtin_engines();
    SSL_load_error_strings();
    SSL_library_init();
    ERR_load_ERR_strings();
    ERR_load_crypto_strings();
    OpenSSL_add_all_algorithms();
    //
    // Set up OPENSSL thread safety callbacks.
    rc = go_init_locks();
    if (rc != 0)
    {
        return rc;
    }
    CRYPTO_set_locking_callback(go_thread_locking_callback);
    CRYPTO_set_id_callback(go_thread_id_callback);

    return 0;
}

long _BIO_get_mem_data(BIO *b, char **pp)
{
    return BIO_get_mem_data(b, pp);
}

void X_OPENSSL_free(void *ref)
{
    OPENSSL_free(ref);
}

EVP_PKEY *parsePrivateKey(char *data, int len)
{
    BIO *bio = NULL;
    EVP_PKEY *key = NULL;
    if ((bio = BIO_new_mem_buf(data, len)) == NULL)
    {
        goto done;
    }
    if ((key = d2i_PrivateKey_bio(bio, NULL)) == NULL)
    {
        goto done;
    }
    goto done;
done:
    BIO_free(bio);
    return key;
}

EVP_PKEY *parsePubKey(char *data, int len)
{
    BIO *bio = NULL;
    EVP_PKEY *key = NULL;
    if ((bio = BIO_new_mem_buf(data, len)) == NULL)
    {
        goto done;
    }
    if ((key = d2i_PUBKEY_bio(bio, NULL)) == NULL)
    {
        goto done;
    }
    goto done;
done:
    BIO_free(bio);
    return key;
}

int Sign(char *priv, int len, unsigned char *hash, size_t hashlen, unsigned char *sig, size_t *siglen)
{
    EVP_PKEY *key = NULL;
    EVP_PKEY_CTX *ctx = NULL;
    int ret = 0;

    if ((key = parsePrivateKey(priv, len)) == NULL)
    {
        goto err;
    }

    if ((ret = EVP_PKEY_set_alias_type(key, EVP_PKEY_SM2)) != 1)
    {
        goto err;
    }
    if ((ctx = EVP_PKEY_CTX_new(key, NULL)) == NULL)
    {
        goto err;
    }
    if ((ret = EVP_PKEY_sign_init(ctx)) != 1)
    {
        goto err;
    }
    if ((ret = EVP_PKEY_CTX_set_signature_md(ctx, EVP_sm3())) != 1)
    {
        goto err;
    }
    if (EVP_PKEY_sign(ctx, NULL, siglen, hash, hashlen) <= 0)
    {
        goto err;
    }
    if (EVP_PKEY_sign(ctx, sig, siglen, hash, hashlen) <= 0)
    {
        goto err;
    }
    goto done;
err:
    ret = -1;
done:
    EVP_PKEY_free(key);
    EVP_PKEY_CTX_free(ctx);
    return ret;
}

int Sm2Sign(char *priv, int len, char *msg, size_t msglen, char *uid, size_t uidlen, unsigned char *sig, size_t *siglen)
{
    EVP_PKEY *key = NULL;
    EVP_MD_CTX *mctx = NULL;
    EVP_PKEY_CTX *ctx = NULL;
    int ret = 0;

    if ((key = parsePrivateKey(priv, len)) == NULL) {
        goto err;
    }

    if ((ret = EVP_PKEY_set_alias_type(key, EVP_PKEY_SM2)) != 1) {
        goto err;
    }

    if ((mctx = EVP_MD_CTX_new()) == NULL)  {
        goto err;
    }
    
    if ((ctx = EVP_PKEY_CTX_new(key, NULL)) == NULL) {
        goto err;
    }
    
    if ((ret = EVP_PKEY_CTX_set1_id(ctx, (const uint8_t *)uid, uidlen)) != 1) {
        goto err;
    }
    
    EVP_MD_CTX_set_pkey_ctx(mctx, ctx);

    if ((ret = EVP_DigestSignInit(mctx, NULL, EVP_sm3(), NULL, key)) != 1) {
        goto err;
    }

    if ((ret = EVP_DigestSignUpdate(mctx, msg, msglen)) != 1) {
        goto err;
    }

    /* Must call to get siglen or SM2err(SM2_F_PKEY_SM2_SIGN, SM2_R_BUFFER_TOO_SMALL); */
    if ((ret = EVP_DigestSignFinal(mctx, NULL, siglen)) != 1) {
        goto err;
    }

    if ((ret = EVP_DigestSignFinal(mctx, sig, siglen)) != 1) {
        goto err;
    }

    goto done;
err:
    ret = -1;
done:
    EVP_PKEY_free(key);
    EVP_MD_CTX_free(mctx);
    EVP_PKEY_CTX_free(ctx);
    return ret;
}

int Verify(char *pub, int len, unsigned char *hash, size_t hashlen, unsigned char *sig, size_t siglen)
{
    EVP_PKEY *key = NULL;
    EVP_PKEY_CTX *ctx = NULL;
    int ret = 0;

    if ((key = parsePubKey(pub, len)) == NULL)
    {
        goto err;
    }

    if ((ret = EVP_PKEY_set_alias_type(key, EVP_PKEY_SM2)) != 1)
    {
        goto err;
    }
    if ((ctx = EVP_PKEY_CTX_new(key, NULL)) == NULL)
    {
        goto err;
    }
    if ((ret = EVP_PKEY_verify_init(ctx)) != 1)
    {
        goto err;
    }
    if ((ret = EVP_PKEY_CTX_set_signature_md(ctx, EVP_sm3())) != 1)
    {
        goto err;
    }
    if (EVP_PKEY_verify(ctx, sig, siglen, hash, hashlen) != 1)
    {
        goto err;
    }
    goto done;
err:
    ret = -1;
done:
    EVP_PKEY_free(key);
    EVP_PKEY_CTX_free(ctx);
    return ret;
}

int Sm2Verify(char *pub, int len, char *msg, size_t msglen, char *uid, size_t uidlen, unsigned char *sig, size_t siglen) {
    EVP_PKEY *key = NULL;
    EVP_MD_CTX *mctx = NULL;
    EVP_PKEY_CTX *ctx = NULL;
    int ret = 0;

    if ((key = parsePubKey(pub, len)) == NULL) {
        goto err;
    }

    if ((ret = EVP_PKEY_set_alias_type(key, EVP_PKEY_SM2)) != 1) {
        goto err;
    }

    if ((mctx = EVP_MD_CTX_new()) == NULL)  {
        goto err;
    }
    
    if ((ctx = EVP_PKEY_CTX_new(key, NULL)) == NULL) {
        goto err;
    }
    
    if ((ret = EVP_PKEY_CTX_set1_id(ctx, (const uint8_t *)uid, uidlen)) != 1) {
        goto err;
    }
    
    EVP_MD_CTX_set_pkey_ctx(mctx, ctx);

    if ((ret = EVP_DigestVerifyInit(mctx, NULL, EVP_sm3(), NULL, key)) != 1) {
        goto err;
    }

    if ((ret = EVP_DigestVerifyUpdate(mctx, msg, msglen)) != 1) {
        goto err;
    }

    if ((ret = EVP_DigestVerifyFinal(mctx, sig, siglen)) != 1) {
        goto err;
    }

    goto done;
err:
    ret = -1;
done:
    EVP_PKEY_free(key);
    EVP_MD_CTX_free(mctx);
    EVP_PKEY_CTX_free(ctx);
    return ret;
}

int marshalPrivateKey(EVP_PKEY *privKey, char **priv, int *privlen)
{
    int ret = 1;
    BIO *bio = NULL;
    if ((bio = BIO_new(BIO_s_mem())) == NULL)
    {
        goto err;
    }
    if ((ret = i2d_PKCS8PrivateKey_bio(bio, privKey, NULL, NULL, 0, NULL, NULL)) != 1)
    {
        goto err;
    }
    if ((*privlen = BIO_ctrl_pending(bio)) <= 0)
    {
        goto err;
    }
    if ((*privlen = BIO_read(bio, *priv, *privlen)) <= 0)
    {
        goto err;
    }
    goto done;
err:
    ret = -1;
done:
    BIO_free(bio);
    return ret;
}

int marshalPublicKey(EVP_PKEY *pubKey, char **pub, int *publen)
{
    int ret = 1;
    BIO *bio = NULL;
    if ((bio = BIO_new(BIO_s_mem())) == NULL)
    {
        goto err;
    }
    if ((ret = i2d_PUBKEY_bio(bio, pubKey)) != 1)
    {
        goto err;
    }
    if ((*publen = BIO_ctrl_pending(bio)) <= 0)
    {
        goto err;
    }
    if ((*publen = BIO_read(bio, *pub, *publen)) <= 0)
    {
        goto err;
    }
    goto done;
err:
    ret = -1;
done:
    BIO_free(bio);
    return ret;
}

int GetPkFromPriv(char *priv, int privlen, char *pub, int *publen)
{
    EVP_PKEY *key = NULL;
    int ret = 0;

    if ((key = parsePrivateKey(priv, privlen)) == NULL)
    {
        goto err;
    }
    if ((ret = marshalPublicKey(key, &pub, publen)) != 1)
    {
        goto err;
    }
    goto done;
err:
    ret = -1;
done:
    EVP_PKEY_free(key);
    return ret;
}

int X_GenerateKey(char *priv, int *privlen, char *pub, int *publen)
{
    int ret = 0;
    EC_KEY *ecpriv = NULL;
    EVP_PKEY *evpriv = NULL;

    // generate sm2 priv Key
    ecpriv = EC_KEY_new_by_curve_name(NID_sm2);
    if (ecpriv == NULL)
    {
        goto err;
    }
    ret = EC_KEY_generate_key(ecpriv);
    if (ret != 1)
    {
        goto err;
    }

    // EVP_PKEY sm2 priv
    if ((evpriv = EVP_PKEY_new()) == NULL)
    {
        goto err;
    }
    if ((ret = EVP_PKEY_set1_EC_KEY(evpriv, ecpriv)) != 1)
    {
        goto err;
    }

    if ((ret = marshalPrivateKey(evpriv, &priv, privlen)) != 1)
    {
        goto err;
    }
    if ((ret = marshalPublicKey(evpriv, &pub, publen)) != 1)
    {
        goto err;
    }
    goto done;
err:
    ret = -1;
done:
    EC_KEY_free(ecpriv);
    EVP_PKEY_free(evpriv);
    return ret;
}