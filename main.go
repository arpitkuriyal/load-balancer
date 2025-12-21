package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"load-balancer/config"
	"load-balancer/internal/backend"
	"load-balancer/internal/healthcheck"
	"load-balancer/internal/pool"
	"load-balancer/internal/strategy"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg, err := config.GetLbConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	sp := &pool.ServerPool{}
	for _, rawURL := range cfg.Backends {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			log.Fatalf("invalid backendUrl %s: %v", rawURL, err)
		}
		sp.Backends = append(sp.Backends, backend.NewBackend(parsedURL))
	}

	var lbStrategy strategy.Strategy
	switch cfg.Strategy {
	case "round-robin":
		lbStrategy = strategy.NewRoundRobin(sp)
	case "least-connection":
		lbStrategy = strategy.NewLeastConnection(sp)
	default:
		log.Fatalf("unsupported strategy: %s till now it only support 'round-robin' and 'least-coonection'", cfg.Strategy)
	}

	go healthcheck.Start(ctx, sp, 5*time.Second)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := lbStrategy.Next()
		if b == nil {
			http.Error(w, "no backend available", http.StatusServiceUnavailable)
			return
		}
		b.Serve(w, r)
	})

	fmt.Println("Load balancer listening on", cfg.Port)
	http.ListenAndServe(cfg.Port, nil)
}
