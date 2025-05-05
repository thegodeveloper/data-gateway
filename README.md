# 🔗 Data Gateway Service

A high-performance, scalable **Data Gateway** written in **Go**, designed to serve as a central access point for microservices to query multiple types of data sources such as **PostgreSQL**, **DynamoDB**, and **MongoDB**.

This service acts as a smart translation layer: microservices send simple query requests to the gateway, and the gateway routes and executes those queries against the appropriate backend — handling the complexity, connectivity, and scaling internally.

---

## ✨ Features

- ⚙️ **Clean Architecture** — Domain-driven, modular, and easy to extend.
- 📡 **Multiple Data Source Support**
    - PostgreSQL (via `database/sql`)
    - DynamoDB (via AWS SDK v2)
    - MongoDB (via official Go MongoDB driver)
- 📈 **Built for Scale** — Easily containerized, scalable via Kubernetes or ECS.
- 🚀 **Pluggable Design** — Add new databases or services in minutes.
- 🌍 **Single Unified Endpoint** — Microservices send structured JSON, no need for separate SDKs.

---

## 📦 Project Structure

```text
.
├── cmd/gateway/               # Main app entry point
├── internal/
│   ├── app/                   # Core business logic
│   ├── config/                # Configuration management
│   ├── datasource/
│   │   ├── dynamodb/
│   │   ├── mongodb/
│   │   └── postgres/
│   ├── domain/                # Interfaces and request models
│   └── transport/http/        # HTTP API layer
├── pkg/common/                # Logging utilities
└── go.mod / go.sum
```

## 🚀 Getting Started

### MongoDB Request

You can make a request like this:

```shell
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "source": "mongodb",
    "params": {
      "database": "orders",
      "collection": "users",
      "filter": { "age": { "$gt": 20 } }
    }
}'
```

### PostgreSQL Request

```shell
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "source": "postgres",
    "params": {
      "query": "SELECT id, name, email FROM users WHERE active = true;"
    }
}'
```

### Response Example

```json
[
  {
    "id": 1,
    "name": "Alice Doe",
    "email": "alice@example.com"
  },
  {
    "id": 2,
    "name": "Bob Smith",
    "email": "bob@example.com"
  }
]
```

## Implementing OpenTelemetry

To implement OpenTelemetry, add the following dependencies:

```shell
go get go.opentelemetry.io/otel@latest
go get go.opentelemetry.io/otel/sdk@latest
go get go.opentelemetry.io/otel/trace@latest
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp@latest
go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin@latest
```

### Update the Go module

```shell
go mod tidy
```

## 🗒️ Notes

- In the code example I have `/query` path, but you can change it to any path you want or you can add more paths for different data sources.