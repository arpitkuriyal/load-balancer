package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "lb",
			Name:      "requests_total",
			Help:      "total requests handled by load balancer",
		},
		[]string{"backend", "method", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "lb",
			Name:      "request_duration_seconds",
			Help:      "request latency",
			Buckets:   []float64{0.05, 0.1, 0.2, 0.5, 1, 2, 5},
		},
		[]string{"backend"},
	)

	BackendUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lb",
			Name:      "backend_up",
			Help:      "Backend health status",
		},
		[]string{"backend"},
	)

	ActiveConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lb",
			Name:      "active_connections",
			Help:      "Active connections per backend",
		},
		[]string{"backend"},
	)

	RateLimitedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "lb",
			Name:      "rate_limited_requests_total",
			Help:      "Requests blocked by rate limiter",
		},
	)
)
