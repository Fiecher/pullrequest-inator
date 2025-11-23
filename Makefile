.PHONY: all run stop deps lint test test-e2e generate migrate-create migrate-up migrate-down coverage help

APP_NAME = pullrequest-inator
DB_URL = postgres://postgres:password@localhost:5432/pullrequest?sslmode=disable
MIGRATE_PATH = database/migrations/pg
DOCKER_COMPOSE = docker-compose


all: deps lint build


run:
	$(DOCKER_COMPOSE) up --build -d
	@echo "Service is running at http://localhost:8080"

stop:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f $(APP_NAME)

deps:
	go mod tidy
	go mod verify

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

generate:
	oapi-codegen -package api -generate types,server,spec -o internal/api/codegen_api.go openapi.yaml
	@echo "API code generated."

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $$name

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

test-e2e:
	go test -v ./test/e2e

help:
	@echo "Available commands:"
	@echo "  make run            - Start application in Docker"
	@echo "  make stop           - Stop Docker containers"
	@echo "  make test-e2e       - Run integration tests"
	@echo "  make generate       - Generate Go code from OpenAPI"
	@echo "  make migrate-create - Create a new DB migration file"
	@echo "  make migrate-up     - Apply DB migrations (locally)"
	@echo "  make migrate-down   - Rollback last DB migration (locally)"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make coverage       - Show code coverage report"