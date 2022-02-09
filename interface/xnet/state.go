package xnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

const CR byte = 0x0d
const LF byte = 0x0a

var SEP = []byte{32}
var deli = []byte{0x0d, 0x0a}

type Client struct {
	c net.Conn
	r *bufio.Reader
	w *bufio.Writer
	s *CacheServer
}

func NewClient(c net.Conn, s *CacheServer) *Client {
	return &Client{
		c: c,
		r: bufio.NewReader(c),
		w: bufio.NewWriter(c),
		s: s,
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
	// println(string(ln), err)
	if err != nil {
		return err
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
		c.w.Write([]byte(fmt.Sprintf("-ERR unknown command `%s`, with args beginning with:", fields[0])))
		c.Write(deli)
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
		if strings.Contains(err.Error(), "Key not found") {
			c.w.Write([]byte("$-1"))
			c.Write(deli)
			return
		}
		c.w.Write([]byte("-ERR "))
		c.Write(deli)
		return
	}
	c.w.Write([]byte("$" + strconv.Itoa(len(v))))
	c.w.Write(deli)
	c.w.Write(v)
	c.Write(deli)
}

func (c *Client) OpSet(k string, v []byte) {
	err := c.s.c.Set(k, v)
	if err != nil {
		log.Println(err)
		c.w.Write([]byte("-ERR \n"))
		c.w.Write(deli)
		c.w.Flush()
		return
	}
	c.w.Write([]byte("+OK"))
	c.Write(deli)
}
