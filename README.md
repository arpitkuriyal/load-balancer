# Load Balancer in Go

A simple HTTP load balancer built in Go to understand how reverse proxies,
load-balancing strategies, health checks, and rate limiting work in real systems.

## Architecture Overview

![Load Balancer Architecture](docs/LB.svg)

### Request Flow
1. Client sends an HTTP request to the load balancer
2. Rate limiter validates the request using a token bucket
3. If allowed, the load-balancing strategy selects a backend
4. Server pool routes the request to a healthy backend
5. Backend response is returned to the client

### Health Checks
Health checks run asynchronously and continuously update backend health
(`IsAlive`). The server pool consults this state when selecting backends.

---

## Features
- Reverse proxy–based HTTP load balancing
- Multiple load-balancing strategies
  - Round Robin
  - Least Connections
- Per-IP rate limiting (Token Bucket algorithm)
- Backend health checks
- Graceful shutdown with OS signals
- Structured logging (Zap)
- Clean, modular Go project structure

## Configuration (config.yaml)
```yaml
port: ":3332"
strategy: "round-robin"
backends:
  - "http://localhost:9001"
  - "http://localhost:9002"
```