# Load Balancer in Go


---

## ✨ Features

*  Reverse proxy–based HTTP load balancing
*  Multiple load‑balancing strategies
  * Round Robin
  * Least Connections
*  **Per‑IP rate limiting (Token Bucket algorithm)**
*  Backend health checks
*  Graceful shutdown with OS signals
*  Structured logging (Zap)
* Clean, modular Go project structure


---

## ⚙️ Configuration (`config.yaml`)

```yaml
port: ":3332"
strategy: "round-robin"
backends:
  - "http://localhost:9001"
  - "http://localhost:9002"
```



