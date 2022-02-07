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

type Conf struct {
	Host string
	Port int
}

type CacheServer struct {
	conf    *Conf
	c       cache.Cache
	ln      net.Listener
	clients []*Client
}

func BootstrapTCPServer(conf *Conf) error {
	var err error
	var c cache.Cache
	c, err = store.NewStore("./badger")
	if err != nil {
		log.Println("can't open store directory")
		panic(err)
	}
	server := CacheServer{
		// c: mem.NewMemCache(),
		conf: conf,
		c:    c,
	}
	addr := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))
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
