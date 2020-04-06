#!/bin/bash
# Regenerate the self-signed certificate for local host. Recent versions of firefox and chrome(ium)
# require a certificate authority to be imported by the browser (localhostCA.pem) while
# the server uses a cert and key signed by that certificate authority.
# Based partly on https://stackoverflow.com/a/48791236
CA_PATH=$(cd `dirname $0`; pwd)/CA
CONF_PATH=$(cd `dirname $0`; pwd)/config
CA_PASSWORD=notsafe

if [ $1 = "ca" ]; then

    echo "Build CA certificate using config at $CONF_PATH, password $CA_PASSWORD"

    # Create a directory to store CA key file.
    if [ ! -d "$CA_PATH" ]; then
        mkdir $CA_PATH
    fi

    # Generate the root certificate authority key with the set password
    openssl genrsa -des3 -passout pass:$CA_PASSWORD -out $CA_PATH/localhostCA.key 2048

    # Generate a root-certificate based on the root-key for importing to browsers.
    openssl req -x509 -new -nodes -key $CA_PATH/localhostCA.key -passin pass:$CA_PASSWORD -config $CONF_PATH/localhostCA.conf -sha256 -days 500 -out $CA_PATH/localhostCA.pem

fi

rm -f "./*"
rm -f "./CA.pem"
rm -f "./ca.cer"
rm -f "./$NAME.key"
rm -f "./$NAME.csr"
rm -f "./$NAME.crt"
rm -f "./$NAME.pem"
ln -s $CA_PATH/localhostCA.pem CA.pem

# Generate a new private key
openssl genrsa -out $NAME.key 2048

# Generate a Certificate Signing Request (CSR) based on that private key (reusing the
# localhostCA.conf details)
openssl req -new -key $NAME.key -out $NAME.csr -config $CONF_PATH/localhostCA.conf

# Create the certificate for the webserver to serve using the $NAME.conf config.
openssl x509 -req -in $NAME.csr -CA $CA_PATH/localhostCA.pem -CAkey $CA_PATH/localhostCA.key -CAcreateserial \
-out $NAME.crt -days 400 -sha256 -extfile $CONF_PATH/localhost.conf -passin pass:$CA_PASSWORD
