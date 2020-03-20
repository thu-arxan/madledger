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

void print_it(const char* label, unsigned char* buff, size_t len)
{
    if(!buff || !len)
        return;
    
    if(label)
        printf("%s %zu:\n", label, len);
    
    for(size_t i=0; i < len; ++i)
        printf("%02x", buff[i]);
   
    printf("\n\n");
}

int testall() {
    int ret = 0;
    EC_KEY *ecpriv = NULL;
    EC_KEY *ecpub = NULL;
    EVP_PKEY *priv = NULL;
    EVP_PKEY *pub = NULL;
    EVP_PKEY *pub2 = NULL;
    const EC_POINT *ecpoint = NULL;
    EVP_PKEY_CTX *sctx = NULL;
    EVP_PKEY_CTX *vctx = NULL;
    EVP_PKEY_CTX *vctx2 = NULL;

    // hash msg
    char * msg = "hello";
    unsigned char data[EVP_MAX_MD_SIZE] = {0};
    unsigned int length;
    unsigned char * hash = NULL;
    unsigned char * sig = NULL;
    char * pk = NULL;
    char * privByte = NULL;
    BIO *out2  =NULL;
    BIO *out = NULL;
    BIO *in = NULL;

    // generate sm2 priv Key
    ecpriv = EC_KEY_new_by_curve_name(NID_sm2);
    if (ecpriv == NULL) {
        goto err;
    }
    ret = EC_KEY_generate_key(ecpriv);
    if (ret != 1) {
        goto err;
    }

    // extrac sm2 public key
    ecpoint = EC_KEY_get0_public_key(ecpriv);
    if (ecpoint == NULL) {
        goto err;
    }
    ecpub = EC_KEY_new_by_curve_name(NID_sm2);
    if (ecpub == NULL) {
        goto err;
    }
    ret = EC_KEY_set_public_key(ecpub, ecpoint);
    if (ret != 1 || ecpub == NULL) {
        goto err;
    }

    // EVP_PKEY sm2 priv, pub
    if((priv = EVP_PKEY_new()) == NULL) {
        goto err;
    }
    if ((ret = EVP_PKEY_set1_EC_KEY(priv, ecpriv)) != 1) {
        goto err;
    }
    if((pub = EVP_PKEY_new()) == NULL) {
        goto err;
    }
    if ((ret = EVP_PKEY_set1_EC_KEY(pub, ecpub)) != 1) {
        goto err;
    }

    size_t msg_length = strlen(msg);
    if ((ret = EVP_Digest(msg, msg_length, data, &length, EVP_sm3(), NULL)) != 1) {
        goto err;
    }
    hash = OPENSSL_malloc(length);
    memcpy(hash, data, length);

    // sign 
    if ((ret = EVP_PKEY_set_alias_type(priv, EVP_PKEY_SM2)) != 1) {
        goto err;
    }
    if ((sctx = EVP_PKEY_CTX_new(priv, NULL)) == NULL) {
        goto err;
    }
    if ((ret = EVP_PKEY_sign_init(sctx)) != 1) {
        goto err;
    }
    if ((ret = EVP_PKEY_CTX_set_signature_md(sctx, EVP_sm3())) != 1) {
        goto err;
    }
    size_t siglen = 0;
    if (EVP_PKEY_sign(sctx, NULL, &siglen, hash, length) <= 0) goto err;
    if (!(sig = (unsigned char *)OPENSSL_malloc(sizeof(unsigned char) * siglen))) goto err;
    if (EVP_PKEY_sign(sctx, sig, &siglen, hash, length) <= 0) goto err;

    // verify
    if ((ret = EVP_PKEY_set_alias_type(pub, EVP_PKEY_SM2)) != 1) {
        goto err;
    }
    if ((vctx = EVP_PKEY_CTX_new(pub, NULL)) == NULL) {
        printf("Pub failed\n");
        goto err;
    }
     if ((ret = EVP_PKEY_verify_init(vctx)) != 1) {
         printf("vtx\n");
        goto err;
    }
    if ((ret = EVP_PKEY_CTX_set_signature_md(vctx, EVP_sm3())) != 1) {
        printf("1\n");
        goto err;
    }
    if (EVP_PKEY_verify(vctx, sig, siglen, hash, length) != 1) {
        printf("Verify failed\n");
        goto err;
    }

    if ((out = BIO_new(BIO_s_mem())) == NULL) {
        goto err;
    }
    if ((ret = i2d_PUBKEY_bio(out, pub)) != 1) {
        printf("i2d failed");
        goto err;
    }
    int len = 0;
    if ((len = BIO_ctrl_pending(out)) <= 0) {
        ret = -1;
        printf("No Data\n");
        goto done;
    }
    pk = (char *) OPENSSL_malloc(len);
    if ((len = BIO_read(out, pk, len)) <= 0) {
        goto err;
    }

    // parse pk
    if ((in = BIO_new_mem_buf(pk, len)) == NULL) {
        goto err;
    }
    if ((pub2 = d2i_PUBKEY_bio(in, NULL)) == NULL) {
        printf("d2i failed\n");
        goto err;
    }
    // verify2
    if ((ret = EVP_PKEY_set_alias_type(pub2, EVP_PKEY_SM2)) != 1) {
        goto err;
    }
    if ((vctx2 = EVP_PKEY_CTX_new(pub2, NULL)) == NULL) {
        printf("Pub failed\n");
        goto err;
    }
     if ((ret = EVP_PKEY_verify_init(vctx2)) != 1) {
         printf("vtx\n");
        goto err;
    }
    if ((ret = EVP_PKEY_CTX_set_signature_md(vctx2, EVP_sm3())) != 1) {
        printf("1\n");
        goto err;
    }
    if (EVP_PKEY_verify(vctx2, sig, siglen, hash, length) != 1) {
        printf("Verify failed\n");
        goto err;
    }

    if ((out2 = BIO_new(BIO_s_mem())) == NULL) {
        goto err;
    }
    if ((ret = i2d_PKCS8PrivateKey_bio(out2, priv, NULL, NULL, 0, NULL, NULL)) != 1) {
        printf("i2d priv failed");
        goto err;
    }
    if ((len = BIO_ctrl_pending(out2)) <= 0) {
        ret = -1;
        printf("No Data\n");
        goto done;
    }
    privByte = (char *) OPENSSL_malloc(len);
    if ((len = BIO_read(out2, privByte, len)) <= 0) {
        goto err;
    }

    goto done;

err:
    printerror();
    goto done;

done:
    EC_KEY_free(ecpub);
    EC_KEY_free(ecpriv);
    EVP_PKEY_free(pub2); 
    EVP_PKEY_free(pub); 
    EVP_PKEY_free(priv);
    EVP_PKEY_CTX_free(sctx);
    EVP_PKEY_CTX_free(vctx);
    EVP_PKEY_CTX_free(vctx2);

    OPENSSL_free(hash);
    OPENSSL_free(sig);
    OPENSSL_free(pk);
    OPENSSL_free(privByte);
    BIO_free(out2);

    BIO_free(out);
    BIO_free(in);
    return ret;
}

int ensm2() {
    int ret = 0;

    EVP_PKEY *priv = NULL;
   
    BIO *out = NULL;
    char *privByte = NULL;

    FILE *fp = NULL;
    if ((fp = fopen("priv.pem", "r")) == NULL) {
        goto err;
    }

    priv = PEM_read_PrivateKey(fp, NULL, NULL, NULL);
    if (priv == NULL) {
        goto err;
    }
    fclose(fp);


    int len = 0;

    if ((out = BIO_new(BIO_s_mem())) == NULL) {
        goto err;
    }
    if ((ret = i2d_PKCS8PrivateKey_bio(out, priv, NULL, NULL, 0, NULL, NULL)) != 1) {
        printf("i2d priv failed");
        goto err;
    }
    if ((len = BIO_ctrl_pending(out)) <= 0) {
        ret = -1;
        printf("No Data\n");
        goto done;
    }
    privByte = (char *) OPENSSL_malloc(len);
    if ((len = BIO_read(out, privByte, len)) <= 0) {
        goto err;
    }

    goto done;

err:
    printerror();
    goto done;

done:
    EVP_PKEY_free(priv);
    OPENSSL_free(privByte);
    BIO_free(out);
    return ret;
}

EVP_PKEY *pem_read_priv(char *filename) {
    FILE *fp = NULL;
    EVP_PKEY *priv= NULL;
    if ((fp = fopen("priv.pem", "r")) == NULL) {
        printf("open failed");
        goto err;
    }

    priv = PEM_read_PrivateKey(fp, NULL, NULL, NULL);
    fclose(fp);
    if (priv == NULL) {
        printf("Read failed");
        goto err;
    }
    goto done;
err:
    printerror();
    goto done;

done:
    return priv;
}

int i2d_Private(EVP_PKEY *priv, char **data, int *length) {
    BIO * out = NULL;
    int ret = 0;
    int len = 0;
    if ((out = BIO_new(BIO_s_mem())) == NULL) {
        goto err;
    }
    // if ((ret = i2d_PKCS8PrivateKey_bio(out, priv, NULL, NULL, 0, NULL, NULL)) != 1) {
    //     goto err;
    // }
    if ((ret = i2d_PrivateKey_bio(out, priv)) != 1) {
        goto err;
    }
    if ((len = BIO_ctrl_pending(out)) <= 0) {
        goto err;
    }

    if ((len = BIO_read(out, *data, len)) <= 0) {
        goto err;
    }
    *length = len;
    goto done;
err:
    ret = -1;
done: 
    BIO_free(out);
    return ret;
}
int i2d_PrivKey_bio(EVP_PKEY *priv, char *data, int *length) {
    return i2d_Private(priv, &data, length);
}

int testEnPriv() {
    int ret = 1;
    EVP_PKEY * priv = NULL;
    char *byte = (char *) OPENSSL_malloc(512);
    int len = 0;
    priv = pem_read_priv("priv.pem");
    if (priv == NULL) {
        ret = -1;
        printf("read failed\n");
        goto done;
    }
    ret = i2d_PrivKey_bio(priv, byte, &len);
    if (ret != 1 || len <= 0) {
        printf("en failed\n");
        ret = -1;
        goto done;
    }

done:
    EVP_PKEY_free(priv);
    OPENSSL_free(byte);
    return ret;
}

int main() {
    int ret = testall();
    if (ret != 1) {
        printf("Error\n");
    }
    ret = ensm2();
    if (ret <= 0) {
        printf("Error2\n");
    }
    ret = testEnPriv();
    if (ret <= 0) {
        printf("Error2\n");
    }
    return 0;
}