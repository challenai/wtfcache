package buffer

import (
	"bufio"
	"bytes"
	"io"
)

type Buffer struct {
	w *bufio.Writer
	r *bufio.Reader
}

const (
	CR byte = 0x0d
	LF byte = 0x0a
)

var DELIMITER = []byte{0x0d, 0x0a}

func NewBuffer(buf io.ReadWriter) *Buffer {
	return &Buffer{
		w: bufio.NewWriter(buf),
		r: bufio.NewReader(buf),
	}
}

func (c *Buffer) WriteBytes(b []byte) {
	c.w.Write(b)
}

func (c *Buffer) WriteString(s string) {
	c.w.Write([]byte(s))
}

func (c *Buffer) WriteStringLine(s string) {
	c.w.Write([]byte(s))
	c.w.Write(DELIMITER)
}

func (c *Buffer) WriteStringEnd(s string) {
	c.w.Write([]byte(s))
	c.w.Write(DELIMITER)
	c.w.Flush()
}

func (c *Buffer) WriteBytesLine(b []byte) {
	c.w.Write(b)
	c.w.Write(DELIMITER)
}

func (c *Buffer) WriteBytesEnd(b []byte) {
	c.w.Write(b)
	c.w.Write(DELIMITER)
	c.w.Flush()
}

func (c *Buffer) ReadStringLine() (string, error) {
	return "", nil
}

func (c *Buffer) ReadBytesLine() ([]byte, error) {
	line, err := c.r.ReadSlice(LF)
	if err != nil {
		return nil, err
	}
	return trimLineCRLF(bytes.TrimSpace(line)), nil
}

func trimLineCRLF(ln []byte) []byte {
	if len(ln) == 0 {
		return ln
	}
	for len(ln) > 0 && (ln[len(ln)-1] == LF || ln[len(ln)-1] == CR) {
		ln = ln[:len(ln)-1]
	}
	return ln
}
