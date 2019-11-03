# etunnel

This service listens on one port and forwards to another.
One site or the other must be using a valid TLS certificate.

## Running the Server

Here what all the options look like, and shows their default values.

    ./etunnel \
    	--server=1.2.3.4:6379 \
    	--connect=127.0.0.1:6379
    	--server-certificate=server.pem \

## Running a Client

    ./etunnel \
    	--server=127.0.0.1:6379 \
    	--connect=1.2.3.4:6379
        --client-certificate=client.pem


## Enable SSL/TLS

Both Client and Server must have a certificate chain signed by the same Certificate Authority.

Please do this but use proper SSL Certificates.
One can rename the files, just pass proper options

    openssl genrsa -out server.key 2048
    openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
    openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
