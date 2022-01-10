package xnet

import (
	"cacheme/cache"
	"cacheme/mem"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	HOST = "localhost"
	PORT = 1235
)

type CacheServer struct {
	c     cache.Cache
	ln    net.Listener
	conns []net.Conn
}

func (cs *CacheServer) handleConn(conn net.Conn) {
	io.Copy(os.Stdout, conn)
}

func BootstrapTCPServer() {
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
			log.Printf("server can't accept connection: %s\n", err.Error())
			continue
		}
		server.conns = append(server.conns, conn)
		go server.handleConn(conn)
	}
}
