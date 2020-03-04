package skcipher

import (
	"net"
)

type packetConn struct {
	net.PacketConn
	Cipher
}

// WriteTo encrypts b and write to addr using the embedded PacketConn.
func (c *packetConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	buf := c.Encode(b)
	_, err := c.PacketConn.WriteTo(buf, addr)
	return len(buf), err
}

// ReadFrom reads from the embedded PacketConn and decrypts into b.
func (c *packetConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, addr, err := c.PacketConn.ReadFrom(b)
	if err != nil {
		return n, addr, err
	}
	c.Decode(b[:n])
	return n, addr, err
}

// NewPacketConn wraps a net.PacketConn with cipher
func NewPacketConn(c net.PacketConn, ciph Cipher) net.PacketConn {
	return &packetConn{PacketConn: c, Cipher: ciph}
}
