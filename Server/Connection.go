package server

import (
	"fmt"
	"log"
	"net"
)

type Connection struct {
	rh *RequestHandler
}

func NewConnection(conn net.Conn) *Connection {
	rh, err := NewRequestHandler(conn)
	if err != nil {
		log.Fatal(err)
	}

	return &Connection{
		rh: rh,
	}
}

func (c *Connection) Start() {
	for {
		req, err := c.rh.Read()
		if err != nil && req != nil {
			req.Head.Error = err.Error()
			if err = c.rh.Write(&req.Head, struct{}{}); err != nil {
				log.Fatal(err)
			}
		}

		go func() {
			// TODO: handle request
			// now just reply a message to client, assume the argv is string
			reply := fmt.Sprintf("grpc: pong %d", req.Head.Service_id)
			if err = c.rh.Write(&req.Head, reply); err != nil {
				return
			}
		}()
	}
}

func (c *Connection) Stop() {
	// TODO
}
