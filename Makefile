PROJECT_NAME := go-txlog-microservices

DOCKER_COMPOSE := docker compose

.PHONY: help test fmt tidy lint build run-api run-kv docker-build up down logs

# Цель по умолчанию: если запустить просто "make", выполнится "help".
.DEFAULT_GOAL := help


# ---------- Справка ----------

# help — выводит список доступных целей и краткое описание каждой.
help:
	@echo "Available make targets:"
	@echo "  make test          - run all Go unit tests"
	@echo "  make fmt           - run gofmt on all Go files"
	@echo "  make tidy          - run go mod tidy"
	@echo "  make build         - build api-gateway and kv-service binaries (local)"
	@echo "  make run-api       - run api-gateway locally"
	@echo "  make run-kv        - run kv-service locally"
	@echo "  make docker-build  - build Docker images via docker compose"
	@echo "  make up            - start all services via docker compose"
	@echo "  make down          - stop all services via docker compose"
	@echo "  make logs          - show docker compose logs"


# ---------- Go команды ----------

test:
	go test ./...

fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

tidy:
	go mod tidy

build:
	mkdir -p bin
	go build -o bin/api-gateway ./services/api-gateway/cmd/api
	go build -o bin/kv-service ./services/kv-service/cmd/kv


# ---------- Локальный запуск без Docker ----------

run-api:
	go run ./services/api-gateway/cmd/api

run-kv:
	go run ./services/kv-service/cmd/kv


# ---------- Docker / Docker Compose ----------

docker-build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f