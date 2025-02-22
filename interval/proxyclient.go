package interval

import (
	"io"
	"log"
	"net"
	"strings"
)

func StartProxyClient(addr, proxyAddr, pem, key string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()
	log.Println("Proxy client listening on ", addr)
	config, err := NewTLSConfig(pem, key)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return err
		}
		go func() {
			defer conn.Close()
			dst, err := DialTLSClient(proxyAddr, config.Clone())
			if err != nil {
				log.Println(err)
				return
			}
			defer dst.Close()
			go io.Copy(dst, conn)
			io.Copy(conn, dst)
		}()
	}
}

func getHost(h, port string) string {
	if strings.IndexByte(h, ':') == -1 {
		return h + ":" + port
	}
	return h
}
