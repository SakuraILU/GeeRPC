package registry

import (
	"log"
	"net/http"
	"time"
)

func HeartBeat(addr, registry string, timeout time.Duration) (err error) {
	c := &http.Client{}
	req, err := http.NewRequest("POST", registry, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-GRPC-Server", addr)

	if _, err = c.Do(req); err != nil {
		log.Fatal(err)
	}
	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		for err == nil {
			<-ticker.C
			if _, err = c.Do(req); err != nil {
				log.Fatal(err)
			}
		}
	}

	return
}
