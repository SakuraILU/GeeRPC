package server

import (
	"log"
	"net"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve(listen_conn net.Listener) (err error) {
	for {
		conn, err := listen_conn.Accept()
		if err != nil {
			log.Fatal(err)
			break
		}

		connection := NewConnection(conn)
		go func() {
			connection.Start()
			defer connection.Stop()
		}()
	}
	return
}
