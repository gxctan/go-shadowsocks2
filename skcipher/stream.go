package skcipher

import (
	"io"
	"net"
)

const BUFFER_SIZE = 1024

type streamConn struct {
	net.Conn
	Cipher
}

// WriteTo reads from the embedded io.Reader, decrypts and writes to w until
// there's no more data to write or when an error occurs. Return number of
// bytes written to w and any error encountered.
func (c *streamConn) WriteTo(w io.Writer) (n int64, err error) {
	for {
		buf := make([]byte, BUFFER_SIZE)
		n, err := c.Read(buf)
		if err != nil {
			if err != io.EOF {
				return 0, err
			} else {
				return 0, nil
			}
		}
		if n > 0 {
			c.Decode(buf[:n])
			nw, err := w.Write(buf[:n])
			if err != nil {
				return 0, err
			}
			if n != nw {
				return 0, io.ErrShortWrite
			}
		}
	}

	return n, err
}

// ReadFrom reads from the given io.Reader until EOF or error, encrypts and
// writes to the embedded io.Writer. Returns number of bytes read from r and
// any error encountered.
func (c *streamConn) ReadFrom(r io.Reader) (n int64, err error) {
	for {
		buf := make([]byte, BUFFER_SIZE)
		n, er := r.Read(buf)
		if n > 0 {
			c.Encode(buf[:n])
			_, ew := c.Write(buf[:n])
			if ew != nil {
				err = ew
				break
			}
		}

		if er != nil {
			if er != io.EOF { // ignore EOF as per io.ReaderFrom contract
				err = er
			}
			break
		}
	}

	return n, err
}

// NewConn wraps a stream-oriented net.Conn with Cipher.
func NewConn(c net.Conn, ciph Cipher) net.Conn {
	return &streamConn{Conn: c, Cipher: ciph}
}
