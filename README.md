# Проект на Go 1.25: распределённый key-value сервис с журналом транзакций (transaction log), реализованный как набор микросервисов.

Цели проекта:
- Практика микросервисной архитектуры (несколько независимых сервисов).
- Транзакционный журнал с постепенной реализацией улучшений.
- Набор практик продакшн-кода: JSON-логи, graceful shutdown, Docker, docker-compose, тесты, метрики Prometheus.

***

## Архитектура

Проект состоит из двух сервисов и общих библиотек:

- **libs/**
    - `logger` — обёртка над zerolog с единым форматом JSON-логов.
    - `txlog` — append-only журнал транзакций (log), используемый kv-service.
- **services/kv-service/**
    - HTTP API для операций с ключами: `/kv/set`, `/kv/get`, `/kv/delete`.
    - Хранит данные в памяти + записывает события в файловый журнал `kv.log`.
- **services/api-gateway/**
    - Внешний API для клиентов: `/api/set`, `/api/get`, `/api/delete`.
    - Проксирует запросы в kv-service, добавляет валидацию и свой слой логирования.

Взаимодействие:

1. Клиент обращается к `api-gateway` по HTTP.
2. api-gateway вызывает kv-service через внутренний HTTP-клиент.
3. kv-service обновляет in-memory store и записывает событие в `txlog`.
4. Все HTTP-сервисы экспортируют метрики Prometheus на `/metrics`.


## Запуск без Docker

### Требования

- Go 1.25+
- make (опционально)

### Тесты и базовые команды

```
make test      # go test ./...
make fmt       # gofmt по проекту
make tidy      # go mod tidy

В одном терминале:

go run ./services/kv-service/cmd/kv
# kv-service слушает :8081

Во втором терминале:

go run ./services/api-gateway/cmd/api
# api-gateway слушает :8080
```

### Пример запросов через curl
```
 Установить значение
curl -s -X POST http://localhost:8080/api/set \
-H "Content-Type: application/json" \
-d '{"key":"user42","value":"Alice"}'

 Прочитать значение
curl -s "http://localhost:8080/api/get?key=user42"

 Удалить ключ
curl -s -X DELETE "http://localhost:8080/api/delete?key=user42"
```

***

## Запуск с Docker + Prometheus

В проекте есть `docker-compose.yml`, поднимающий:

- `kv-service` (порт 8081)
- `api-gateway` (порт 8080)
- `prometheus` (порт 9090), который скрейпит `/metrics` обоих сервисов.

Конфигурация Prometheus описана в `prometheus.yml`.

### Запуск

docker-compose up --build

### После старта:

api-gateway: http://localhost:8080

kv-service: http://localhost:8081

Prometheus: http://localhost:9090


### Проверка метрик:
```
 api-gateway
curl -s http://localhost:8080/metrics | grep http_requests_total

 kv-service
curl -s http://localhost:8081/metrics | grep http_requests_total
```


***

## Observability и особенности реализации

### Метрики Prometheus

Оба сервиса экспортируют:

- `/metrics` — стандартный endpoint Prometheus client_golang.
- Счётчик HTTP-запросов:
  - `http_requests_total{handler="api_set",method="POST",status="200"}` и др.

(При желании можно добавить histogram по длительности запросов.)

### Журнал транзакций (txlog)

Библиотека `libs/txlog`:

- Формат записей:
  - `Op len(key) len(value) keyBytes valueBytes "\n"`
- Ограничения размеров:
  - `MaxKeySize` и `MaxValueSize`, ошибки `ErrKeyTooLarge`, `ErrValueTooLarge`.
- Безопасное закрытие:
  - `Sync()` + `Close()` перед shutdown.
- Простая compaction:
  - `CompactLogFile(path)` переписывает лог, оставляя только последние `set` для живых ключей.

### Graceful shutdown

Оба сервиса:

- Слушают сигналы `SIGINT`/`SIGTERM`.
- Вызывают `http.Server.Shutdown` с таймаутом.
- Корректно закрывают `txlog.FileLog` (kv-service) и логируют результат.

***

## Базовый Makefile

## Makefile

Минимально полезные цели:

```makefile
PROJECT_NAME=go-txlog-microservices

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make test    - run unit tests"
	@echo "  make tidy    - run go mod tidy"
	@echo "  make fmt     - run gofmt on all go files"

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
```
***

## Структура проекта

```text
.
├── libs/
│   ├── logger/                # Общий JSON-логгер (zerolog-обёртка)
│   └── txlog/                 # Журнал транзакций (append-only log)
│       ├── txlog.go
│       └── txlog_test.go
└── services/
    ├── kv-service/            # Внутренний key-value сервис
    │   ├── cmd/kv/            # Точка входа (main.go)
    │   ├── internal/
    │   │   ├── http/          # HTTP-хендлеры: /kv/set, /kv/get, /kv/delete
    │   │   ├── metrics/       # Prometheus-метрики kv-service
    │   │   ├── server/        # Конструктор http.Server
    │   │   └── store/         # In-memory хранилище + работа с txlog
    │   └── Dockerfile
    └── api-gateway/           # Внешний API для клиентов
        ├── cmd/api/           # Точка входа (main.go)
        ├── internal/
        │   ├── client/        # HTTP-клиент для общения с kv-service
        │   ├── http/          # HTTP-хендлеры: /api/set, /api/get, /api/delete
        │   ├── metrics/       # Prometheus-метрики api-gateway
        │   └── server/        # Конструктор http.Server
        ├── test/apigateway_test/
        │   └── e2e_api_kv_test.go  # End-to-end тест через реальный HTTP
        └── Dockerfile
```
