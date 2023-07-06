package server

import (
	"net"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve(conn net.Conn) (err error) {
	rh, err := NewRequestHandler(conn)
	if err != nil {
		return err
	}
	for {
		if err = rh.Handle(); err != nil {
			return
		}
	}
}
