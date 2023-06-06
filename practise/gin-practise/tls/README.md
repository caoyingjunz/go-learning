# How to Test

## Generate RSA private key and digital certificate

## Starting server with tls
```shell
go run main.go
```

## Curl with cert
```shell
curl --cacert ./cert/tls.crt https://127.0.0.1:443/welcome

# 正常回显
{"status":"success"}
```
