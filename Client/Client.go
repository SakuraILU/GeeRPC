package client

import (
	"context"
	"encoding/json"
	"fmt"
	codec "grpc/Codec"
	"log"
	"net"
	"sync"
)

type Client struct {
	conn     net.Conn
	sendlock sync.RWMutex

	seq_id uint
	idlock sync.RWMutex
	cm     *CallManager
}

func NewClient(network, addr string, opt *codec.Option) (c *Client, err error) {
	var conn net.Conn
	conn, err = net.DialTimeout(network, addr, opt.ConnTimeout)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = json.NewEncoder(conn).Encode(&opt); err != nil {
		log.Fatal(err)
	}
	c = &Client{
		conn:     conn,
		sendlock: sync.RWMutex{},
		seq_id:   0,
		idlock:   sync.RWMutex{},
		cm:       NewCallManager(conn, opt.Codec_type),
	}
	return
}

func (c *Client) Start() {
	c.cm.Start()
}

func (c *Client) Stop() {
	c.cm.Stop()
}

func (c *Client) Call(ctx context.Context, service_method string, args interface{}, reply interface{}) (err error) {
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

	withLock(&c.sendlock,
		func() {
			if err = c.cm.codecr.Write(call.Head, call.Argv); err != nil {
				if err = c.cm.RemoveCall(call.Head.Service_id); err != nil {
					return
				}
				log.Fatalf("%s", err)
			}
		})

	select {
	case <-ctx.Done():
		err = fmt.Errorf("call timeout")
		c.cm.RemoveCall(call.Head.Service_id)
	case <-call.Done_chan:
		if call.Head.Error != "" {
			err = fmt.Errorf(call.Head.Error)
		}
	}

	return
}
