#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <openssl/evp.h>
#include <openssl/ec.h>

int ECDSA_SIG_Encode(const ECDSA_SIG *sig, unsigned char **bytes) {
    return i2d_ECDSA_SIG(sig, bytes);
}

ECDSA_SIG *ECDSA_SIG_Decode(const unsigned char **pp, long len) {
    return d2i_ECDSA_SIG(NULL, pp, len);
}
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

int sm2(char *msg) {
    // sign msg
    unsigned char hash[EVP_MAX_MD_SIZE] = {0};
    unsigned int size;
    // 获取hash
    EVP_Digest(msg, strlen(msg), hash, &size, EVP_sm3(), NULL);

    print_it("hashed", hash, size);

    int     ret;
    ECDSA_SIG *sig;
    EC_KEY    *eckey;
    eckey = EC_KEY_new_by_curve_name(NID_sm2);
    EC_KEY_set_asn1_flag(eckey, OPENSSL_EC_NAMED_CURVE);
    

    if (eckey == NULL) {
        geterror();
        return 0;
    }
    // 生成公私钥对
    if (EC_KEY_generate_key(eckey) != 1) {
        geterror();
        return 0;
    }
    // 进行签名
    sig = ECDSA_do_sign(hash, size, eckey);
    if (sig == NULL)
    {
        geterror();
        return 0;
    }
    // 对签名做验证
    ret = ECDSA_do_verify(hash, size, sig, eckey);
    if (ret != 1) {
        printf("verify failed: %d\n", ret);
        geterror();
        return 0;
    }
    printf("verify passed\n");

    // encode to char
    // 将签名转换为char *
    unsigned char * sigdig = NULL;
    int sig_len = i2d_ECDSA_SIG(sig, &sigdig);;
    if (sig_len == 0) {
        geterror();
    }
    print_it("sig", sigdig, sig_len);


    BIGNUM *r=NULL;

    BIGNUM *s=NULL;
    ECDSA_SIG *newsig = ECDSA_SIG_new();
    if (newsig == NULL) {
        geterror();
        return ret;
    }
    char *sig_r = "2692C65D701711AFEE3536E7572C59D72F52AFE41F4F84C68673FA60A7E16E6D";
    char *sig_s = "86891D2794FCAB018053211C5B18FF88B676CCC69E1BF761F2EDB9F65D653AA7";
    ret = BN_hex2bn(&r, (char *)sig_r);
    if (ret == 0) {
        geterror();
        return ret;
    }
    ret = BN_hex2bn(&s, (char *)sig_s);
    if (ret == 0) {
        geterror();
        return ret;
    }
    // 304502202692c65d701711afee3536e7572c59d72f52afe41f4f84c68673fa60a7e16e6d02210086891d2794fcab018053211c5b18ff88b676ccc69e1bf761f2edb9f65d653aa7
    // 304502202692c65d701711afee3536e7572c59d72f52afe41f4f84c68673fa60a7e16e6d02210086891d2794fcab018053211c5b18ff88b676ccc69e1bf761f2edb9f65d653aa7
    ECDSA_SIG_set0(newsig,r,s);
    unsigned char * signDer = NULL;
    int signDerLen = 0;
    
    signDerLen = i2d_ECDSA_SIG(newsig, &signDer);
    print_it("NEWsig", signDer, signDerLen);

    // decode sig
    ECDSA_SIG *sig_new = NULL;
    // 解析签名
    // sig_new = d2i_ECDSA_SIG(NULL, &sigdig, sig_len);
    sig_new = ECDSA_SIG_Decode(&sigdig, sig_len);
    if (sig_new == NULL) {
        geterror();
        return 0;
    }

    geterror();

    // 再次编码
    sigdig = NULL;
    sig_len = ECDSA_SIG_Encode(sig_new, &sigdig);
    if (sig_len == 0) {
        geterror();
    }
    print_it("sig", sigdig, sig_len);

    // 验证签名
    ret = ECDSA_do_verify(hash, size, sig_new, eckey);
    if (ret != 1) {
        printf("verify failed: %d\n", ret);
        geterror();
        return 0;
    }
    printf("verify2 passed\n");

    const EC_POINT *ec_point = NULL; 
    ec_point = EC_KEY_get0_public_key(eckey);
    if (ec_point == NULL) {
        geterror();
        return 0;
    }
    EC_GROUP *ec_group = EC_KEY_get0_group(eckey);
    char * pubhex = EC_POINT_point2hex(ec_group, ec_point, POINT_CONVERSION_UNCOMPRESSED, NULL);
    printf("%ld PUB: %s\n", strlen(pubhex), pubhex);

    // decode
    EC_KEY *pubKey =  EC_KEY_new_by_curve_name(NID_sm2);
    EC_GROUP *group =  EC_KEY_get0_group(pubKey);
    EC_POINT *point = EC_POINT_hex2point(group, pubhex, NULL, NULL);
    if (point == NULL) {
        geterror();
        return 0;
    }
    ret = EC_KEY_set_public_key(pubKey, point);
    if (ret <= 0) {
        geterror();
        return 0;
    }

    ret = ECDSA_do_verify(hash, size, sig_new, pubKey);
    if (ret != 1) {
        printf("verify failed: %d\n", ret);
        geterror();
        return 0;
    }
    printf("verify3 passed\n");



    // encode private key
    // unsigned char * priv_str = NULL;
    // int priv_length = i2d_ECPrivateKey(eckey, NULL);
    // printf("length: %d\n", priv_length);
    // priv_str = OPENSSL_malloc(priv_length);
    // priv_length = i2d_ECPrivateKey(eckey, &priv_str);
    // printf("length: %d\n", priv_length);

    // printf("%s\n", priv_str);

    // print_it("priv", priv_str, priv_length);
    // // encode public key
    // unsigned char * pub_str = NULL, *p;
    // EC_KEY_set_asn1_flag(eckey, OPENSSL_EC_NAMED_CURVE);

    // int pub_length = i2o_ECPublicKey(eckey, NULL);
    // if (pub_length <= 0) {
    //     geterror();
    //     return;
    // }
    // printf("length: %d\n", pub_length);
    // pub_str = OPENSSL_malloc(pub_length);
    // p = pub_str;
    // pub_length = i2o_ECPublicKey(eckey, &p);
    // printf("length: %d\n", pub_length);

    // printf("%s\n", p);

    // print_it("priv", p, pub_length);
    // // decode private key

    // // decode public key
    // EC_KEY *pubKey =  EC_KEY_new_by_curve_name(NID_sm2);
    // EC_KEY_set_asn1_flag(pubKey, OPENSSL_EC_NAMED_CURVE);

    // pubKey = o2i_ECPublicKey(&pubKey, &p, pub_length);
    // if (pubKey == NULL) {
    //     geterror();
    //     return;
    // }

    // verify sig
    return 1;
}

int main() {
    printf("OpenSSL version: %s\n",OPENSSL_VERSION_TEXT);
    sm2("hello");
    return 0;
}