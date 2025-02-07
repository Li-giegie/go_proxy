/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_proxy",
	Short: "http/https正向代理",
	Long:  `http/https Forward Proxy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, _ := cmd.Flags().GetString("addr")
		level, _ := cmd.Flags().GetUint32("loglevel")
		mode, _ := cmd.Flags().GetString("logmode")
		maxConnNum, _ := cmd.Flags().GetInt("maxconnnum")
		cacheSize, _ := cmd.Flags().GetInt("logfilecachesize")
		err := InitLog(level, mode, cacheSize)
		if err != nil {
			return err
		}
		Serve(addr, maxConnNum)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("addr", "a", ":1080", "listen address")
	rootCmd.Flags().Uint32P("loglevel", "l", 4, "log level: [0~6]")
	rootCmd.Flags().StringP("logmode", "m", "null", "log out mode: [null|stdout|$filename]")
	rootCmd.Flags().IntP("maxconnnum", "n", 100, "max connection number")
	rootCmd.Flags().IntP("logfilecachesize", "s", 0, "If the mode is Output to File, write to the buffer first, and enable the feature if the buffer size is greater than 16")
}

var file *os.File
var writer *bufio.Writer
var connNum int32

func InitLog(level uint32, mode string, logFileCacheSize int) error {
	if level < 0 || level > 6 {
		return fmt.Errorf("log level out of range: 0~6")
	}
	switch mode {
	case "null":
		logrus.SetOutput(io.Discard)
	case "stdout":
		logrus.SetOutput(os.Stdout)
	default:
		var err error
		file, err = os.OpenFile(mode, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logrus.Errorf("open file %s error", mode)
			return fmt.Errorf("open file %s error", mode)
		}
		var w io.Writer = file
		if logFileCacheSize > 16 {
			writer = bufio.NewWriter(file)
			w = writer
		}
		logrus.SetOutput(w)
	}
	logrus.SetLevel(logrus.Level(level))
	return nil
}

func Serve(addr string, maxConnNum int) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Errorf("listen err %v", err)
		return
	}
	defer func() {
		_ = l.Close()
		if writer != nil {
			_ = writer.Flush()
		}
		if file != nil {
			_ = file.Sync()
			_ = file.Close()
		}
		logrus.Infoln("server exit")
	}()
	logrus.Infof("Listening on %s", addr)
	exit := false
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		s := <-sigchan
		logrus.Warning("receive signal:", s.String())
		l.Close()
		exit = true
	}()
	for {
		count := 1
		for !exit && maxConnNum > 0 && atomic.LoadInt32(&connNum) >= int32(maxConnNum) {
			logrus.Warnf("connection number exceeds max conn num(%d) sleep %ds", connNum, count*3)
			time.Sleep(time.Second * 3)
			count++
		}
		conn, err := l.Accept()
		if err != nil {
			if conn != nil {
				_ = conn.Close()
			}
			if exit {
				return
			}
			return
		}
		go handle(conn)
	}
}

func handle(client net.Conn) {
	logrus.Infof("New connection %s", client.RemoteAddr())
	atomic.AddInt32(&connNum, 1)
	defer func() {
		v := recover()
		if v != nil {
			logrus.Panicf("Recover panic:%v", v)
		}
		logrus.Infof("Close client %v", client.RemoteAddr())
		client.Close()
		atomic.AddInt32(&connNum, -1)
	}()
	// 用来存放客户端数据的缓冲区
	var b [1024]byte
	//从客户端获取数据
	n, err := client.Read(b[:])
	if err != nil {
		logrus.Errorf("Read err %v", err)
		return
	}
	header, err := parseHTTPRequest(b[:n])
	if err != nil {
		logrus.Errorf("Parse header err %v", err)
		return
	}
	logrus.Debugf("header %#v", header)
	remoteAddr := header.Host
	// 如果方法是 CONNECT，则为 https 协议
	if header.Method == http.MethodConnect && remoteAddr == "" {
		remoteAddr = header.URI
	}
	//获得了请求的 host 和 port，向服务端发起 tcp 连接
	server, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		logrus.Errorf("dial peer err %v", err)
		return
	}
	defer server.Close()
	//如果使用 https 协议，需先向客户端表示连接建立完毕
	if header.Method == http.MethodConnect {
		if _, err = client.Write([]byte(header.Version + " 200 Connection established\r\n\r\n")); err != nil {
			logrus.Errorf("https connect reply client err %v\n", err)
			return
		}
	} else {
		if _, err = server.Write(b[:n]); err != nil {
			logrus.Errorf("https write server err %v", err)
			return
		}
	}
	logrus.Debugf("dial peer success %s\n", remoteAddr)
	go func() {
		if _, sErr := io.Copy(server, client); sErr != nil {
			logrus.Errorf("copy server to client err %v", sErr)
		}
	}()
	if _, err = io.Copy(client, server); err != nil {
		logrus.Errorf("copy client to server err %v", err)
	}
}

type Header struct {
	Method  string
	URI     string
	Version string
	Host    string
}

func parseHTTPRequest(p []byte) (*Header, error) {
	reader := bytes.NewBuffer(p)
	// 解析请求行
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Fields(requestLine)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", requestLine)
	}
	req := &Header{
		Method:  parts[0],
		URI:     parts[1],
		Version: parts[2],
	}
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		headerLine = strings.TrimSpace(headerLine)
		if headerLine == "" {
			break
		}
		headerParts := strings.SplitN(headerLine, ": ", 2)
		if len(headerParts) != 2 {
			return nil, fmt.Errorf("invalid header line: %s", headerLine)
		}
		if headerParts[0] == "Host" {
			if bytes.LastIndexByte([]byte(headerParts[1]), ':') == -1 {
				if req.Method == http.MethodConnect {
					req.Host = headerParts[1] + ":443"
				} else {
					req.Host = headerParts[1] + ":80"
				}
			} else {
				req.Host = headerParts[1]
			}
			break
		}
	}
	if req.Method != http.MethodConnect && req.Host == "" {
		return nil, fmt.Errorf("invalid header not found Host: %s", requestLine)
	}
	return req, nil
}
