# Architecture

## Overview

The platform is a **modular monolith**: one Go binary organized into domain
modules. Each module follows a lightweight layered structure:

```
HTTP request
   └── middlewares (request ID, logging, recover, CORS, rate limit)
        └── module handler   → validate HTTP, map status codes
             └── module service → business logic, events, orchestration
                  └── module repository → MongoDB queries
                       └── MongoDB
```

Infrastructure concerns live in `pkg/` (config, logging, mongodb, rabbitmq,
httpx). Domain modules live in `internal/` and depend only on `pkg/`.

## Current iteration (1) scope

- Docker Compose stack: `backend`, `frontend`, `mongodb`, `rabbitmq`
- Env-driven configuration with safe defaults
- MongoDB + RabbitMQ connections with ping-based health checks
- Liveness (`/healthz`) and readiness (`/readyz`) endpoints
- Structured logging (slog) with request IDs and per-request logs
- Gin middlewares: request ID, request logging, panic recovery, CORS, in-memory
  per-IP rate limiting
- Stable JSON envelopes: `{"data": ...}` for success, `{"error": {...}}` for errors
- Frontend placeholder (React + TS + Vite) served behind nginx, which proxies
  `/api/` to the backend

## Design decisions

1. **Modular monolith** — KISS/YAGNI. Module boundaries are explicit so modules
   can be split into services later only if growth justifies it.
2. **Denormalized computed fields** — fields such as `realCost` and `margin`
   will be stored on `products` and recalculated at write time so dashboards
   never run heavy aggregations on every read.
3. **Soft delete for products** — product history must survive for analysis.
4. **RabbitMQ used sparingly** — wired and health-checked now; consumers are
   added when asynchronous flows (alerts, analytics) need them.
5. **Mocks-only tests** — repositories are thin interfaces; unit tests use
   manual fakes. Real-Mongo verification happens manually against the compose
   stack (see `scripts/` in later iterations).

## MongoDB schema (planned)

Referenced documents are used for relationships (many-to-many junction
`supplier_products`, time series `price_history`/`stock_history`). Embedded
documents are used when data only makes sense within the parent (`variants`,
`order.items`). Full details land with iteration 2.

| Collection          | Key fields                                    |
|---------------------|-----------------------------------------------|
| `users`             | email (unique), passwordHash, role            |
| `products`          | sku (unique), status, sellingPrice, variants  |
| `suppliers`         | name (unique), country, integrationMethod     |
| `supplier_products` | productId+supplierId (unique), cost, shipping |
| `price_history`     | supplierProductId+timestamp (indexed)         |
| `stock_history`     | supplierProductId+timestamp (indexed)         |
| `alerts`            | type, severity, entityId, status, createdAt   |

## Security

- Secrets via environment variables only (`.env.example` documents them).
- JWT auth for a single internal admin user lands in iteration 3.
- Rate limiting, input validation and restricted CORS are active from the start.

## Iteration plan

1. Bootstrap (this iteration)
2. Domain models + repositories
3. Services + REST API v1 + pricing + events
4. Alerts + Dashboard
5. Frontend application
6. Background worker + RabbitMQ consumers
7. Docs, OpenAPI, seeds, polish
