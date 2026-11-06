# Go Enterprise Microservice Boilerplate

Standar *boilerplate* microservice Go kelas industri berbasis **Clean Architecture & Event-Driven Architecture (EDA)** sesuai standar arsitektur microservice Blibli. Didesain untuk performa tinggi, skalabilitas independen di Kubernetes, keandalan pengiriman pesan (*guaranteed at-least-once delivery*), dan isolasi tanggung jawab yang ketat (*Strict Layer Boundaries*).

Teknologi utama yang digunakan:
- **Framework API**: [Echo v4](https://echo.labstack.com/) (REST HTTP) & [gRPC](https://grpc.io/)
- **Messaging & EDA**: [Apache Kafka (`segmentio/kafka-go`)](https://github.com/segmentio/kafka-go) dengan Transactional Event Queue & Dead Letter Queue (DLQ)
- **Database & Cache**: [GORM (PostgreSQL)](https://gorm.io/) & [Go-Redis v9](https://github.com/redis/go-redis) (Cache-Aside Pattern)
- **Job Scheduler**: [Cron v3 (`robfig/cron/v3`)](https://github.com/robfig/cron) (Jadwal deklaratif via YAML)
- **Dependency Injection**: [Uber Dig](https://github.com/uber-go/dig)
- **Structured Logging & Tracing**: [Uber Zap](https://github.com/uber-go/zap) dengan `X-Request-ID` Context Flow
- **Observability**: Prometheus Metrics (`/metrics`) & Kubernetes Probes (`/health/liveness`, `/health/readiness`)

Seluruh struktur direktori dan penamaan layer di dalam template ini bersifat **baku (standardized)** sehingga dapat langsung diduplikasi ke repository microservice baru tanpa membuang atau mengubah hirarki folder.

---

## 1. Aturan Baku 4-Tier Layer (Strict Layer Boundaries)

Microservice ini membagi kode ke dalam 4 layer utama di dalam folder `internal/` yang saling terisolasi:

```text
               ┌────────────────────────────────────────────────────────┐
               │              LAYER 1: CONTROLLER (INBOUND)            │
               │   • internal/controller/http (Echo REST API)          │
               │   • internal/controller/grpc (gRPC Service Server)    │
               └───────────────────────────┬────────────────────────────┘
                                           │ memanggil
                                           ▼
               ┌────────────────────────────────────────────────────────┐
               │              LAYER 2: SERVICE (CORE DOMAIN)           │
               │   • internal/service/<domain> (Logika Bisnis)         │
               │   • internal/service/event/producer (Kirim Kafka/DLQ) │
               │   • internal/service/event/consumer (Terima Kafka)    │
               │   • internal/service/scheduler (Cron Jobs)            │
               └───────────────────────────┬────────────────────────────┘
                                           │
                        ┌──────────────────┴──────────────────┐
                        ▼                                     ▼
        ┌───────────────────────────────┐     ┌───────────────────────────────┐
        │  LAYER 3: REPOSITORY (DB)     │     │  LAYER 4: OUTBOUND (EXTERNAL) │
        │  • internal/repository/<dom>  │     │  • internal/outbound/cache    │
        │  • internal/repository/tx     │     │  • internal/outbound/http     │
        │  • internal/repository/queue  │     │  • internal/outbound/grpc     │
        └───────────────┬───────────────┘     └───────────────┬───────────────┘
                        ▼                                     ▼
                  [PostgreSQL]                       [Redis / External API]
```

### Ringkasan Tanggung Jawab & Larangan per Layer

| Layer (`internal/`) | Tanggung Jawab Utama | Aturan Ketat (Strict Boundary Rules) |
| :--- | :--- | :--- |
| **`controller`** | **Murni Inbound API** (HTTP REST & gRPC Server). Menerima request eksternal, validasi input DTO, dan meneruskan ke `service`. | ❌ **Dilarang** menaruh Kafka Consumer, Cron Scheduler, atau query Database di dalam `controller`. |
| **`service`** | **Pusat Logika Bisnis, EDA Event, & Scheduler**. Mengorkestrasi aturan bisnis, memanggil `repository` untuk data, `outbound` untuk cache/API eksternal, dan mengelola penerimaan/pengiriman event Kafka. | 💡 Satu-satunya layer yang boleh mengoordinasikan `repository` dan `outbound`. |
| **`repository`** | **Murni Akses Database (PostgreSQL)**. Berisi query GORM, `transaction.Manager` (`WithTransaction`), dan tabel antrean event (`event_queue`). | ❌ **Dilarang** memanggil Redis, Kafka, atau HTTP/gRPC Client di dalam `repository`. |
| **`outbound`** | **Murni Integrasi Keluar & Caching**. Berisi Redis Cache (`outbound/cache`), Cross-Service REST Client (`outbound/http`), dan Cross-Service gRPC Client (`outbound/grpc`). | 💡 Setara levelnya dengan `repository`, hanya boleh dipanggil oleh `service`. |

---

## 2. Diagram Cara Kerja Template

### Diagram 1: Arsitektur 4 Entrypoint & Alur Komponen (Flowchart)

Aplikasi memiliki 4 entrypoint independen di `cmd/` yang dapat dijalankan sebagai pod/container terpisah di Kubernetes:

```mermaid
flowchart TD
    subgraph Clients ["Klien & Dunia Luar"]
        BROWSER["Web / Mobile Apps"]
        K8S["Kubernetes Ingress / Gateway"]
        PROMETHEUS["Prometheus Scraper"]
    end

    subgraph Entrypoints ["4 Entrypoints Independen (cmd/)"]
        CMD_API["cmd/api\n(HTTP Server :9001)"]
        CMD_GRPC["cmd/grpc\n(gRPC Server :9002)"]
        CMD_WORKER["cmd/worker\n(Kafka Consumer)"]
        CMD_SCHED["cmd/scheduler\n(Cron Jobs)"]
    end

    subgraph ControllerLayer ["Layer 1: Controller (internal/controller - Inbound API)"]
        CTRL_HTTP["Echo HTTP Controller\n(REST Routes, Validation)"]
        CTRL_GRPC["gRPC Service Server\n(Protobuf Endpoints)"]
        OBS_ROUTES["Observability Routes\n(/metrics, /health/*)"]
    end

    subgraph ServiceLayer ["Layer 2: Service (internal/service - Business & EDA)"]
        SVC_USER["User & Auth Business Service"]
        SVC_PROD["Event Producer Service\n(Kafka Writer + DLQ)"]
        SVC_CONS["Event Consumer Service\n(Worker Loop + Retry)"]
        SVC_SCHED["Cron Scheduler Service\n(Event Queue Relay)"]
    end

    subgraph StorageLayer ["Layer 3: Repository & Layer 4: Outbound"]
        REPO_DB["Repository (PostgreSQL)\nUser, Auth, Event Queue, TxManager"]
        OUT_CACHE["Outbound Cache (Redis)\nCache-Aside Pattern"]
        OUT_HTTP["Outbound HTTP Client\nPartner REST APIs"]
        OUT_GRPC["Outbound gRPC Client\nPartner gRPC Services"]
    end

    subgraph ExternalSystems ["Sistem Eksternal & Infrastruktur"]
        POSTGRES[("PostgreSQL 16")]
        REDIS[("Redis 7 (Cache)")]
        KAFKA[["Apache Kafka / Redpanda"]]
        PARTNER_SVC[["Partner Microservices"]]
    end

    %% Client traffic
    BROWSER --> CMD_API
    K8S --> CMD_GRPC
    PROMETHEUS --> CMD_API

    %% Entrypoints routing
    CMD_API --> CTRL_HTTP
    CMD_API --> OBS_ROUTES
    CMD_GRPC --> CTRL_GRPC
    CMD_WORKER --> SVC_CONS
    CMD_SCHED --> SVC_SCHED

    %% Controller to Service
    CTRL_HTTP --> SVC_USER
    CTRL_GRPC --> SVC_USER

    %% Service to Storage & Outbound
    SVC_USER --> REPO_DB
    SVC_USER --> OUT_CACHE
    SVC_USER --> OUT_HTTP
    SVC_USER --> OUT_GRPC
    SVC_USER --> SVC_PROD

    %% Worker and Scheduler Flow
    SVC_CONS --> SVC_USER
    SVC_CONS --> SVC_PROD
    SVC_SCHED --> REPO_DB
    SVC_SCHED --> SVC_PROD

    %% Storage to Infra
    REPO_DB --> POSTGRES
    OUT_CACHE --> REDIS
    SVC_PROD --> KAFKA
    KAFKA -.-> SVC_CONS
    OUT_HTTP --> PARTNER_SVC
    OUT_GRPC --> PARTNER_SVC
```

---

### Diagram 2: Sequence Diagram — Inbound API, Cache-Aside & Transactional Event Queue

Diagram ini mengilustrasikan bagaimana request baca (Read) dan tulis (Write) diproses secara atomik tanpa risiko *dual-write failure*:

```mermaid
sequenceDiagram
    autonumber
    actor Client as Klien / Frontend
    participant Ctrl as Controller (Echo REST)
    participant Svc as Service (User Service)
    participant OutCache as Outbound (Redis Cache)
    participant Repo as Repository (PostgreSQL + Tx)
    participant Sched as Service (Cron Scheduler)
    participant Prod as Service (Event Producer)
    participant Kafka as Apache Kafka ("user.events")

    Note over Client,OutCache: ALUR BACA (Read Flow - Cache-Aside Pattern)
    Client->>Ctrl: GET /v1/users (Headers: X-Request-ID, Bearer JWT)
    Ctrl->>Svc: GetUsers(ctx)
    Svc->>OutCache: GetUsersCache(ctx)
    alt Cache Hit
        OutCache-->>Svc: Kembalikan data dari Redis
    else Cache Miss
        OutCache-->>Svc: Nilai cache nil / miss
        Svc->>Repo: GetUserRepository(ctx)
        Repo-->>Svc: Query SELECT dari PostgreSQL
        Svc->>OutCache: SetUsersCache(ctx, users, TTL: 15m)
    end
    Svc-->>Ctrl: Data Users DTO
    Ctrl-->>Client: HTTP 200 OK (JSON)

    Note over Client,Kafka: ALUR TULIS (Write Flow - Atomic DB Transaction & Event Queue)
    Client->>Ctrl: POST /v1/register (Payload: Name, Username, Password)
    Ctrl->>Svc: CreateUser(ctx, payload)
    Svc->>Repo: TxManager.WithTransaction(ctx, fn)
    activate Repo
    Repo->>Repo: 1. INSERT INTO "MST_USER" (Data Bisnis)
    Repo->>Repo: 2. INSERT INTO "EVENT_QUEUE" (Status: PENDING)
    Repo-->>Svc: Commit Transaksi (Keduanya sukses bersamaan)
    deactivate Repo
    Svc->>OutCache: InvalidateUsersCache(ctx)
    Svc-->>Ctrl: User Response DTO
    Ctrl-->>Client: HTTP 201 Created (JSON)

    Note over Sched,Kafka: RELAY ASINKRON (Background Cron Scheduler)
    loop Setiap interval (configs/dev.yaml)
        Sched->>Repo: FetchPendingEvents(ctx, batchSize: 50)
        Repo-->>Sched: Daftar event berstatus PENDING
        Sched->>Prod: PublishRaw(ctx, topic, payload)
        Prod->>Kafka: Kirim pesan ke topic "user.events"
        Sched->>Repo: MarkProcessed(ctx, eventID) (Status: PROCESSED)
    end
```

---

### Diagram 3: Sequence Diagram — Kafka Event Consumer (Retry & Dead Letter Queue - DLQ)

Diagram ini mengilustrasikan worker penerima event yang dilengkapi dengan mekanisme *retry backoff* dan penerusan otomatis ke DLQ jika batas kegagalan terlampaui:

```mermaid
sequenceDiagram
    autonumber
    participant Kafka as Kafka ("user.events")
    participant Cons as Service (Event Consumer Worker)
    participant Svc as Service (Business Logic)
    participant Outbound as Outbound (HTTP / gRPC Client)
    participant Partner as Partner Microservice
    participant Prod as Service (Event Producer)
    participant DLQ as Kafka DLQ ("user.events.dlq")

    Kafka->>Cons: FetchMessage(ctx)
    loop Retry up to MaxRetries (dengan Exponential Backoff)
        Cons->>Svc: HandleUserCreated(ctx, eventPayload)
        Svc->>Outbound: FetchData / SyncToPartner(ctx)
        Outbound->>Partner: Request Antar-Service (+ X-API-Key)
        alt Partner Sukses (HTTP 200 / gRPC OK)
            Partner-->>Outbound: Respon Sukses
            Outbound-->>Svc: Hasil
            Svc-->>Cons: Selesai tanpa error
        else Partner Gagal / Timeout
            Partner-->>Outbound: Error (HTTP 5xx / Connection Error)
            Outbound-->>Svc: Error
            Svc-->>Cons: Error (picu retry loop)
        end
    end

    alt Pemrosesan Berhasil
        Cons->>Kafka: CommitMessages(ctx, msg)
    else Batas Retry Terlampaui (Max Retries Exhausted)
        Note over Cons,DLQ: Pindahkan pesan rusak ke Dead Letter Queue (DLQ)
        Cons->>Prod: PublishDLQ(ctx, dlqTopic, msg.Value, reason)
        Prod->>DLQ: Tulis ke "user.events.dlq" + Header "x-error-reason"
        Cons->>Kafka: CommitMessages(ctx, msg) (Agar antrean tidak macet)
    end
```

---

### Diagram 4: Observability, Health Checks & Tracing Flow

```mermaid
flowchart LR
    subgraph Monitoring ["Monitoring & Orkestrasi"]
        K8S_KUBELET["Kubernetes Kubelet\n(Liveness & Readiness)"]
        PROM_SERVER["Prometheus Server\n(Metrics Collector)"]
    end

    subgraph AppServer ["cmd/api (:9001)"]
        MIDDLEWARE["Echo Zap Middleware\n(Extracts & Injects X-Request-ID)"]
        
        subgraph Endpoints ["Observability Endpoints"]
            LIVE["GET /health/liveness\n(Process is Running -> 200 OK)"]
            READY["GET /health/readiness\n(Pings DB & Redis -> 200/503)"]
            METRIC["GET /metrics\n(Latency, Request Count, Memory, GC)"]
        end
    end

    subgraph Dependencies ["Dependensi"]
        PG[("PostgreSQL")]
        RD[("Redis")]
    end

    K8S_KUBELET -- "Setiap 10s" --> LIVE
    K8S_KUBELET -- "Setiap 5s" --> READY
    READY -.-> PG
    READY -.-> RD

    PROM_SERVER -- "Scrape setiap 15s" --> METRIC
    MIDDLEWARE --> METRIC
```

---

## 3. Struktur Direktori Proyek

```text
go-template/
├── .air.toml                    # Konfigurasi Live Hot-Reload untuk development lokal
├── .github/workflows/           # CI/CD GitHub Actions (Test, Coverage, Build Check)
├── .golangci.yml                # Konfigurasi linter Go berstandar industri
├── .mockery.yaml                # Konfigurasi generator mock otomatis (mockery)
├── cmd/                         # 4 Binary Entrypoint Independen
│   ├── api/main.go              # 1. HTTP REST API Server (:9001)
│   ├── grpc/main.go             # 2. Standalone gRPC Server (:9002)
│   ├── worker/main.go           # 3. Kafka Event Consumer Worker
│   └── scheduler/main.go        # 4. Cron Job Scheduler (jadwal via configs/*.yaml)
├── configs/                     # File Konfigurasi Deklaratif (dev.yaml, prod.yaml)
├── deployments/                 # Dockerfile (Multi-stage) & docker-compose.yml
├── docs/                        # Dokumentasi API (OpenAPI 3.0 & Swagger UI di /swagger)
├── internal/                    # Layer Internal Aplikasi (Private)
│   ├── controller/              # [LAYER 1] Inbound Controller (HTTP & gRPC)
│   │   ├── http/                # Routes Echo, Handler Auth & User
│   │   ├── grpc/                # gRPC Server Handler & Interceptor
│   │   └── model/               # Request & Response DTOs
│   ├── service/                 # [LAYER 2] Core Business Logic, EDA & Cron
│   │   ├── auth/                # Business logic registrasi & autentikasi
│   │   ├── user/                # Business logic pengelolaan user
│   │   ├── event/               # Producer & Consumer Kafka (Retry + DLQ)
│   │   └── scheduler/           # Cron runner untuk Event Queue relay
│   ├── repository/              # [LAYER 3] Akses Database PostgreSQL Murni
│   │   ├── transaction/         # GORM Transaction Manager (WithTransaction)
│   │   ├── model/               # GORM Database Entities
│   │   ├── auth/                # Repository Auth DB
│   │   ├── user/                # Repository User DB
│   │   └── event_queue/         # Repository Tabel EVENT_QUEUE (Transactional Outbox)
│   ├── outbound/                # [LAYER 4] Integrasi Eksternal & Caching
│   │   ├── cache/               # Outbound Redis (Cache-Aside Pattern)
│   │   ├── http/                # Outbound REST Client ke service lain
│   │   └── grpc/                # Outbound gRPC Client ke service lain
│   └── middleware/              # Zap Logger, JWT Auth, API Key, Prometheus
├── migrations/                  # File DDL SQL PostgreSQL murni
├── pkg/                         # Pustaka & Resource Bersama (Shared Resources)
│   ├── config/                  # Parser konfigurasi YAML & Env
│   ├── di/                      # Container Dependency Injection (Uber Dig)
│   ├── event/                   # Format Envelope CloudEvents JSON
│   ├── hash/                    # Helper Bcrypt password hashing
│   ├── logger/                  # Uber Zap Logger dengan Context Tracing
│   ├── resource/                # Resource Pool: PostgreSQL, Redis, Kafka, HTTP/gRPC Client
│   ├── response/                # Standar Response JSON API
│   └── validator/               # Struct validator payload HTTP
├── proto/                       # Kontrak Protocol Buffers (.proto) & Generated Go code
├── Makefile                     # Perintah otomatisasi build, run, test, dan linter
└── README.md                    # Dokumentasi utama proyek
```

---

## 4. Panduan Menjalankan Aplikasi

### Persiapan Awal
Salin environment template dan unduh seluruh dependency Go:

```bash
cp .env.example .env
make tidy
```

### Opsi A: Menjalankan Seluruh Stack via Docker Compose
Menjalankan PostgreSQL, Redis, Kafka (Redpanda), Kafka-UI (`:8080`), HTTP API (`:9001`), gRPC (`:9002`), Worker, dan Scheduler sekaligus:

```bash
make docker-up
```
Untuk menghentikan:

```bash
make docker-down
```

### Opsi B: Menjalankan Secara Independen di Lokal

| Perintah | Deskripsi | URL / Port |
| :--- | :--- | :--- |
| `make dev` | Menjalankan HTTP API dengan **Live Hot-Reload** (Air) | `http://localhost:9001` |
| `make run` | Menjalankan HTTP REST API biasa | `http://localhost:9001` |
| `make run-grpc` | Menjalankan server gRPC standalone | `localhost:9002` |
| `make run-worker` | Menjalankan Kafka Event Consumer Worker | Background CLI |
| `make run-scheduler`| Menjalankan Cron Scheduler Job | Background CLI |

---

## 5. Testing, Code Quality & Build

### Unit Testing & Coverage Report
Seluruh unit test diletakkan berdampingan dengan filenya (*colocated unit test*) dengan pengujian mock di setiap layer:
```bash
make test           # Menjalankan seluruh unit test
make test-coverage  # Menjalankan unit test & membuat visual report coverage.html (Coverage: 84.3%)
```
Buka file `coverage.html` di browser untuk melihat baris kode mana saja yang telah terlindungi test.

### Static Code Analysis (Linter)
Memeriksa kualitas kode Go, unhandled errors, dan potensi *memory leaks*:
```bash
make lint
```
### Build Binary Produksi
Meng-compile keempat microservice menjadi file binary mandiri di dalam folder `./bin/`:
```bash
make build-all
```
File yang dihasilkan:
- `./bin/api` (HTTP REST API Server)
- `./bin/grpc` (gRPC Server)
- `./bin/worker` (Kafka Consumer Worker)
- `./bin/scheduler` (Cron Job Scheduler)

---

## 6. Observability, Metrik & Health Probes

Microservice ini sudah terintegrasi dengan standar monitoring modern:

1. **Prometheus Metrics**:
   - `GET http://localhost:9001/metrics`
   - Mengekspos histogram latensi request HTTP, status code counter, jumlah goroutine, thread, dan alokasi memori GC.
2. **Kubernetes Liveness Probe**:
   - `GET http://localhost:9001/health/liveness`
   - Mengembalikan HTTP `200 OK` selama proses aplikasi berjalan. Digunakan K8s untuk me-restart pod jika terjadi deadlock.
3. **Kubernetes Readiness Probe**:
   - `GET http://localhost:9001/health/readiness`
   - Melakukan live-ping serentak ke PostgreSQL dan Redis. Mengembalikan HTTP `200 OK` jika semua siap, atau HTTP `503 Service Unavailable` jika ada dependensi yang bermasalah.
4. **Distributed Tracing Header (`X-Request-ID`)**:
   - Setiap request HTTP secara otomatis diekstrak header `X-Request-ID`-nya (atau dibuatkan UUID baru jika belum ada).
   - ID tersebut disuntikkan ke dalam `context.Context` dan dicatat pada setiap baris log Uber Zap melalui helper `logger.FromContext(ctx, log)`.
5. **Interactive Swagger UI**:
   - Buka `http://localhost:9001/swagger` untuk menguji API secara langsung dari browser.

---

## 7. Glosarium Penamaan Baku

| Istilah Baku | Arti / Konsep Asal | Mengapa Dinamakan Demikian? |
| :--- | :--- | :--- |
| **`repository/transaction`** | Database Transaction (`BEGIN`, `COMMIT`, `ROLLBACK`) | Menggantikan singkatan teknis `tx`. Menandakan fungsi wrapper `WithTransaction` untuk menjamin konsistensi ACID beberapa query. |
| **`repository/event_queue`** | Transactional Outbox Pattern | Menggantikan istilah teknis `outbox`. Menegaskan bahwa ini adalah tabel antrean event (`EVENT_QUEUE`) di database yang menunggu dikirimkan ke Kafka oleh scheduler. |
| **`service/scheduler/event_queue`** | Outbox Relay Job | Job cron yang bertugas membaca baris `PENDING` di `EVENT_QUEUE` lalu meneruskannya ke Kafka topic. |
| **`service/event/producer`** | Message Producer | Modul pengirim pesan ke broker Kafka (+ pengiriman ke topic DLQ). |
| **`service/event/consumer`** | Message Consumer Worker | Modul penerima pesan dari broker Kafka dengan logika retry & dead letter queue. |
| **`outbound/*`** | Outbound Adapter (Hexagonal Architecture) | Mengelompokkan semua integrasi ke luar sistem (Redis Cache, REST HTTP Client ke partner, gRPC Client ke partner). Berada di level layer yang setara dengan `repository`. |
