package main

import (
	registry "grpc/Registry"
	"net/http"
	"time"
)

func startRegistry() {
	reg := registry.NewRegistry(10 * time.Second)

	// http.Handle(rpath, reg)
	// http.ListenAndServe(raddr, nil)
	// [INFO]: For a single http server, Equivalent to above
	server_mux := http.NewServeMux()
	server_mux.Handle(rpath, reg)
	server := &http.Server{
		Addr:    raddr,
		Handler: server_mux,
	}
	server.ListenAndServe()
}
