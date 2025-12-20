package strategy

import (
	"load-balancer/internal/backend"
	"load-balancer/internal/pool"
)

type LeastConnection struct {
	pool *pool.ServerPool
}

func NewLeastConnection(p *pool.ServerPool) *LeastConnection {
	return &LeastConnection{pool: p}
}

func (lc *LeastConnection) Next() *backend.Backend {
	lc.pool.Lock()
	defer lc.pool.Unlock()

	minConn := -1
	var selected *backend.Backend

	for _, b := range lc.pool.Backends {
		alive, conns := b.Snapshot()
		if !alive {
			continue
		}

		if minConn == -1 || conns < minConn {
			minConn = conns
			selected = b
		}
	}

	return selected
}
