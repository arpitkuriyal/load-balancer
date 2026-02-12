package healthCheck

import (
	"context"
	utils "load-balancer/internal/logger"
	"load-balancer/internal/pool"
	"time"

	"go.uber.org/zap"
)

func Start(ctx context.Context, pool *pool.ServerPool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	utils.Log.Info(
		"health check scheduler started",
		zap.Duration("interval", interval),
	)

	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			utils.Log.Info("health check scheduler stopped")
			return
		case <-ticker.C:
			for _, b := range pool.Backends {
				go b.HealthCheck(ctx)
			}
		}

	}
}
