package xnet

import (
	"cacheme/cache"
	"cacheme/mem"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"
)

const (
	HOST = "localhost"
	PORT = 1235
)

type CacheServer struct {
	c       cache.Cache
	ln      net.Listener
	clients []*Client
}

func BootstrapTCPServer() error {
	var err error
	server := CacheServer{
		c: mem.NewMemCache(),
	}
	addr := net.JoinHostPort(HOST, strconv.Itoa(PORT))
	server.ln, err = net.Listen("tcp", addr)
	defer server.ln.Close()
	if err != nil {
		log.Fatalf("server can't listen to: %s", addr)
	}
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Printf("temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				return fmt.Errorf("listener.Accept() error - %s", err)
			}
			break
		}
		client := NewClient(conn)
		server.clients = append(server.clients, client)
		go client.Start()
	}
	return nil
}
