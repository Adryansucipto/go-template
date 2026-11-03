APP_NAME := go-template
API_PATH := ./cmd/api
WORKER_PATH := ./cmd/worker
SCHEDULER_PATH := ./cmd/scheduler

.PHONY: help run run-worker run-scheduler build build-worker build-scheduler build-all test tidy clean docker-up docker-down

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

run: ## Run the HTTP API server only
	go run $(API_PATH)/main.go

run-worker: ## Run the standalone Kafka Event Consumer worker
	go run $(WORKER_PATH)/main.go

run-scheduler: ## Run the standalone Cron Scheduler
	go run $(SCHEDULER_PATH)/main.go

build: ## Build the API binary into ./bin/api
	go build -o ./bin/api $(API_PATH)

build-worker: ## Build the Worker binary into ./bin/worker
	go build -o ./bin/worker $(WORKER_PATH)

build-scheduler: ## Build the Scheduler binary into ./bin/scheduler
	go build -o ./bin/scheduler $(SCHEDULER_PATH)

build-all: build build-worker build-scheduler ## Build all 3 binaries (api, worker, scheduler)

test: ## Run all unit tests
	go test -v -race ./...

tidy: ## Tidy go modules
	go mod tidy

clean: ## Remove build artifacts
	rm -rf ./bin ./tmp

docker-up: ## Start Postgres, Kafka, Kafka-UI, API, Worker, and Scheduler using Docker Compose
	docker compose -f deployments/docker-compose.yml up --build -d

docker-down: ## Stop Docker Compose services
	docker compose -f deployments/docker-compose.yml down
