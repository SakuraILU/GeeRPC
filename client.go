package main

import (
	"fmt"
	client "grpc/Client"
	codec "grpc/Codec"
	"log"
	"net"
	"sync"
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
	client := client.NewClient(conn, codec.GOBTYPE)
	client.Start()
	defer client.Stop()

	wg := sync.WaitGroup{}
	wg.Add(50)
	for n := 0; n < 50; n++ {
		go func() {
			for i := 0; i < 15; i++ {
				var reply string
				err = client.Call("T.ping", "hello", &reply)
				fmt.Println(reply)
				if err != nil {
					log.Fatal(err)
				}
				time.Sleep(time.Second * 2)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
