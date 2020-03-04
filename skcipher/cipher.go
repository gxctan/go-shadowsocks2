package skcipher

import (
	"net"
)

// 自定义编解码器
type Cipher struct {
	// 编码用的密码
	encodePassword *password
	// 解码用的密码
	decodePassword *password
}

// 加密原数据
func (cipher *Cipher) Encode(bs []byte) []byte {
	for i, v := range bs {
		bs[i] = cipher.encodePassword[v]
	}
	return nil
}

// 解码加密后的数据到原数据
func (cipher *Cipher) Decode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.decodePassword[v]
	}
}

// 新建一个编码解码器
func NewCipher(pass string) *Cipher {
	encodePassword, _ := parsePassword(pass)
	decodePassword := &password{}
	for i, v := range encodePassword {
		encodePassword[i] = v
		decodePassword[v] = byte(i)
	}
	return &Cipher{
		encodePassword: encodePassword,
		decodePassword: decodePassword,
	}
}

// 包装原始连接
func (cipher *Cipher) StreamConn(c net.Conn) net.Conn {
	return NewConn(c, *cipher)
}

func (cipher *Cipher) PacketConn(c net.PacketConn) net.PacketConn {
	return NewPacketConn(c, *cipher)
}
