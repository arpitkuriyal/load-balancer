package healthcheck

import (
	"load-balancer/internal/pool"
	"time"
)

func Start(pool *pool.ServerPool, interval time.Duration) {
	for {
		for _, b := range pool.Backends {
			go b.HealthCheck()
		}
		time.Sleep(interval)
	}
}
