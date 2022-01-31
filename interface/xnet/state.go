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

var SEP = []byte{32}

type Client struct {
	c net.Conn
	r *bufio.Reader
	w *bufio.Writer
	// is a valid cache protocol
	valid bool
}

func NewClient(c net.Conn) *Client {
	return &Client{
		c:     c,
		r:     bufio.NewReader(c),
		w:     bufio.NewWriter(c),
		valid: false,
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
	fields := bytes.Split(ln, SEP)
	if len(fields) == 0 {
		// write to connection
		return errors.New("wrong parameters")
	}
	// we can only use bytes.Equal to compare
	switch {
	case bytes.Equal(fields[0], []byte("GET")):
		c.OpGet()
	case bytes.Equal(fields[0], []byte("SET")):
		c.OpSet()
	default:
		c.Write([]byte("- UNKNOWN COMMAND"))
	}
	return nil
}

func (c *Client) Write(b []byte) {
	c.w.Write(b)
	c.w.Flush()
}

func (c *Client) OpGet() {
	c.Write([]byte("GET a key\n"))
}

func (c *Client) OpSet() {
	c.Write([]byte("SET a key\n"))
}
