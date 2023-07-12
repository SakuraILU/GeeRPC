package main

import (
	server "grpc/Server"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
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
	// randomly sleep 1~3 seconds
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
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
	wg := sync.WaitGroup{}
	wg.Add(len(saddrs))
	for _, addr := range saddrs {
		go func(addr string) {
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatal(err)
			}
			s := server.NewServer()
			s.RegisterService(&Str{})
			s.RegisterService(&Sort{})
			err = s.Serve(lis, "http://127.0.0.1:9999/grpc/registry", 5*time.Second)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(addr)
	}
	wg.Wait()
}
