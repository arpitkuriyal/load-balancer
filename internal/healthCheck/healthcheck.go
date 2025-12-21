package healthcheck

import (
	"context"
	"load-balancer/internal/pool"
	"time"
)

func Start(ctx context.Context, pool *pool.ServerPool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, b := range pool.Backends {
				go b.HealthCheck(ctx)
			}
		}

	}
}
