package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"load-balancer/internal/backend"
	"load-balancer/internal/healthcheck"
	"load-balancer/internal/pool"
	"load-balancer/internal/strategy"
)

func main() {
	sp := &pool.ServerPool{}

	u1, _ := url.Parse("http://localhost:9001")
	u2, _ := url.Parse("http://localhost:9002")

	b1 := backend.NewBackend(u1)
	b2 := backend.NewBackend(u2)

	sp.Backends = append(sp.Backends, b1, b2)

	lbStrategy := strategy.NewRoundRobin(sp)

	go healthcheck.Start(sp, 5*time.Second)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := lbStrategy.Next()
		if b == nil {
			http.Error(w, "no backend available", http.StatusServiceUnavailable)
			return
		}
		b.Serve(w, r)
	})

	log.Println("Load balancer listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
