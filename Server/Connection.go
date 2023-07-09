package server

import (
	service "grpc/Service"
	"io"
	"log"
	"net"
	"time"
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
			if err == io.EOF {
				// client close connection and trigger EOF error,
				// so we just log the error and return
				log.Println(err)
				return
			}
			req.Head.Error = err.Error()
			if err = c.rh.Write(&req.Head, struct{}{}); err != nil {
				log.Fatal(err)
			}
			continue
		}

		// TODO: goroutine worker pool
		// async handle request
		done := make(chan bool)
		work := func() {
			if err := svice.Call(method, req.Argv, req.Replyv); err != nil {
				req.Head.Error = err.Error()
			}

			if err = c.rh.Write(&req.Head, req.Replyv.Interface()); err != nil {
				// client may close connection, so we just log the error and return
				log.Println(err)
				return
			}
			close(done)
		}

		go func() {
			go work()

			select {
			case <-done:
				return
			case <-time.After(c.rh.timeout):
				req.Head.Error = "handle timeout"
				if err = c.rh.Write(&req.Head, struct{}{}); err != nil {
					// client may close connection, so we just log the error and return
					log.Println(err)
					return
				}
			}
		}()
	}
}

func (c *Connection) Stop() {
	// TODO
}
