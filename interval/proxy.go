package interval

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
)

type Dialer interface {
	Dial() (net.Conn, error)
}

func StartForward(addr string, dialer Dialer) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logrus.Info("start forward", l.Addr())
	defer func() {
		l.Close()
		logrus.Info("stop forward", l.Addr())
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		s := <-c
		logrus.Info("receive signal", s.String())
		l.Close()
	}()
	for {
		src, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			defer src.Close()
			dst, err := dialer.Dial()
			if err != nil {
				logrus.Errorf("failed to dial: %v", err)
				return
			}
			defer dst.Close()
			go io.Copy(dst, src)
			io.Copy(src, dst)
		}()
	}
}

func StartProxy(listener net.Listener) error {
	defer func() {
		listener.Close()
		logrus.Info("stop proxy", listener.Addr())
	}()
	logrus.Info("start proxy", listener.Addr())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		s := <-c
		logrus.Info("receive signal", s.String())
		listener.Close()
	}()
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
				logrus.Errorf("failed to read request: %v", err)
				return
			}
			var dst net.Conn
			if req.Method != "CONNECT" {
				dst, err = net.Dial("tcp", DefaultHost(req.Host, "443"))
			} else {
				dst, err = net.Dial("tcp", DefaultHost(req.Host, "80"))
			}
			if err != nil {
				logrus.Errorf("failed to dial: %v", err)
				return
			}
			if req.Method == "CONNECT" {
				fmt.Fprintf(conn, "HTTP/%d.%d 200 Connection established\r\n\r\n", req.ProtoMajor, req.ProtoMinor)
			}
			defer dst.Close()
			go io.Copy(dst, src)
			io.Copy(conn, dst)
		}()
	}
}
