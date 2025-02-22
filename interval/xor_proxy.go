package interval

import (
	"net"
)

type XORConn struct {
	net.Conn
	key         []byte
	readKeyPos  int
	writeKeyPos int
}

// 异或加密和解密函数，直接在原 data 上操作
func (x *XORConn) xorEncrypt(data []byte) {
	for i := range data {
		data[i] = data[i] ^ x.key[x.writeKeyPos%len(x.key)]
		x.writeKeyPos++
	}
}

func (x *XORConn) xorDecrypt(data []byte) {
	for i := range data {
		data[i] = data[i] ^ x.key[x.readKeyPos%len(x.key)]
		x.readKeyPos++
	}
}

// 实现 net.Conn 接口的 Read 方法
func (x *XORConn) Read(b []byte) (n int, err error) {
	n, err = x.Conn.Read(b)
	if err != nil {
		return
	}
	x.xorDecrypt(b[:n])
	return
}

// 实现 net.Conn 接口的 Write 方法
func (x *XORConn) Write(b []byte) (n int, err error) {
	x.xorEncrypt(b)
	return x.Conn.Write(b)
}

type XORForward struct {
	ProxyAddr string
	Key       []byte
}

func (f *XORForward) Dial() (net.Conn, error) {
	conn, err := net.Dial("tcp", f.ProxyAddr)
	return &XORConn{Conn: conn, key: f.Key}, err
}

type XORProxy struct {
	Key []byte
	net.Listener
}

func (p *XORProxy) Accept() (net.Conn, error) {
	conn, err := p.Listener.Accept()
	return &XORConn{Conn: conn, key: p.Key}, err
}
