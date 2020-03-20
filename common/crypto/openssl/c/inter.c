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

void print_it(const char *label, const unsigned char *buff, size_t len)
{
    if (!buff || !len)
        return;

    if (label)
        printf("%s %zu:\n", label, len);

    for (size_t i = 0; i < len; ++i)
        printf("%02x", buff[i]);

    printf("\n\n");
}
EVP_PKEY *ParsePrivateKey(char *data, int len)
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

EVP_PKEY *ParsePubKey(char *data, int len)
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

    if ((key = ParsePrivateKey(priv, len)) == NULL)
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

int Verify(char *pub, int len, unsigned char *hash, size_t hashlen, unsigned char *sig, size_t siglen)
{
    EVP_PKEY *key = NULL;
    EVP_PKEY_CTX *ctx = NULL;
    int ret = 0;

    if ((key = ParsePubKey(pub, len)) == NULL)
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

int MarshalPrivateKey(EVP_PKEY *privKey, char **priv, int *privlen)
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

int MarshalPublicKey(EVP_PKEY *pubKey, char **pub, int *publen)
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

    if ((key = ParsePrivateKey(priv, privlen)) == NULL)
    {
        goto err;
    }
    if ((ret = MarshalPublicKey(key, &pub, publen)) != 1)
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

    if ((ret = MarshalPrivateKey(evpriv, &priv, privlen)) != 1)
    {
        goto err;
    }
    if ((ret = MarshalPublicKey(evpriv, &pub, publen)) != 1)
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

int main()
{
    char *priv = OPENSSL_malloc(256);
    char *pub = OPENSSL_malloc(256);
    char *sig = OPENSSL_malloc(256);
    int privlen = 0;
    int publen = 0;
    int ret = 0;
    if ((ret = X_GenerateKey(priv, &privlen, pub, &publen)) != 1)
    {
        printf("Failed gen\n");
        goto err;
    }
    char *msg = "hello";
    size_t msglen = strlen(msg);
    size_t siglen = 0;
    if ((ret = Sign(priv, privlen, (unsigned char *)msg, msglen, (unsigned char *)sig, &siglen)) != 1)
    {
        printf("Failed sign\n");
        goto err;
    }
    if ((ret = Verify(pub, publen, (unsigned char *)msg, msglen, (unsigned char *)sig, siglen)) != 1)
    {
        printf("Failed berify\n");

        goto err;
    }
    print_it("priv", (const unsigned char *)priv, privlen);
    print_it("pub", (const unsigned char *)pub, publen);
    print_it("sig", (const unsigned char *)sig, siglen);
    goto done;
err:
    ERR_print_errors_fp(stderr);
done:
    OPENSSL_free(priv);
    OPENSSL_free(pub);
    OPENSSL_free(sig);

    return 0;
}