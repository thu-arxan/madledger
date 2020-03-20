src="uid_sm2.c"
if [ "$1" != "" ]; then
    src="$1"
fi
# gcc -Wall  -o main -L/usr/local/openssl/lib -lcrypto -lssl -I/usr/local/openssl/include  util.c $src
gcc -g -o main  $src -Wall -L/usr/local/openssl/lib -lcrypto -lssl -I/usr/local/openssl/include