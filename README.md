# TaskScheduler (Golang + Redis + Asynq)

A backend service to schedule and process delayed and recurring tasks using [Asynq](https://github.com/hibiken/asynq) with Redis.  
Features:
- Delayed & recurring tasks
- Built-in retries and exponential backoff
- Structured logging with zerolog
- Prometheus metrics endpoint (for monitoring)
- Docker + docker-compose for quick local setup
- Example task handlers and graceful shutdown

## Quick start (local, development)

Requirements:
- Go 1.21+
- Redis
- Docker (optional)

1. Start Redis (using docker-compose):
```bash
docker-compose up -d
```

2. Build and run:
```bash
go mod tidy
go run ./cmd/server
```

3. In another terminal, run the worker:
```bash
go run ./cmd/worker
```

4. Use the included CLI or HTTP endpoints (future) to enqueue tasks. See `tasks/example` for usage.

## Project structure

```
cmd/
  server/      # runs scheduler + metrics + enqueuer demo
  worker/      # runs Asynq worker to process tasks
internal/
  config/      # configuration via viper
  tasks/       # task definitions and handlers
pkg/
  asynqclient/ # helpers to create client/server
Dockerfile
docker-compose.yml
README.md
```

## Example: enqueue a task
See `cmd/server/main.go` for a simple demo that enqueues:
- a one-off task delayed by 30s
- schedules a recurring task (every minute) using Asynq Scheduler
