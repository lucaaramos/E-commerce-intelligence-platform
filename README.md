# E-Commerce Intelligence & Automation Platform

An internal platform to manage and optimize an e-commerce dropshipping business.
It centralizes products, suppliers, prices, stock and competitors to answer
actionable questions like *"what is the real cost?"*, *"which supplier is
cheapest?"* and *"what margin will I get?"*.

The system evolves from a simple MVP toward an intelligence center, following a
**useful → automated → intelligent → scalable** path.

## Status

| Iteration | Scope | Status |
|---|---|---|
| 1 | Bootstrap: docker-compose, config, Mongo/RabbitMQ wiring, health checks, logging, middleware, skeleton | ✅ Done |
| 2 | Domain models + repositories (products, suppliers, supplier_products, price_history, stock_history) | Pending |
| 3 | Services + REST API v1, pricing service, price/stock events | Pending |
| 4 | Alerts + Dashboard aggregation | Pending |
| 5 | Frontend: Vite + React + TypeScript + Tailwind | Pending |
| 6 | Worker + RabbitMQ consumers | Pending |
| 7 | Docs, OpenAPI, seeds, Makefile polish | Pending |

## Stack

- **Backend**: Go (Gin), REST API, MongoDB (mongo-driver v2), RabbitMQ (amqp091-go)
- **Frontend**: React, TypeScript, Vite, Tailwind CSS
- **Infra**: Docker, Docker Compose

## Quick start

```bash
cp .env.example .env       # adjust secrets
make up                    # docker compose up -d --build
```

Services:

| Service  | URL                                    |
|----------|----------------------------------------|
| Backend  | http://localhost:8080                  |
| Frontend | http://localhost:5173                  |
| Mongo    | mongodb://localhost:27017 (ecommerce_intel) |
| RabbitMQ | http://localhost:15672 (guest/guest)   |

Health:

```bash
curl http://localhost:8080/healthz   # liveness
curl http://localhost:8080/readyz    # readiness (mongodb + rabbitmq)
```

## Development

```bash
make backend-test   # go test ./...  (backend)
make frontend-build # npm run build  (frontend)
make logs           # follow container logs
make down           # stop the stack
make clean          # stop and remove volumes
```

Run the backend locally against the dockerized Mongo/RabbitMQ:

```bash
cd backend && go run ./cmd/server
```

## Repository layout

```
backend/
  cmd/server/          # application entrypoint
  internal/            # domain modules (products, suppliers, ...)
  pkg/
    config/            # env-driven configuration
    logging/           # slog structured logging
    mongodb/           # Mongo client + health
    rabbitmq/          # RabbitMQ client + health
    httpx/             # gin router, middlewares, responses, health endpoints
frontend/
  src/                 # React + TypeScript application
docs/
  architecture.md      # design decisions and MongoDB schema
```

## Security

- Secrets are never hardcoded; they come from environment variables (see `.env.example`).
- JWT auth is added in a later iteration (single internal admin user).
- Rate limiting per IP is already active.

## License

Internal project.
