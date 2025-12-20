package pool

import (
	"load-balancer/internal/backend"
	"sync"
)

type ServerPool struct {
	Backends     []*backend.Backend
	CurrentIndex int
	mux          sync.Mutex
}

func (p *ServerPool) Lock()   { p.mux.Lock() }
func (p *ServerPool) Unlock() { p.mux.Unlock() }
