package interval

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type TLSForward struct {
	ProxyAddr string
	Pem       string
	Key       string
	c         *tls.Config
}

func (x *TLSForward) Dial() (net.Conn, error) {
	if x.c == nil {
		conf, err := NewTLSConfig(x.Pem, x.Key)
		if err != nil {
			return nil, err
		}
		x.c = conf
	}
	return DialTLSClient(x.ProxyAddr, x.c.Clone())
}

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
				dst, err = net.Dial("tcp", DefaultHost(req.Host, "443"))
				if err != nil {
					fmt.Println("Error connecting to server:", err)
					return
				}
				defer dst.Close()
				fmt.Fprintf(conn, "HTTP/%d.%d 200 Connection established\r\n\r\n", req.ProtoMajor, req.ProtoMinor)
				go io.Copy(dst, src)
				io.Copy(conn, dst)
			} else {
				dst, err = net.Dial("tcp", DefaultHost(req.Host, "80"))
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

func DefaultHost(h, port string) string {
	if strings.IndexByte(h, ':') == -1 {
		return h + ":" + port
	}
	return h
}

func NewTLSConfig(pem, key string, poolPem ...string) (*tls.Config, error) {
	var config = new(tls.Config)
	cert, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		return nil, err
	}
	config.Certificates = []tls.Certificate{cert}
	config.InsecureSkipVerify = true
	if len(poolPem) > 0 {
		clientCertPool := x509.NewCertPool()
		for _, s := range poolPem {
			data, err := os.ReadFile(s)
			if err != nil {
				return nil, err
			}
			if ok := clientCertPool.AppendCertsFromPEM(data); !ok {
				return nil, errors.New("failed to append client cert")
			}
		}
		config.ClientCAs = clientCertPool
	}
	return config, nil
}

func DialTLSClient(addr string, conf *tls.Config) (net.Conn, error) {
	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewTLSListen(addr, pem, key, clientCert string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		return nil, err
	}
	cerdata, err := os.ReadFile(clientCert)
	if err != nil {
		return nil, err
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(cerdata)
	if !ok {
		return nil, errors.New("failed to parse root certificate")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	return tls.Listen("tcp", addr, config)
}
