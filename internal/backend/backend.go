package backend

import (
	"context"
	utils "load-balancer/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
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

func (b *Backend) HealthCheck(ctx context.Context) {
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		checkCtx,
		http.MethodGet,
		b.Url.String(),
		nil,
	)

	if err != nil {
		utils.Log.Debug(
			"health check request creation failed",
			zap.String("backend", b.Url.String()),
			zap.Error(err),
		)
		b.SetAlive(false)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Log.Warn(
			"health check failed",
			zap.String("backend", b.Url.String()),
			zap.Error(err),
		)
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
