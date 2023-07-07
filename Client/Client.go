package client

import (
	"encoding/json"
	"fmt"
	codec "grpc/Codec"
	"log"
	"net"
	"sync"
)

type Client struct {
	conn     net.Conn
	connlock sync.RWMutex

	seq_id uint
	idlock sync.RWMutex
	cm     *CallManager
}

func NewClient(conn net.Conn, tp codec.CodecType) *Client {
	opt, err := codec.NewOption(tp)
	if err != nil {
		log.Fatal(err)
	}
	if err = json.NewEncoder(conn).Encode(&opt); err != nil {
		log.Fatal(err)
	}
	return &Client{
		conn:     conn,
		connlock: sync.RWMutex{},
		seq_id:   0,
		idlock:   sync.RWMutex{},
		cm:       NewCallManager(conn, tp),
	}
}

func (c *Client) Start() {
	c.cm.Start()
}

func (c *Client) Stop() {
	c.cm.Stop()
}

func (c *Client) Call(service_method string, args interface{}, reply interface{}) (err error) {
	var call *Call
	withLock(&c.idlock,
		func() {
			call = &Call{
				Head: &codec.Head{
					Service_id:     c.seq_id,
					Service_method: service_method,
					Error:          "",
				},
				Argv:      args,
				Reply:     reply,
				Done_chan: make(chan bool),
			}
			c.seq_id++
		})

	if err = c.cm.AddCall(call); err != nil {
		return
	}

	withLock(&c.connlock,
		func() {
			if err = c.cm.codecr.Write(call.Head, call.Argv); err != nil {
				if err = c.cm.RemoveCall(call.Head.Service_id); err != nil {
					return
				}
				log.Fatalf("%s", err)
			}
		})

	<-call.Done_chan
	if call.Head.Error != "" {
		err = fmt.Errorf("ret error")
	}

	return
}
