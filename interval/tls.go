package interval

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"os"
)

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
