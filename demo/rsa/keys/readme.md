[来源](https://idersan.github.io/2017/06/11/%E7%94%9F%E6%88%90ssl%E8%AF%81%E4%B9%A6/)

# 虚拟CA

```sh
# CA rsa 原始数据
openssl genrsa -out ca.key 2048
# 加密CA原始数据. -name指定某种加密曲线, 可以通过-list_curves显示
openssl ecparam -genkey -name secp384r1 -out ca.pem
# 生成CA证书 
openssl req -config conf/keys.cnf -newkey rsa:2048 -x509 -days 3650 -key ca.pem -out ca.crt 
```

# Server端证书

```sh
# 创建一个RSA密钥串
openssl genrsa -out server.key 2048
# 加密
openssl ecparam -genkey -name secp384r1 -out server.key

# 证书
openssl req -config conf/keys.cnf -new -key server.key -out server_reqout.txt 
# 签名. go自1.18开始, 不支持sha1签名的证书, 加密算法换成sha256
openssl x509 -req -in server_reqout.txt -days 3650 -sha256 -CAcreateserial -CA ca.crt -CAkey ca.pem -out server.crt -extfile conf/extfile.cnf

```

# Client端证书

```sh
# 创建一个RSA密钥串
openssl genrsa -out client.key 2048
# 加密
openssl ecparam -genkey -name secp384r1 -out client.key
# 证书
openssl req -config conf/keys.cnf -new -key client.key -out client_reqout.txt 
# 签名
openssl x509 -req -in client_reqout.txt -days 3650 -sha256 -CAcreateserial -CA ca.crt -CAkey ca.pem -out client.crt
```