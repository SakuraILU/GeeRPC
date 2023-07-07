package client

import (
	codec "grpc/Codec"
)

type Call struct {
	Head      *codec.Head
	Argv      interface{}
	Reply     interface{}
	Done_chan chan bool
}

func NewCall(head *codec.Head, body interface{}) *Call {
	return &Call{
		Head:      head,
		Argv:      body,
		Done_chan: make(chan bool),
	}
}

func (c *Call) Done() {
	// close write chan, read operation is still ok
	close(c.Done_chan)
}
