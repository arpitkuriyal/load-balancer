package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type backend struct {
	Url     *url.URL
	IsAlive bool
	mux     sync.RWMutex
	// runtime load (for strategies like least-connections)
	activeConnections int
	proxy             *httputil.ReverseProxy
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

func (b *backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.IsAlive = alive
	b.mux.Unlock()
}

func (b *backend) healthCheck() {
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
func startHealthCheck(pool *serverPool, interval time.Duration) {
	for {
		for _, b := range pool.backends {
			go b.healthCheck()
		}
		time.Sleep(interval)
	}
}

func addBackend(u *url.URL) *backend {
	proxy := httputil.NewSingleHostReverseProxy(u)
	return &backend{
		Url:     u,
		IsAlive: true,
		proxy:   proxy,
	}
}

func (b *backend) Server(w http.ResponseWriter, r *http.Request) {
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
func main() {
	pool := &serverPool{}
	u1, _ := url.Parse("http://localhost:9001")
	u2, _ := url.Parse("http://localhost:9002")

	pool.backends = append(pool.backends,
		addBackend(u1),
		addBackend(u2),
	)

	go startHealthCheck(pool, 5*time.Second)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		backend := pool.NextBackend()
		if backend == nil {
			http.Error(w, "no backend available", http.StatusServiceUnavailable)
			return
		}
		backend.Server(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
