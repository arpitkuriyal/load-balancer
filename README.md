# Load Balancer in Go

A simple HTTP load balancer built in Go to understand how reverse proxies, load-balancing strategies, health checks, and rate limiting work in real systems.


## Features

- Reverse proxy–based HTTP load balancing
- Multiple load‑balancing strategies
  - Round Robin
  - Least Connections
- Per‑IP rate limiting (Token Bucket algorithm)
- Backend health checks
- Graceful shutdown with OS signals
- Structured logging (Zap)
- Clean, modular Go project structure


---

## Configuration (`config.yaml`)

```yaml
port: ":3332"
strategy: "round-robin"
backends:
  - "http://localhost:9001"
  - "http://localhost:9002"
```



