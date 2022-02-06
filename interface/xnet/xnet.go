package xnet

import (
	"cacheme/cache"
	"cacheme/store"
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
	var c cache.Cache
	c, err = store.NewStore("./badger")
	if err != nil {
		log.Println("can't open store directory")
		panic(err)
	}
	server := CacheServer{
		// c: mem.NewMemCache(),
		c: c,
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
		client := NewClient(conn, &server)
		server.clients = append(server.clients, client)
		go client.Start()
	}
	return nil
}
