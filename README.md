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
