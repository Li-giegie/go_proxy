# go_proxy
go_proxy是一个http\https简单的正向代理，最终以二进制可执行程序运行，跨平台

## 构建
1. 源码构建，进入项目目录
```go
go build .
```
2. 远程构建

go install github.com/Li-giegie/go_proxy@latest

## 使用
go_proxy -h
```go
http/https Forward Proxy

Usage:
go_proxy [flags]

Flags:
-a, --addr string            listen address (default ":1080")
-h, --help                   help for go_proxy
-s, --logfilecachesize int   If the mode is Output to File, write to the buffer first, and enable the feature if the buffer size is greater than 16
-l, --loglevel uint32        log level: [0~6] (default 4)
-m, --logmode string         log out mode: [null|stdout|$filename] (default "null")
-n, --maxconnnum int         max connection number (default 100)
```
### 示例
1) ./go_proxy：启动默认配置的正向代理服务不输出日志
2) ./go_proxy -m run.log -s 4096 -l 6：输出日志到run.log文件中，缓冲区大小为4096，日志等级为6 (trace)
3) ./go_proxy -n 10：服务器最大接受10个连接

### http/s 加密隧道代理
代理的地址是隧道客户端
```
请求 <--> proxyclient <--> proxyserver <--> 目的服务端
```
双向证书认证
#### 生成证书
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

#### 启动proxy服务端
```
go_proxy proxyserver --listen :1080 --pem ./pem/server.pem --key ./pem/server.key --clientpem ./pem/client.pem
```
#### 启动proxy客户端
```
go_proxy proxyclient -l :1080 --proxy :2080 --pem ./pem/client.pem --key ./pem/client.key
```