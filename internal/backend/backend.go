package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	Url               *url.URL
	IsAlive           bool
	mux               sync.RWMutex
	activeConnections int
	proxy             *httputil.ReverseProxy
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.IsAlive = alive
	b.mux.Unlock()
}

func (b *Backend) Serve(w http.ResponseWriter, r *http.Request) {
	b.mux.Lock()
	b.activeConnections++
	b.mux.Unlock()

	defer func() {
		b.mux.Lock()
		b.activeConnections--
		b.mux.Unlock()
	}()

	b.proxy.ServeHTTP(w, r)
}

func NewBackend(u *url.URL) *Backend {
	proxy := httputil.NewSingleHostReverseProxy(u)
	return &Backend{
		Url:     u,
		IsAlive: true,
		proxy:   proxy,
	}
}

func (b *Backend) HealthCheck() {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(b.Url.String())
	if err != nil {
		b.SetAlive(false)
		return
	}
	defer resp.Body.Close()
	b.SetAlive(true)
}

func (b *Backend) Snapshot() (bool, int) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.IsAlive, b.activeConnections
}
