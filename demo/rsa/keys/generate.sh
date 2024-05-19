# # 生成 CA 的 RSA 密钥
# openssl genrsa -out ca.key 2048
# # 生成 CA 证书
# openssl req -config conf/keys.cnf -new -x509 -days 3650 -key ca.key -out ca.crt

# # 为服务器生成 RSA 密钥
# openssl genrsa -out server.key 2048
# # 生成服务器的证书请求
# openssl req -config conf/keys.cnf -new -key server.key -out server_reqout.txt
# # 使用 CA 签名服务器的证书
# openssl x509 -req -in server_reqout.txt -days 3650 -sha256 -CAcreateserial -CA ca.crt -CAkey ca.key -out server.crt -extfile conf/extfile.cnf

# # 为客户端生成 RSA 密钥
# openssl genrsa -out client.key 2048
# # 生成客户端的证书请求
# openssl req -config conf/keys.cnf -new -key client.key -out client_reqout.txt
# # 使用 CA 签名客户端的证书
# openssl x509 -req -in client_reqout.txt -days 3650 -sha256 -CAcreateserial -CA ca.crt -CAkey ca.key -out client.crt

openssl genrsa -out server.key 2048

openssl req -new -key server.key -out server.csr -config conf/keys.cnf

openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
