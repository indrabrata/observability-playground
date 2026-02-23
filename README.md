# Observability Playground

A hands-on learning project for understanding the **three pillars of observability** : **Logs**, **Metrics**, and **Traces** — in a real Go backend system.

---

## Goal

Learn how to instrument a backend service and wire it into a full observability stack using open-source tools. The app itself is a simple Product CRUD API; the real focus is on how signals flow from the application into the monitoring stack.

---

## Tech Stack

| Layer             | Tool                                                                                          |
| ----------------- | --------------------------------------------------------------------------------------------- |
| Language          | Go 1.24                                                                                       |
| HTTP Router       | [chi](https://github.com/go-chi/chi)                                                          |
| Database          | SQLite                                                                                        |
| Query generation  | [sqlc](https://sqlc.dev)                                                                      |
| DB Migration      | [goose](https://github.com/pressly/goose)                                                     |
| Logs              | [Zap](https://github.com/uber-go/zap) + [Lumberjack](https://github.com/natefinch/lumberjack) |
| Metrics           | [Prometheus client](https://github.com/prometheus/client_golang)                              |
| Traces            | [OpenTelemetry Go SDK](https://opentelemetry.io/docs/languages/go/)                           |
| DB Tracing        | [otelsql](https://github.com/XSAM/otelsql) (OTel wrapper for SQLite driver)                   |
| Collector / Agent | [Grafana Alloy](https://grafana.com/oss/alloy/)                                               |
| Log storage       | [Loki](https://grafana.com/oss/loki/)                                                         |
| Trace storage     | [Tempo](https://grafana.com/oss/tempo/)                                                       |
| Metric storage    | [Prometheus](https://prometheus.io)                                                           |
| Visualization     | [Grafana](https://grafana.com/oss/grafana/)                                                   |
| API Docs          | Swagger (swaggo)                                                                              |

---

## System Design

```
┌────────────────────────────────────────────────────────────────┐
│                        Go Application                          │
│                                                                │
│   ┌──────────┐   ┌──────────┐   ┌─────────────────────────┐    │
│   │  Handler │──▶│ Service  │──▶│  Repository (sqlc)      │    │
│   └──────────┘   └──────────┘   └───────────┬─────────────┘    │
│        │               │                    │                  │
│   OTel Span       OTel Span          ┌──────▼──────┐           │
│        └───────────────┘             │   otelsql   │           │
│                  │                   │ (OTel spans │           │
│                  │                   │  for DB ops)│           │
│                  │                   └──────┬──────┘           │
│                  │                          │                  │
│                  │                   ┌──────▼──────┐           │
│                  │                   │   SQLite    │           │
│                  │                   └─────────────┘           │
│                  │                                             │
│          ┌───────▼────────┐                                    │
│          │  Middleware    │                                    │
│          │  - RequestID   │                                    │
│          │  - Metrics     │──▶ /metrics (Prometheus format)    │
│          │  - Request log │──▶ ./logs/<date>.log (JSON)        │
│          └────────────────┘                                    │
│                                                                │
│  Traces ──────────────────────────────────────▶ OTLP gRPC      │
└────────────────────────────────────────────────────┬───────────┘
                                                     │
                              ┌──────────────────────▼──────────────────────┐
                              │               Grafana Alloy                 │
                              │                                             │
                              │  ┌─────────────────┐  ┌──────────────────┐  │
                              │  │ loki.source.file │  │ otelcol.receiver │ │
                              │  │  (tail log files)│  │  .otlp (gRPC)    │ │
                              │  └────────┬─────────┘  └────────┬─────────┘ │
                              └───────────┼─────────────────────┼───────────┘
                                          │                     │
                              ┌───────────▼──────┐  ┌──────────▼──────────┐
                              │      Loki        │  │       Tempo         │
                              │  (log storage)   │  │   (trace storage)   │
                              └───────────┬──────┘  └──────────┬──────────┘
                                          │                     │
                              ┌───────────▼─────────────────────▼──────────┐
                              │                  Grafana                   │
                              │   Datasources: Prometheus · Loki · Tempo   │
                              └──────────────────────┬─────────────────────┘
                                                     │
                         ┌───────────────────────────┘
                         │
              ┌──────────▼──────────┐
              │     Prometheus       │
              │  scrapes /metrics    │◀──── Node Exporter (host metrics)
              └─────────────────────┘
```

### Signal Flow

| Pillar      | Path                                                                                                  |
| ----------- | ----------------------------------------------------------------------------------------------------- |
| **Logs**    | Zap writes JSON logs to `./logs/<date>.log` → Alloy tails the file → pushes to Loki                   |
| **Metrics** | Prometheus middleware records request count & latency → exposed at `/metrics` → Prometheus scrapes it |
| **Traces**  | OTel spans created in handler & service → exported via OTLP gRPC to Alloy → forwarded to Tempo        |

### Ports

| Service    | Port            | Purpose                            |
| ---------- | --------------- | ---------------------------------- |
| App        | `8080`          | REST API + `/metrics` + Swagger UI |
| Grafana    | `3000`          | Dashboards                         |
| Prometheus | `9090`          | Metric storage & query             |
| Loki       | `3100`          | Log storage                        |
| Tempo      | `3200` / `4318` | Trace storage (HTTP)               |
| Alloy      | `4317`          | OTLP gRPC receiver                 |
| Alloy      | `12345`         | Alloy debug UI                     |

---

## Project Structure

```
.
├── main.go                        # Wires everything together
├── handler/                       # HTTP handlers (OTel spans)
├── service/                       # Business logic (OTel spans)
├── repository/                    # sqlc-generated DB layer
├── middleware/
│   ├── metrics.go                 # Prometheus counter + histogram
│   ├── request.go                 # Structured request logging
│   └── request_id.go              # Injects X-Request-ID header
├── infrastructure/
│   ├── zap_log.go                 # Zap + Lumberjack setup (log rotation)
│   ├── prometheus_metric.go       # Prometheus registry setup
│   ├── open_telemetry_trace.go    # OTel TracerProvider → Alloy via gRPC
│   ├── sqlite3_db.go              # SQLite connection (otelsql-instrumented)
│   └── migration.go               # Goose auto-migration on startup
├── monitoring/
│   ├── alloy/config.alloy         # Alloy: tail logs → Loki, receive traces → Tempo
│   ├── prometheus/prometheus.yml  # Prometheus scrape config
│   ├── loki/config.yaml           # Loki storage config
│   ├── tempo/config.yaml          # Tempo storage config
│   ├── datasource/datasources.yaml# Grafana datasource provisioning
│   └── grafana/dashboard.json     # Pre-built Grafana dashboard
├── sql/
│   ├── migrations/                # Goose SQL migration files
│   └── queries/                   # sqlc query definitions
├── model/                         # Request / response structs
├── constant/                      # App name, package constants
├── common/                        # Shared utilities (response interceptor)
├── docker-compose.yaml            # Full stack: app + monitoring
├── Dockerfile
├── Makefile
└── sqlc.yaml
```

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) & Docker Compose
- Go 1.24+ (for local development only)

### Run the full stack

```bash
docker compose up -d
```

| URL                              | What's there            |
| -------------------------------- | ----------------------- |
| `http://localhost:8080/swagger/` | Swagger API docs        |
| `http://localhost:3000`          | Grafana (admin / admin) |
| `http://localhost:9090`          | Prometheus              |
| `http://localhost:12345`         | Alloy debug UI          |

### Run locally (dev mode)

```bash
cp .env.example .env
# Uncomment godotenv.Load() in main.go
make run
```

---

## API Endpoints

| Method   | Path             | Description         |
| -------- | ---------------- | ------------------- |
| `GET`    | `/health`        | Health check        |
| `GET`    | `/metrics`       | Prometheus metrics  |
| `POST`   | `/products`      | Create a product    |
| `GET`    | `/products`      | List all products   |
| `GET`    | `/products/{id}` | Get a product by ID |
| `PUT`    | `/products/{id}` | Update a product    |
| `DELETE` | `/products/{id}` | Delete a product    |
| `GET`    | `/swagger/*`     | Swagger UI          |

## Observability Details

### Logs (Zap → Loki)

- Structured JSON logs written to `./logs/<dd-MM-yyyy>.log`
- Log rotation via Lumberjack (max 1 GB per file, 30 backups, 90-day retention)
- Log level controlled by `LOG_LEVEL` env var (`DEBUG`, `INFO`, `WARN`, `ERROR`)
- Alloy tails the log directory and pushes entries to Loki with labels `job=observability-playground`

### Metrics (Prometheus)

Two custom metrics are recorded by `MetricsMiddleware` on every request:

| Metric            | Type      | Labels                         |
| ----------------- | --------- | ------------------------------ |
| `total_requests`  | Counter   | `method`, `endpoint`, `status` |
| `request_latency` | Histogram | `method`, `endpoint`           |

Prometheus scrapes `/metrics` on the app container. Node Exporter provides host-level system metrics.

### Traces (OpenTelemetry → Tempo)

- Spans are created at the **handler** and **service** layers
- DB queries are also traced via `otelsql` (wraps the SQLite driver)
- Exported over OTLP gRPC to Alloy (`:4317`), which forwards them to Tempo
- W3C TraceContext propagation is enabled for distributed tracing compatibility

---

## Environment Variables

See [`.env.example`](.env.example) for all available variables with inline descriptions.
