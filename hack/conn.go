package hack

import (
	"bytes"
	"net"
	"time"
)

type Identity int

const (
	IdentityUnknow Identity = iota
	IdentityHttp
	IdentityHttps
)

type IdentityConn struct {
	conn       net.Conn
	tempBuffer *bytes.Buffer
}

func NewIdentityConn(conn net.Conn) *IdentityConn {
	return &IdentityConn{
		conn:       conn,
		tempBuffer: bytes.NewBuffer([]byte{}),
	}
}

func (c *IdentityConn) Identify() (Identity, error) {
	temp := make([]byte, 0x800)
	n, err := c.conn.Read(temp)
	if err != nil {
		return IdentityUnknow, err
	}

	c.tempBuffer.Write(temp[:n])

	/*
		about tls protocol
		references `https://tls.ulfheim.net/`
		it will be help
	*/

	if temp[0] == 0x16 {
		return IdentityHttps, nil
	}

	return IdentityHttp, nil
}

/*
	Hack conn interface
*/
func (c *IdentityConn) Read(b []byte) (n int, err error) {
	// read cache first
	if c.tempBuffer.Len() > 0 {
		return c.tempBuffer.Read(b)
	}

	return c.conn.Read(b)
}

func (c *IdentityConn) Write(b []byte) (n int, err error) {
	n, err = c.conn.Write(b)
	return
}

func (c *IdentityConn) Close() error {
	return c.conn.Close()
}

func (c *IdentityConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *IdentityConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *IdentityConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *IdentityConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *IdentityConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
