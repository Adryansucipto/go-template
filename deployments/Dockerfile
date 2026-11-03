# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/scheduler ./cmd/scheduler

# Final stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/api /app/bin/api
COPY --from=builder /app/bin/worker /app/bin/worker
COPY --from=builder /app/bin/scheduler /app/bin/scheduler
COPY --from=builder /app/configs /app/configs

EXPOSE 9001

ENTRYPOINT ["/app/bin/api"]
