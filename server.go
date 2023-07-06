package main

import (
	server "grpc/Server"
	"log"
	"net"
)

func startServer() {
	// short path

	listen_conn, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("service is listen on", listen_conn.Addr())
	svic := server.NewServer()

	conn, err := listen_conn.Accept()
	if err != nil {
		log.Println(err)
	}
	svic.Serve(conn)
}
