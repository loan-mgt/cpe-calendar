# cpe-calendar
Sync your CPE calendar with your personal one

# Setup

1. Frist step is to setup your own `.env` use as reference `exemple.env`
2.
```
openssl genrsa -out secret/private.pem 2048
openssl rsa -in secret/private.pem -pubout > static/public.pem
```