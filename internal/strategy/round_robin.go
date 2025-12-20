package strategy

import (
	"load-balancer/internal/backend"
	"load-balancer/internal/pool"
)

type RoundRobin struct {
	pool *pool.ServerPool
}

func NewRoundRobin(p *pool.ServerPool) *RoundRobin {
	return &RoundRobin{pool: p}
}

func (rr *RoundRobin) Next() *backend.Backend {
	rr.pool.Lock()
	defer rr.pool.Unlock()

	n := len(rr.pool.Backends)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		rr.pool.CurrentIndex = (rr.pool.CurrentIndex + 1) % n
		b := rr.pool.Backends[rr.pool.CurrentIndex]

		if b.IsAlive {
			return b
		}
	}

	return nil
}
