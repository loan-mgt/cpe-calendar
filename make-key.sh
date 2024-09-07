#!/bin/sh
(ls secret/private.pem > /dev/null 2>&1  && echo 'key already exists') || ( echo 'creating key' && openssl genrsa -out secret/private.pem 2048)
openssl rsa -in secret/private.pem -pubout > static/public.pem