package main

import "cacheme/interface/xnet"

// import "cacheme/interface/xhttp"

func main() {
	// xhttp.BoostrapHTTPServer()
	xnet.BootstrapTCPServer()
}
