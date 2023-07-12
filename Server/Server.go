package server

import (
	"fmt"
	registry "grpc/Registry"
	service "grpc/Service"
	"log"
	"net"
	"time"
)

type Server struct {
	svices map[string]*service.Service
}

func NewServer() *Server {
	return &Server{
		svices: make(map[string]*service.Service),
	}
}

func (s *Server) Serve(listen_conn net.Listener, registry_addr string, timeout time.Duration) (err error) {
	go registry.HeartBeat(listen_conn.Addr().String(), registry_addr, timeout)

	for {
		conn, err := listen_conn.Accept()
		if err != nil {
			log.Fatal(err)
			break
		}

		connection := NewConnection(conn, s.svices)
		go func() {
			connection.Start()
			defer connection.Stop()
		}()
	}

	return
}

func (s *Server) RegisterService(any interface{}) (err error) {
	svice, err := service.NewService(any)
	if err != nil {
		return
	}

	if _, ok := s.svices[svice.GetName()]; ok {
		return fmt.Errorf("service %s already registered", svice.GetName())
	}
	s.svices[svice.GetName()] = svice

	return
}
