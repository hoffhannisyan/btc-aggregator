# BTC Price Aggregation Service

Go microservice that polls BTC/USD prices from Coinbase, Kraken, and CoinDesk, aggregates them using median strategy, and exposes HTTP endpoints with Prometheus metrics.

## Package Structure

```
cmd/server/          — application entrypoint
internal/config/     — environment-based configuration
internal/fetcher/    — price source interface and implementations
internal/aggregator/ — median price calculation with stale fallback
internal/poller/     — background polling with retry and backoff
internal/api/        — HTTP handlers (/price, /health, /metrics)
internal/metrics/    — Prometheus metrics definitions
```

## How to Run

### Local

```bash
go run cmd/server/main.go
```

### Docker

```bash
docker compose up --build
```

## Configuration

All settings have sensible defaults and can be overridden via environment variables.

| Variable          | Default | Description         |
| ----------------- | ------- | ------------------- |
| `POLL_INTERVAL`   | `10s`   | Polling interval    |
| `REQUEST_TIMEOUT` | `5s`    | Per-request timeout |
| `MAX_RETRIES`     | `3`     | Max retry attempts  |
| `HTTP_PORT`       | `8080`  | Server port         |

### Override locally

```bash
POLL_INTERVAL=5s REQUEST_TIMEOUT=3s MAX_RETRIES=2 HTTP_PORT=9090 go run cmd/server/main.go
```

### Override in docker-compose.yml

```yaml
services:
  btc-aggregator:
    build: .
    ports:
      - "9090:9090"
    environment:
      - POLL_INTERVAL=5s
      - REQUEST_TIMEOUT=3s
      - MAX_RETRIES=2
      - HTTP_PORT=9090
    restart: unless-stopped
    stop_grace_period: 15s
```

## API Endpoints

```bash
curl http://localhost:8080/price
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

## Testing

```bash
go test ./... -v
go test ./... -cover
```
