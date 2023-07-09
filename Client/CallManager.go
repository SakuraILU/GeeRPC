package client

import (
	codec "grpc/Codec"
	"log"
	"net"
	"sync"
)

type CallManager struct {
	calls   map[uint]*Call
	maplock sync.RWMutex

	codecr codec.ICodec

	exit_chan chan bool
}

func NewCallManager(conn net.Conn, tp codec.CodecType) *CallManager {
	codec_newfun, ok := codec.CodecNewFuncs[tp]
	if !ok {
		log.Fatal("wrong codec type")
	}
	return &CallManager{
		calls:   make(map[uint]*Call),
		maplock: sync.RWMutex{},

		codecr: codec_newfun(conn),

		exit_chan: make(chan bool),
	}
}

func (cm *CallManager) Start() {
	go cm.receiver()
}

func (cm *CallManager) Stop() {
	close(cm.exit_chan)

	withLock(&cm.maplock,
		func() {
			for _, call := range cm.calls {
				call.Done()
			}
		})
}

func (cm *CallManager) AddCall(call *Call) (err error) {
	withLock(&cm.maplock,
		func() {
			id := call.Head.Service_id
			if _, ok := cm.calls[id]; ok {
				log.Fatal("add: id already exsit, something wrong...")
			}
			cm.calls[id] = call
		})

	return
}

func (cm *CallManager) RemoveCall(id uint) (err error) {
	withLock(&cm.maplock,
		func() {
			if _, ok := cm.calls[id]; !ok {
				log.Fatal("remove: id not exist, something wrong")
			}
			delete(cm.calls, id)
		})
	return
}

func (cm *CallManager) GetCall(id uint) (call *Call, err error) {
	withRLock(&cm.maplock,
		func() {
			_call, ok := cm.calls[id]
			if !ok {
				log.Fatal("get: id not exist, something wrong")
			}
			call = _call
		})

	return
}

func (cm *CallManager) receiver() {
	var err error
	for {
		select {
		case <-cm.exit_chan:
			return
		default:
			var head codec.Head
			if err = cm.codecr.ReadHead(&head); err != nil {
				log.Fatal("receiver read head")
			}

			withLock(&cm.maplock,
				func() {
					if head.Error != "" {
						log.Printf("id %d error: %s\n", head.Service_id, head.Error)
						cm.codecr.ReadBody(nil)
						return
					}
					call, ok := cm.calls[head.Service_id]
					if !ok {
						log.Printf("id %d is not exist, may timeout\n", head.Service_id)
						cm.codecr.ReadBody(nil)
						return
					}

					if err = cm.codecr.ReadBody(call.Reply); err != nil {
						log.Fatalf("%s", err)
					}
					delete(cm.calls, head.Service_id)
					call.Done()
				})
		}
	}
}

func withLock(lock sync.Locker, f func()) {
	lock.Lock()
	defer lock.Unlock()
	f()
}

func withRLock(lock *sync.RWMutex, f func()) {
	lock.RLock()
	defer lock.RUnlock()
	f()
}
