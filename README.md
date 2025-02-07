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


