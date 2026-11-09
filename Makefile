APP_NAME := go-template
API_PATH := ./cmd/api
GRPC_PATH := ./cmd/grpc
WORKER_PATH := ./cmd/worker
SCHEDULER_PATH := ./cmd/scheduler

.PHONY: help dev run run-grpc run-worker run-scheduler build build-grpc build-worker build-scheduler build-all test test-coverage lint mock tidy clean docker-up docker-down

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

dev: ## Run HTTP API server with live hot-reload using Air
	@which air > /dev/null || (echo "Air not found. Installing air..." && go install github.com/air-verse/air@latest)
	air -c .air.toml

run: ## Run the HTTP API server only (:9001)
	go run $(API_PATH)/main.go

run-grpc: ## Run the standalone gRPC server (:9002)
	go run $(GRPC_PATH)/main.go

run-worker: ## Run the standalone Kafka Event Consumer worker
	go run $(WORKER_PATH)/main.go

run-scheduler: ## Run the standalone Cron Scheduler
	go run $(SCHEDULER_PATH)/main.go

build: ## Build the API binary into ./bin/api
	go build -o ./bin/api $(API_PATH)

build-grpc: ## Build the gRPC binary into ./bin/grpc
	go build -o ./bin/grpc $(GRPC_PATH)

build-worker: ## Build the Worker binary into ./bin/worker
	go build -o ./bin/worker $(WORKER_PATH)

build-scheduler: ## Build the Scheduler binary into ./bin/scheduler
	go build -o ./bin/scheduler $(SCHEDULER_PATH)

build-all: build build-grpc build-worker build-scheduler ## Build all 4 binaries (api, grpc, worker, scheduler)

test: ## Run all unit tests
	go test -v $(shell go list -f '{{if len .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

test-coverage: ## Run all unit tests and generate HTML coverage report
	go test -v -coverprofile=coverage.out $(shell go list -f '{{if len .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint: ## Run golangci-lint
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

mock: ## Generate mocks using mockery
	go run github.com/vektra/mockery/v2@v2.46.3

tidy: ## Tidy go modules
	go mod tidy

clean: ## Remove build artifacts
	rm -rf ./bin ./tmp

docker-up: ## Start Postgres, Kafka, Kafka-UI, API, gRPC, Worker, and Scheduler
	docker compose -f deployments/docker-compose.yml up --build -d

docker-down: ## Stop Docker Compose services
	docker compose -f deployments/docker-compose.yml down
