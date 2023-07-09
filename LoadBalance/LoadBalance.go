package loadbalance

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type LBMode int

const (
	RANDOM LBMode = iota
	ROUNDROBIN
)

type LoadBalance struct {
	addrs []string
	lock  sync.RWMutex

	mode LBMode
	rg   *rand.Rand
	idx  int
}

func NewLoadBalance(addrs []string, mode LBMode) *LoadBalance {
	return &LoadBalance{
		addrs: addrs,
		lock:  sync.RWMutex{},
		mode:  mode,
		rg:    rand.New(rand.NewSource(time.Now().Unix())),
		idx:   0,
	}
}

func (lb *LoadBalance) Refresh() error {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	lb.addrs = nil
	return nil
}

func (lb *LoadBalance) UpdateAll(addrs []string, err error) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	lb.addrs = addrs
}

func (lb *LoadBalance) Get() (addr string, err error) {
	lb.lock.RLock()
	defer lb.lock.RUnlock()
	if lb.addrs == nil {
		err = fmt.Errorf("empty addrs")
		return
	}
	size := len(lb.addrs)
	switch lb.mode {
	case RANDOM:
		addr = lb.addrs[lb.rg.Intn(size)]
	case ROUNDROBIN:
		addr = lb.addrs[lb.idx]
		lb.idx = (lb.idx + 1) % size
	}
	return
}

func (lb *LoadBalance) GetAll() (addrs []string, err error) {
	lb.lock.RLock()
	defer lb.lock.RUnlock()

	if lb.addrs == nil {
		err = fmt.Errorf("empty addrs")
		return
	}
	return lb.addrs, nil
}

func (lb *LoadBalance) Add(addrs ...string) (err error) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	if lb.addrs == nil {
		return fmt.Errorf("empty addrs")
	}
	lb.addrs = append(lb.addrs, addrs...)
	return
}
