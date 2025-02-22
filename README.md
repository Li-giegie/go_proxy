# go_proxy
go_proxy是一个http\https代理工具

## 构建
1. 源码构建，进入项目目录
```go
go build .
```
2. 远程构建

go install github.com/Li-giegie/go_proxy@latest

## 使用
本代理分为两种模式： 转发模式、隧道模式，转发模式下从发送方到转发服务器为明文传输，如果访问外网可能被防火墙拦截，进而无法正常使用；
隧道模式发出的请求先进入隧道入口经过加密在传输到隧道服务端处理，这段数据是加密过的很大程度上防火墙不会拦截。

```go
go_proxy -h
http/s proxy server

Usage:
go_proxy [command]

Available Commands:
completion  Generate the autocompletion script for the specified shell
help        Help about any command
server      http/s proxy server (no tunnels)
tlsclient   http/s proxy tls tunnel client
tlsserver   http/s proxy tls tunnel server
version     version
xorclient   http/s proxy XOR encryption tunnel client
xorserver   http/s proxy XOR encryption tunnel server

Flags:
-h, --help   help for go_proxy
```
### 转发模式
```
请求 <---> proxyserver <---> 目的服务端
```
启动服务
```
go_proxy server -addr :1080
```
### 隧道模式
```
请求 <--> proxyclient <--> proxyserver <--> 目的服务端
```
代理入口为隧道入口 (隧道客户端)
#### 1.TLS
使用tls隧道双向认证
1. 启动proxy服务端
```
go_proxy tlsserver --addr :1080 --pem ./pem/server.pem --key ./pem/server.key --cpem ./pem/client.pem
```
2. 启动proxy客户端
```
go_proxy tlsclient --addr :1080 --proxy :2080 --pem ./pem/client.pem --key ./pem/client.key
```
#### 2.秘钥
使用秘钥加密隧道
1. 启动隧道服务端
```
go_proxy xorserver --addr :1080 --key a1!3&*dhiuSDHASd
```
2. 启动隧道客户端
```
go_proxy xorclient --addr :2080 --key a1!3&*dhiuSDHASd --proxy :1080
```

#### 生成证书
可使用Openssl工具生成证书
```
1、 生成服务器端的私钥
openssl genrsa -out server.key 2048
2、 生成服务器端证书
openssl req -new -x509 -key server.key -out server.pem -days 3650
3、 生成客户端的私钥
openssl genrsa -out client.key 2048
4、 生成客户端的证书
openssl req -new -x509 -key client.key -out client.pem -days 3650
```
