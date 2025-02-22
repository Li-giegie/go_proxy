package interval

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func StartProxyServer(addr, pem, key, clientPem string) error {
	listener, err := NewTLSListen(addr, pem, key, clientPem)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("Proxy server listening on ", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			defer conn.Close()
			src := bufio.NewReader(conn)
			req, err := http.ReadRequest(src)
			if err != nil {
				fmt.Println("Error reading request:", err)
				return
			}
			var dst net.Conn
			if req.Method == "CONNECT" {
				dst, err = net.Dial("tcp", getHost(req.Host, "443"))
				if err != nil {
					fmt.Println("Error connecting to server:", err)
					return
				}
				defer dst.Close()
				fmt.Fprintf(conn, "HTTP/%d.%d 200 Connection established\r\n\r\n", req.ProtoMajor, req.ProtoMinor)
				go io.Copy(dst, src)
				io.Copy(conn, dst)
			} else {
				dst, err = net.Dial("tcp", getHost(req.Host, "80"))
				if err != nil {
					return
				}
				defer dst.Close()
				go req.WriteProxy(dst)
				io.Copy(conn, dst)
			}
		}()
	}
}
