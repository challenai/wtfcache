package main

import (
	"cacheme/interface/xnet"
	"log"
)

// import "cacheme/interface/xhttp"

func main() {
	log.SetFlags(log.Lshortfile)
	// xhttp.BoostrapHTTPServer()
	xnet.BootstrapTCPServer()
}
