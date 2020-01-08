#!/bin/bash
# Regenerate the self-signed certificate for local host. Recent versions of firefox and chrome(ium)
# require a certificate authority to be imported by the browser (localhostCA.pem) while
# the server uses a cert and key signed by that certificate authority.
# Based partly on https://stackoverflow.com/a/48791236
CA_PASSWORD=ThisShouldBeStrongButString
CN=localhost


# Generate the root certificate authority key with the set password
openssl genrsa -des3 -passout pass:$CA_PASSWORD -out CA.key 4096

# Generate a root-certificate based on the root-key for importing to browsers.
openssl req -x509 -new -key CA.key -passin pass:$CA_PASSWORD -days 365 -out CA.crt -subj "/CN=${CN}"

# Generate a new private key
openssl genrsa -out server.key 4096

# Generate a Certificate Signing Request (CSR) based on that private key (reusing the
# localhostCA.conf details)
openssl req -new -key server.key -out server.csr -subj "/CN=${CN}"

# Create the certificate for the webserver to serve using the localhost.conf config.
openssl x509 -req -in server.csr -CA CA.crt -CAkey CA.key -set_serial 01 \
-out server.crt -days 365 -passin pass:$CA_PASSWORD

#Convert key to pem for grpc to work.
openssl pkcs8 -topk8 -nocrypt -passin pass:$CA_PASSWORD -in server.key -out server.pem