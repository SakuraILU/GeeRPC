package server

import (
	"fmt"
	codec "grpc/Codec"
	"log"
	"net"
	"reflect"
)

type RequestHandler struct {
	conn    net.Conn
	codecer codec.ICodec
}

func NewRequestHandler(conn net.Conn) (rh *RequestHandler, err error) {
	opt, err := codec.ParseOption(conn)

	if err != nil {
		log.Println("Init request handler failed")
		return nil, err
	}
	var codecer codec.ICodec
	switch opt.Codec_type {
	case codec.GOBTYPE:
		codecer, err = codec.NewGobCodec(conn), nil
	case codec.JSONTYPE:
		codecer, err = codec.NewJsonCodec(conn), nil
	}
	if err != nil {
		log.Println("Init request handler failed")
		return nil, err
	}

	rh = &RequestHandler{
		conn:    conn,
		codecer: codecer,
	}

	return rh, nil
}

func (rh *RequestHandler) Read() (r *codec.Request, err error) {
	var head codec.Head
	if err := rh.codecer.ReadHead(&head); err != nil {
		return nil, err
	}

	argv := reflect.New(reflect.TypeOf(""))
	if err = rh.codecer.ReadBody(argv.Interface()); err != nil {
		return
	}
	log.Printf("argv: %v\n", argv.Elem().Interface())

	r = &codec.Request{
		Head: head,
		Argv: argv,
	}
	return
}

func (rh *RequestHandler) Write(head *codec.Head, body interface{}) (err error) {
	err = rh.codecer.Write(head, body)
	return
}

func (rh *RequestHandler) Handle() (err error) {
	req, err := rh.Read()
	if err != nil {
		return
	}

	// TODO: handle request
	// now just reply a message to client, assume the argv is string
	reply := fmt.Sprintf("grpc: pong %d", req.Head.Service_id)
	if err = rh.Write(&req.Head, reply); err != nil {
		return
	}

	return
}
