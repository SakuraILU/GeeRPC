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
	opt, err := codec.NewOption(codec.GOBTYPE, time.Second*3, time.Second*3)
	if err != nil {
		log.Fatal(err)
	}

	addrs := []string{"127.0.0.1:10000", "127.0.0.1:10001",
		"127.0.0.1:10002"}
	client, err := client.NewXClient("tcp", addrs, opt)
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
				ctx, _ := context.WithTimeout(context.Background(), opt.ConnTimeout)
				err = client.Call(ctx, "Str.Reverse", "hello", &reply)
				fmt.Println(reply)
				if err != nil {
					log.Fatal(err)
				}
				args := []int{1, 3, 2, 4, 5, 6, 7, 8, 9, 0}
				reply2 := make([]int, len(args))
				ctx, _ = context.WithTimeout(context.Background(), opt.ConnTimeout)
				err = client.Call(ctx, "Sort.BubbleSort", args, &reply2)
				fmt.Println(reply2)
				time.Sleep(time.Second * 2)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
