package main

import (
	"cacheme/interface/xhttp"
	"cacheme/interface/xnet"
	"flag"
	"fmt"
	"log"
	"os"
)

const USAGE = `
Usage: wtfcache [options]

Common Options:
    -h, --help                       show help
    -v, --version                    print version
	--tcp                            start tcp server
	--http                           start http server
	--http_host                      http host
	--http_port                      http port
	--tcp_host                       tcp host
	--tcp_port                       tcp port
`

func main() {
	var (
		showHelp    bool
		showVersion bool
		useTCP      bool
		useHTTP     bool
	)
	log.SetFlags(log.Lshortfile)

	httpConf := &xhttp.Conf{}
	tcpConf := &xnet.Conf{}

	fs := flag.NewFlagSet("wtfcache", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Print(USAGE)
		os.Exit(0)
	}
	fs.BoolVar(&showHelp, "h", false, "show help")
	fs.BoolVar(&showVersion, "v", false, "print version")
	fs.BoolVar(&useHTTP, "http", false, "start http server")
	fs.BoolVar(&useTCP, "tcp", false, "start tcp server")
	fs.StringVar(&httpConf.Host, "http_host", "localhost", "http host")
	fs.IntVar(&httpConf.Port, "http_port", 1234, "tcp host")
	fs.StringVar(&tcpConf.Host, "tcp_host", "localhost", "tcp host")
	fs.IntVar(&tcpConf.Port, "tcp_port", 1235, "tcp port")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Println("config: can't parse command line flags: ", err)
		return
	}

	if useHTTP {
		xhttp.BoostrapHTTPServer(httpConf)
	}

	if useTCP || !useHTTP {
		xnet.BootstrapTCPServer(tcpConf)
	}
}
