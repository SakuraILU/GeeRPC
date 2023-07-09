package client

import (
	"context"
	"fmt"
	codec "grpc/Codec"
	loadbalance "grpc/LoadBalance"
	"log"
	"reflect"
)

type XClient struct {
	clients map[string]*Client
	lb      *loadbalance.LoadBalance
}

func NewXClient(network string, addrs []string, opt *codec.Option) (xc *XClient, err error) {
	xc = &XClient{
		clients: make(map[string]*Client),
		lb:      loadbalance.NewLoadBalance(addrs, loadbalance.ROUNDROBIN),
	}

	for _, addr := range addrs {
		c, err := NewClient(network, addr, opt)
		if err != nil {
			log.Println(err)
		}
		xc.clients[addr] = c
	}
	return
}

func (xc *XClient) Start() {
	for _, c := range xc.clients {
		c.Start()
	}
}

func (xc *XClient) Stop() {
	for _, c := range xc.clients {
		c.Stop()
	}
}

func (xc *XClient) Call(ctx context.Context, service_method string, args interface{}, reply interface{}) (err error) {
	addr, err := xc.lb.Get()
	if err != nil {
		log.Println(err)
		return
	}

	c, ok := xc.clients[addr]
	if !ok {
		log.Println("no client")
		return
	}

	err = c.Call(ctx, service_method, args, reply)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (xc *XClient) BroadCastctx(ctx context.Context, service_method string, args interface{}, reply interface{}) (err error) {
	if reflect.TypeOf(reply).Kind() != reflect.Ptr {
		log.Fatal("reply must be a ptr")
	}

	done_chan := make(chan interface{})
	for _, c := range xc.clients {
		go func(c *Client) {
			// reader may close channel after reading a valid reply,
			// so we need to recover the panic triggered by writing to a closed channel
			defer func() {
				recover()
			}()

			clone_reply := reflect.New(reflect.TypeOf(reply).Elem()).Interface()
			err := c.Call(ctx, service_method, args, clone_reply)
			if err != nil {
				done_chan <- err
				return
			}
			done_chan <- clone_reply
		}(c)
	}

	for i := 0; i < len(xc.clients); i++ {
		ret := <-done_chan

		switch v := ret.(type) {
		case error:
			// no need to do anything
		default:
			if reflect.TypeOf(v).Kind() != reflect.Ptr {
				log.Fatal("reply must be a ptr")
			}
			// *reply = *v
			reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(v).Elem())
			// [INFO]: a unrecommend and unsafe way to close a channel.
			// usually writer take the responsibility to close channel,
			// but here reader just read one reply value,
			// and writer don't know whether reader has read the reply value or not.
			// So, reader take the responsibility to close channel.
			// Meanwhile, writer will trigger panic when write to a closed channel
			// thus, we need to recover their panic in their goroutine
			close(done_chan)
			return
		}
	}

	return fmt.Errorf("all servers failed")
}
