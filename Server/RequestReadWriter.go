package server

import (
	"fmt"
	codec "grpc/Codec"
	service "grpc/Service"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)

type RequestReadWriter struct {
	codecer  codec.ICodec
	sendlock sync.RWMutex
	svices   map[string]*service.Service
}

func NewRequestReadWriter(conn net.Conn, svices map[string]*service.Service) (rh *RequestReadWriter, err error) {
	// read option to get codec type
	opt, err := codec.ParseOption(conn)
	if err != nil || !opt.IsValid() {
		log.Println("Init request handler failed")
		return nil, err
	}

	codec_newfun, ok := codec.CodecNewFuncs[opt.Codec_type]
	if !ok {
		log.Fatal("wrong codec type")
	}
	rh = &RequestReadWriter{
		codecer: codec_newfun(conn),
		svices:  svices,
	}

	return rh, nil
}

func (rh *RequestReadWriter) ReadAndParse() (req *codec.Request, svice *service.Service, method *service.Method, err error) {
	errHandle := func(msg string) error {
		// send error msg back to client
		req.Head.Error = msg
		if err = rh.Write(&req.Head, struct{}{}); err != nil {
			log.Fatal(err)
		}
		return fmt.Errorf(msg)
	}

	req = &codec.Request{}

	// read Head
	var head codec.Head
	if err = rh.codecer.ReadHead(&head); err != nil {
		return
	}
	req.Head = head

	// get service and method name from Head
	names := strings.Split(req.Head.Service_method, ".")
	if len(names) != 2 {
		err = errHandle("invalid service_method name")
		return
	}
	sname, mname := names[0], names[1]
	// find service and method
	svice, ok := rh.svices[sname]
	if !ok {
		err = errHandle("invalid service name")
		return
	}
	method, err = svice.GetMethod(mname)
	if err != nil {
		err = errHandle("invalid method name")
		return
	}
	// read Body(argv), note ReadBody need a pointer,
	// so we need to get the pointer of argv if it is not a pointer
	req.Argv = method.NewArgv()
	argvp := req.Argv
	if req.Argv.Kind() != reflect.Ptr {
		argvp = req.Argv.Addr()
	}
	if err = rh.codecer.ReadBody(argvp.Interface()); err != nil {
		err = errHandle("read argv failed")
		return
	}

	// init reply
	req.Replyv = method.NewReply()

	return
}

func (rh *RequestReadWriter) Write(head *codec.Head, body interface{}) (err error) {
	// lock to avoid concurrent write
	rh.sendlock.Lock()
	defer rh.sendlock.Unlock()

	err = rh.codecer.Write(head, body)
	return
}
