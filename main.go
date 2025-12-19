package main

import (
	"fmt"
	"net/url"
	"sync"
)

type backend struct {
	Url     *url.URL
	IsAlive bool
	mux     sync.RWMutex
	// runtime load (for strategies like least-connections)
	activeConnections int
}

type serverPool struct {
	backends     []*backend
	currentIndex int
	strategy     string
	mux          sync.Mutex
}

func (sp *serverPool) NextBackend() *backend {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	n := len(sp.backends)
	if n == 0 {
		return nil
	}

	for range n {
		sp.currentIndex = (sp.currentIndex + 1) % n
		b := sp.backends[sp.currentIndex]

		b.mux.RLock()
		alive := b.IsAlive
		b.mux.RUnlock()

		if alive {
			return b
		}
	}

	return nil
}

func main() {
	fmt.Println("hello")
}
