package main

import (
	"encoding/json"
	codec "grpc/Codec"
	"log"
	"net"
	"reflect"
	"time"
)

func startClient() {
	// connect to server
	conn, err := net.Dial("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("connect to server success")
	// send option
	opt, err := codec.NewOption(codec.GOBTYPE)
	if err != nil {
		log.Println(err)
		return
	}
	buf, err := json.Marshal(opt)
	conn.Write(buf)
	log.Println("read request success")

	codecer := codec.NewGobCodec(conn)

	for {
		// send request
		req := &codec.Request{
			Head: codec.Head{
				Service_id:     1,
				Service_method: "Hello",
				Error:          "",
			},
			Argv: reflect.ValueOf("grpc: ping"),
		}
		err = codecer.Write(&req.Head, req.Argv.Interface())
		if err != nil {
			log.Println(err)
			return
		}
		// log.Println("send request success")
		// read response
		var head codec.Head
		err = codecer.ReadHead(&head)
		if err != nil {
			log.Println(err)
			return
		}
		// log.Println("read response head success")
		var body string
		err = codecer.ReadBody(&body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(body)

		time.Sleep(time.Second * 2)
	}
}
