package server

import (
	service "grpc/Service"
	"io"
	"log"
	"net"
)

type Connection struct {
	rh *RequestReadWriter
}

func NewConnection(conn net.Conn, svices map[string]*service.Service) *Connection {
	rh, err := NewRequestReadWriter(conn, svices)
	if err != nil {
		log.Fatal(err)
	}

	return &Connection{
		rh: rh,
	}
}

func (c *Connection) Start() {
	for {
		req, svice, method, err := c.rh.ReadAndParse()
		if err != nil {
			req.Head.Error = err.Error()
			if err = c.rh.Write(&req.Head, struct{}{}); err != nil {
				if err == io.EOF {
					// client may close connection, so we just log the error and return
					log.Println(err)
					return
				} else {
					// other error may be caused by server, so we log and fatal, defense coding
					log.Fatal(err)
				}
			}
			continue
		}

		// TODO: goroutine worker pool
		// async handle request
		go func() {
			if err := svice.Call(method, req.Argv, req.Replyv); err != nil {
				req.Head.Error = err.Error()
			}

			if err = c.rh.Write(&req.Head, req.Replyv.Interface()); err != nil {
				log.Fatal(err)
			}
		}()
	}
}

func (c *Connection) Stop() {
	// TODO
}
