package main

import (
	"log"
	"os"
)

func main() {
	// arg: --server/ --client
	// read arg
	// if server
	// startServer()
	// else
	// startClient()
	arg := os.Args[1]

	log.SetFlags(log.Lshortfile)
	if arg == "server" {
		startServer()
	} else {
		startClient()
	}
}
