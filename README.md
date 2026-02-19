# Load Balancer in Go

A HTTP load balancer built in Go to understand real backend system tradeoffs such as reverse proxying, load-balancing strategies, health checks, per-IP rate limiting, graceful shutdown, and Prometheus-based observability.

---

## Architecture Overview

![Load Balancer Architecture](docs/LB.svg)

### Request Flow
1. Client sends an HTTP request to the load balancer
2. Per-IP rate limiter validates the request using a token bucket algorithm
3. If allowed, the load-balancing strategy selects a backend
4. The server pool routes the request to a healthy backend
5. Backend response is proxied back to the client

### Health Checks
Health checks run asynchronously at fixed intervals and continuously update
backend health (`IsAlive`).

The server pool consults this state while selecting backends, ensuring traffic is
not routed to unhealthy instances.

---

## Features

- Reverse proxy–based HTTP load balancing
- Pluggable load-balancing strategies
  - Round Robin
  - Least Connections
- Per-IP rate limiting using Token Bucket algorithm
- Periodic backend health checks
- Graceful shutdown using OS signals and context cancellation
- Prometheus metrics for observability
- Structured logging with Zap

---

## Metrics & Observability

The load balancer exposes Prometheus metrics at:


### Available Metrics
- `lb_requests_total{backend,method,status}` – total requests handled
- `lb_request_duration_seconds` – request latency histogram
- `lb_backend_up` – backend health status
- `lb_active_connections` – active connections per backend
- `lb_rate_limited_requests_total` – requests blocked by rate limiter

These metrics provide visibility into traffic patterns, backend health, and
overall system performance.

---

## Configuration (`config.yaml`)

```yaml
lb_port: ":3332"
strategy: "round-robin"
backends:
  - "http://localhost:9001"
  - "http://localhost:9002"
```

### Configuration Options
- `lb_port` – Listening port (configurable; can be overridden via environment in production)
- `strategy` – Load-balancing strategy (round-robin, least-connection)
- `backends` – List of backend service URLs


## How to Run
### Start backend servers (example)
```bash
python3 -m http.server 9001
python3 -m http.server 9002
```
### Run the load balancer
```bash
go run main.go
```
### Send requests
```bash
curl http://localhost:3332
```

### View Prometheus metrics
```bash
curl http://localhost:3332/metrics
```