# Проект на Go 1.25: распределённый key-value сервис с журналом транзакций (transaction log), реализованный как набор микросервисов.

Цели проекта:
- Практика микросервисной архитектуры (несколько независимых сервисов).
- Транзакционный журнал с постепенной реализацией улучшений.
- Набор практик продакшн-кода: JSON-логи, graceful shutdown, Docker, docker-compose, тесты.

***

## Базовый Makefile

`Makefile` в корне (минимально полезный набор).

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
