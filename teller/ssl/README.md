To generate certs, do:

`openssl ecparam -genkey -name secp384r1 -out server.key`
`openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`

And then

`curl https://localhost:443/hello -k`
