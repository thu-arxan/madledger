#include <stdio.h>
#include <string.h>
#include "util.h"
#include <openssl/evp.h>

int sm3_hash(const char *msg, size_t size, unsigned char * hash, unsigned int * hashlength) {
    return EVP_Digest(msg, size, hash, hashlength, EVP_sm3(), NULL);
}
int main(int argc, char const *argv[])
{
    /* code */
    char * msg = "hello";
    // if (argc > 1) {
    //     msg = argv[1];
    // }
    size_t size = strlen(msg);
    unsigned char hash[EVP_MAX_MD_SIZE] = {0};
    unsigned int length;
    sm3_hash(msg, size, hash, &length);

    print_it("hash", hash, length);

    return 0;
}
