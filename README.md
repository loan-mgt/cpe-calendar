# cpe-calendar
Sync your CPE calendar with your personal one


```
openssl genrsa -out secret/private.pem 2048
openssl rsa -in secret/private.pem -pubout > static/public.pem
```