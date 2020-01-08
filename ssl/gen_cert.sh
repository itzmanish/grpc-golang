#!/bin/bash
# Regenerate the self-signed certificate for local host. Recent versions of firefox and chrome(ium)
# require a certificate authority to be imported by the browser (localhostCA.pem) while
# the server uses a cert and key signed by that certificate authority.
# Based partly on https://stackoverflow.com/a/48791236
CA_PASSWORD=ThisShouldBeStrongButString

# Generate the root certificate authority key with the set password
openssl genrsa -des3 -passout pass:$CA_PASSWORD -out CA.key 4096

# Generate a root-certificate based on the root-key for importing to browsers.
openssl req -x509 -new -nodes -key CA.key -passin pass:$CA_PASSWORD -config localhostCA.conf -sha256 -days 1825 -out CA.crt

# Generate a new private key
openssl genrsa -out server.key 4096

# Generate a Certificate Signing Request (CSR) based on that private key (reusing the
# localhostCA.conf details)
openssl req -new -key server.key -out server.csr -config localhostCA.conf

# Create the certificate for the webserver to serve using the localhost.conf config.
openssl x509 -req -in server.csr -CA CA.crt -CAkey CA.key -CAcreateserial \
-out server.crt -days 1024 -sha256 -extfile localhost.conf -passin pass:$CA_PASSWORD

openssl pkcs8 -topk8 -nocrypt -passin pass:$CA_PASSWORD -in server.key -out server.pem