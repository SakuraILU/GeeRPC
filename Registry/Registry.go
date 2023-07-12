package registry

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Server struct {
	addr          string
	register_time time.Time
}

type Registry struct {
	servers map[string]*Server
	lock    sync.RWMutex
	timeout time.Duration // 0 means no timeout, aways keep alive
}

func NewRegistry(timeout time.Duration) *Registry {
	return &Registry{
		servers: make(map[string]*Server),
		lock:    sync.RWMutex{},
		timeout: timeout,
	}
}

func (r *Registry) addServer(addr string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.servers[addr] = &Server{
		addr:          addr,
		register_time: time.Now(),
	}

	log.Printf("register server from %s\n", addr)
}

func (r *Registry) aliveServers() (addrs []string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for addr, s := range r.servers {
		if s.register_time.Add(r.timeout).After(time.Now()) || r.timeout == 0 {
			addrs = append(addrs, addr)
		} else {
			delete(r.servers, addr)
		}
	}
	return
}

func (r *Registry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// note: in http, if w is not set anythin, it will return 200 OK

	switch req.Method {
	case "POST":
		addr := req.Header.Get("X-GRPC-Server")
		if addr == "" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			r.addServer(addr)
		}
	case "GET":
		addrs := r.aliveServers()
		w.Header().Set("X-GRPC-Servers", strings.Join(addrs, ","))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
