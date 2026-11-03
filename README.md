# Go Clean Architecture + EDA + Scheduler Boilerplate

Production-ready REST API, Event-Driven Consumer, and Cron Scheduler microservice boilerplate in Go using **Echo**, **Apache Kafka (`segmentio/kafka-go`)**, **Cron (`robfig/cron/v3`)**, **GORM (PostgreSQL)**, **Uber Dig** (Dependency Injection), and **Uber Zap** (Structured Logging).

## Project Structure (3 Loosely-Coupled Entrypoints)

```text
go-template/
├── cmd/
│   ├── api/                 # 1. HTTP REST API Server entrypoint
│   ├── worker/              # 2. Kafka Event Consumer Worker entrypoint
│   └── scheduler/           # 3. Cron Job Scheduler entrypoint (configured via YAML)
├── configs/                 # Environment YAML configuration files (supports ${ENV} expansion)
├── deployments/             # Dockerfile & docker-compose.yml (Postgres + Kafka + API + Worker + Scheduler)
├── docs/                    # Embedded OpenAPI 3.0 spec & Swagger UI (/swagger)
├── internal/
│   ├── controller/
│   │   ├── http/            # Inbound HTTP handlers, DTOs, & routes (used by cmd/api)
│   │   ├── event/           # Inbound Kafka Event Consumers (used by cmd/worker)
│   │   └── scheduler/       # Inbound Cron Jobs (used by cmd/scheduler)
│   ├── middleware/          # HTTP middlewares (JWT Auth, Zap Logger, Recover, RequestID)
│   ├── usecase/             # Business logic layer
│   └── repository/          # Data access & outbound messaging layer
│       ├── db/              # PostgreSQL GORM repositories (User, Auth, Outbox)
│       └── messaging/       # Outbound Kafka Producer & DLQ Publisher
├── migrations/              # SQL schema migrations (MST_USER, SESSION, OUTBOX_EVENTS)
├── pkg/                     # Shared, stateless packages
│   ├── config/              # YAML + Environment config loader (App, DB, JWT, Kafka, Scheduler)
│   ├── di/                  # Uber Dig DI container registration
│   ├── event/               # Standard CloudEvents JSON envelope
│   ├── hash/                # Bcrypt password hashing helpers
│   ├── logger/              # Uber Zap logger initialization
│   ├── resource/            # Database & Kafka connection initialization
│   ├── response/            # Standardized HTTP JSON response & error formatting
│   └── validator/           # Echo request payload validator
├── .env.example             # Sample environment variables
└── Makefile                 # Common developer commands
```

## Architecture & Separation of Concerns

```text
                       ┌──────────────────────────┐
                       │      CONFIG & LOGGER     │
                       └─────────────┬────────────┘
                                     ↓
                       ┌──────────────────────────┐
                       │         RESOURCE         │
                       │  (PostgreSQL + Kafka)    │
                       └─────────────┬────────────┘
                                     ↓
            ┌────────────────────────┴────────────────────────┐
            ↓                                                 ↓
┌───────────────────────┐                         ┌───────────────────────┐
│     DB REPOSITORY     │                         │  MESSAGING REPOSITORY │
│ (User, Auth, Outbox)  │                         │ (Kafka Producer/DLQ)  │
└───────────┬───────────┘                         └───────────┬───────────┘
            └────────────────────────┬────────────────────────┘
                                     ↓
                       ┌──────────────────────────┐
                       │         USECASE          │
                       │     (Business Logic)     │
                       └─────────────┬────────────┘
                                     ↓
     ┌───────────────────────────────┼───────────────────────────────┐
     ↓                               ↓                               ↓
┌────────────────────┐      ┌────────────────────┐      ┌────────────────────┐
│  HTTP CONTROLLER   │      │  EVENT CONTROLLER  │      │SCHEDULER CONTROLLER│
│   (cmd/api/main)   │      │  (cmd/worker/main) │      │(cmd/scheduler/main)│
└────────────────────┘      └────────────────────┘      └────────────────────┘
```

## Quick Start

1. **Copy environment variables**:
   ```bash
   cp .env.example .env
   ```
2. **Start full stack with Docker Compose**:
   ```bash
   make docker-up
   ```
3. **Run each service independently locally**:
   ```bash
   make run            # Starts HTTP API Server (:9001)
   make run-worker     # Starts Kafka Event Consumer Worker
   make run-scheduler  # Starts Cron Scheduler (from configs/*.yaml)
   ```
