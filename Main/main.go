package main

import (
	"log"
	"os"
)

var raddr string = "127.0.0.1:9999"
var rpath string = "/grpc/registry"
var saddrs []string = []string{"127.0.0.1:10000", "127.0.0.1:10001",
	"127.0.0.1:10002", "127.0.0.1:10005", "127.0.0.1:11000"}

func main() {
	// arg: --server/ --client
	// read arg
	// if server
	// startServer()
	// else if client
	// startClient()
	// else if registry
	// startRegistry()
	arg := os.Args[1]

	log.SetFlags(log.Lshortfile)
	if arg == "server" {
		startServer()
	} else if arg == "client" {
		startClient()
	} else if arg == "registry" {
		startRegistry()
	} else {
		log.Fatal("unsupported parameter")
	}
}
