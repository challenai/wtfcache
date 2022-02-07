package xnet

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
)

const LF byte = 10
const CR byte = 13

var SEP = []byte{32}

type Client struct {
	c net.Conn
	r *bufio.Reader
	w *bufio.Writer
	// is a valid cache protocol
	valid bool
	s     *CacheServer
}

func NewClient(c net.Conn, s *CacheServer) *Client {
	return &Client{
		c:     c,
		r:     bufio.NewReader(c),
		w:     bufio.NewWriter(c),
		valid: false,
		s:     s,
	}
}

func (c *Client) Start() {
	// magic: `cch `
	magic := []byte("cch ")
	header := make([]byte, 4)
	n, err := io.ReadFull(c.c, header)
	if n != 4 || !bytes.Equal(magic, header) {
		log.Println("wrong client")
		c.Write([]byte("wrong protocol"))
		c.c.Close()
		return
	}
	if err != nil {
		log.Println(err)
		c.c.Close()
		return
	}
	c.Write([]byte("OK\n"))
	for {
		err = c.Read()
		if err != nil {
			break
		}
	}
}

func trimCRLF(ln []byte) []byte {
	if len(ln) == 0 {
		return ln
	}
	for len(ln) > 0 && (ln[len(ln)-1] == LF || ln[len(ln)-1] == CR) {
		ln = ln[:len(ln)-1]
	}
	return ln
}

func (c *Client) Read() error {
	ln, err := c.r.ReadSlice(LF)
	if err != nil {
		return err
	}
	if !c.valid {
		c.valid = true
		return nil
	}
	ln = bytes.Trim(ln, " ")
	ln = trimCRLF(ln)
	fields := bytes.Split(ln, SEP)
	if len(fields) == 0 {
		// write to connection
		return errors.New("wrong parameters")
	}
	// we can only use bytes.Equal to compare
	switch {
	case bytes.Equal(fields[0], []byte("GET")):
		c.OpGet(string(fields[1]))
	case bytes.Equal(fields[0], []byte("SET")):
		c.OpSet(string(fields[1]), fields[2])
	default:
		c.Write([]byte("- UNKNOWN COMMAND \n"))
	}
	return nil
}

func (c *Client) Write(b []byte) {
	c.w.Write(b)
	c.w.Flush()
}

func (c *Client) OpGet(k string) {
	v, err := c.s.c.Get(k)
	if err != nil {
		log.Println(err)
		c.Write([]byte("- ERROR \n"))
		return
	}
	c.w.Write(v)
	c.Write([]byte{LF})
}

func (c *Client) OpSet(k string, v []byte) {
	err := c.s.c.Set(k, v)
	if err != nil {
		log.Println(err)
		c.Write([]byte("- ERROR \n"))
		return
	}
	c.Write([]byte("+ OK \n"))
}
