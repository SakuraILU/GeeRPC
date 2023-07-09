package main

import (
	server "grpc/Server"
	"log"
	"net"
	"strings"
	"time"
)

// define a simple add struct service

type Str struct {
}

func (s *Str) Lower2upper(argv string, reply *string) error {
	// convert lower case to upper case
	*reply = strings.ToUpper(argv)
	return nil
}

func (s *Str) Upper2lower(argv string, reply *string) error {
	// convert upper case to lower case
	*reply = strings.ToLower(argv)
	return nil
}

func (s *Str) Reverse(argv string, reply *string) error {
	// reverse string
	*reply = ""
	for i := len(argv) - 1; i >= 0; i-- {
		*reply += string(argv[i])
	}
	return nil
}

type Sort struct {
}

func (s *Sort) BubbleSort(argv []int, reply *[]int) error {
	// bubble sort
	time.Sleep(10 * time.Second)
	*reply = make([]int, len(argv))
	copy(*reply, argv)
	for i := 0; i < len(*reply); i++ {
		for j := i + 1; j < len(*reply); j++ {
			if (*reply)[i] > (*reply)[j] {
				(*reply)[i], (*reply)[j] = (*reply)[j], (*reply)[i]
			}
		}
	}
	return nil
}

func startServer() {
	listen_conn, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("service is listen on", listen_conn.Addr())

	if err != nil {
		log.Println(err)
	}

	sver := server.NewServer()
	sver.RegisterService(&Str{})
	sver.RegisterService(&Sort{})
	sver.Serve(listen_conn)
}
