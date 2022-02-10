package xnet

import (
	"bytes"
	"cacheme/utils/buffer"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

var SEPERATOR = []byte{0x20}

type Client struct {
	s   *CacheServer
	buf *buffer.Buffer
}

func NewClient(c net.Conn, s *CacheServer) *Client {
	return &Client{
		buf: buffer.NewBuffer(c),
		s:   s,
	}
}

func (c *Client) Start() {
	var err error
	for {
		err = c.Read()
		if err != nil {
			break
		}
	}
}

func (c *Client) Read() error {
	ln, err := c.buf.ReadBytesLine()
	if err != nil {
		return err
	}

	fields := bytes.Split(ln, SEPERATOR)
	if len(fields) == 0 {
		// write to connection
		return errors.New("wrong parameters")
	}
	// we can only use bytes.Equal to compare
	switch {
	case bytes.Equal(fields[0], []byte("GET")):
		c.OpGet(fields)
	case bytes.Equal(fields[0], []byte("SET")):
		c.OpSet(fields)
	default:
		c.buf.WriteStringEnd(fmt.Sprintf("-ERR unknown command `%s`, with args beginning with:", fields[0]))
	}
	return nil
}

func (c *Client) OpGet(fields [][]byte) {
	if len(fields) < 2 {
		c.buf.WriteStringEnd("-ERR get args not enough")
		return
	}
	v, err := c.s.c.Get(string(fields[1]))
	if err != nil {
		if strings.Contains(err.Error(), "Key not found") {
			c.buf.WriteStringEnd("$-1")
			return
		}
		c.buf.WriteStringEnd("-ERR")
		return
	}
	c.buf.WriteStringLine(fmt.Sprintf("$%d", len(v)))
	c.buf.WriteBytesEnd(v)
}

func (c *Client) OpSet(fields [][]byte) {
	if len(fields) < 3 {
		c.buf.WriteStringEnd("-ERR set args not enough")
		return
	}
	err := c.s.c.Set(string(fields[1]), fields[2])
	if err != nil {
		log.Println(err)
		c.buf.WriteStringEnd("-ERR")
		return
	}
	c.buf.WriteStringEnd("+OK")
}
