package server

import (
	"net"
)

type Server struct {
	conn net.Conn
}

func NewServer(conn net.Conn) *Server {
	return &Server{
		conn: conn,
	}
}

func (s *Server) Serve() (err error) {
	rh, err := NewRequestHandler(s.conn)
	if err != nil {
		return err
	}
	for {
		if err = rh.Handle(); err != nil {
			return
		}
	}
}
