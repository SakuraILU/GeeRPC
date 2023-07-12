package main

import (
	"context"
	"fmt"
	client "grpc/Client"
	codec "grpc/Codec"
	"log"
	"sync"
	"time"
)

func startClient() {
	// in sever, we set a sleep random 0~6 in Sort.BubbleSort,
	// so timeout 4 may cause timeout error and kill the goroutine
	opt, err := codec.NewOption(codec.GOBTYPE, time.Second*4, time.Second*10)
	if err != nil {
		log.Fatal(err)
	}

	// the http path of registry is raddr + rpath
	client, err := client.NewXClientRegistry("http://"+raddr+rpath, opt)
	if err != nil {
		log.Fatal(err)
	}

	client.Start()
	defer client.Stop()

	wg := sync.WaitGroup{}
	wg.Add(50)
	for n := 0; n < 50; n++ {
		go func() {
			for i := 0; i < 15; i++ {
				var reply string
				ctx, _ := context.WithTimeout(context.Background(), opt.Conn_timeout)
				err = client.Call(ctx, "Str.Reverse", "hello", &reply)
				fmt.Println(reply)
				if err != nil {
					log.Fatal(err)
				}
				args := []int{1, 3, 2, 4, 5, 6, 7, 8, 9, 0}
				reply2 := make([]int, len(args))
				ctx, _ = context.WithTimeout(context.Background(), opt.Conn_timeout)
				err = client.BroadCastctx(ctx, "Sort.BubbleSort", args, &reply2)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(reply2)
				time.Sleep(time.Second * 2)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
