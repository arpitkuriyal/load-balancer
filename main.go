package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"load-balancer/config"
	"load-balancer/internal/backend"
	"load-balancer/internal/healthcheck"
	utils "load-balancer/internal/logger"
	"load-balancer/internal/pool"
	"load-balancer/internal/strategy"

	"go.uber.org/zap"
)

func main() {
	utils.InitLogger(os.Getenv("LOG_ENV"))
	defer utils.Sync()

	utils.Log.Info("starting load balancer")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.GetLbConfig("config.yaml")
	if err != nil {
		utils.Log.Fatal("failed to load config", zap.Error(err))
	}

	utils.Log.Info(
		"config loaded",
		zap.Strings("backends", cfg.Backends),
		zap.String("strategy", cfg.Strategy),
		zap.String("port", cfg.Port),
	)

	sp := &pool.ServerPool{}
	for _, rawURL := range cfg.Backends {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			utils.Log.Fatal(
				"invalid backend Urls",
				zap.String("url", rawURL),
				zap.Error(err),
			)
		}

		sp.Backends = append(sp.Backends, backend.NewBackend(parsedURL))

		utils.Log.Info(
			"backend registered",
			zap.String("backend", parsedURL.String()),
		)
	}

	var lbStrategy strategy.Strategy
	switch cfg.Strategy {
	case "round-robin":
		lbStrategy = strategy.NewRoundRobin(sp)
	case "least-connection":
		lbStrategy = strategy.NewLeastConnection(sp)
	default:
		utils.Log.Fatal("unsupported strategy", zap.String("strategy", cfg.Strategy))

	}

	utils.Log.Info("load balancing strategy initialized", zap.String("strategy", cfg.Strategy))

	go healthcheck.Start(ctx, sp, 5*time.Second)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := lbStrategy.Next()
		if b == nil {
			http.Error(w, "no backend available", http.StatusServiceUnavailable)
			return
		}

		utils.Log.Debug(
			"request forwarded",
			zap.String("backend", b.Url.String()),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
		)

		b.Serve(w, r)
	})
	server := &http.Server{
		Addr: cfg.Port,
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-shutdownCh
		utils.Log.Info("shutdown signal received", zap.String("signal", sig.String()))

		cancel()

		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTimeout()

		if err := server.Shutdown(ctxTimeout); err != nil {
			utils.Log.Error("server shutdown failed", zap.Error(err))
		}

		utils.Log.Info("graceful shutdown completed")
	}()

	utils.Log.Info("load balancer listening", zap.String("port", cfg.Port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.Log.Fatal("http server error", zap.Error(err))
	}

}
