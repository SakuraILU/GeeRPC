package loadbalance

import (
	"log"
	"net/http"
	"strings"
)

type RegistryLoadBalance struct {
	raddr string
	lb    *LoadBalance
}

func NewRegistryLoadBalance(raddr string, mode LBMode) *RegistryLoadBalance {
	saddrs := getAliveServers(raddr)
	return &RegistryLoadBalance{
		raddr: raddr,
		lb:    NewLoadBalance(saddrs, mode),
	}
}

func getAliveServers(raddr string) []string {
	c := &http.Client{}

	req, err := http.NewRequest("GET", raddr, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(err)
	}

	saddrs := resp.Header.Get("X-GRPC-Servers")
	return strings.Split(saddrs, ",")
}

func (rlb *RegistryLoadBalance) Refresh() error {
	return rlb.lb.Refresh()
}

func (rlb *RegistryLoadBalance) UpdateAll() (err error) {
	saddrs := getAliveServers(rlb.raddr)
	return rlb.lb.UpdateAll(saddrs)
}

func (rlb *RegistryLoadBalance) Get() (addr string, err error) {
	return rlb.lb.Get()
}

func (rlb *RegistryLoadBalance) GetAll() (addrs []string, err error) {
	return rlb.lb.GetAll()
}
