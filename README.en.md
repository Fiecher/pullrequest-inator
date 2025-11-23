# pullrequest-inator

- [README RU ver](https://github.com/Fiecher/pullrequest-inator/blob/main/README.en.md)

A microservice designed for automatically assigning reviewers to Pull Requests (PRs).
It provides functionality for managing teams and team members.
Interaction with the service is exclusively via an HTTP API.

Developed as part of the technical assignment for
the [Avito-Tech Autumn-2025](https://github.com/avito-tech/tech-internship/tree/main/Tech%20Internships/Backend/Backend-trainee-assignment-autumn-2025)
internship.

## Tech Stack

- Language: [Go 1.24+](https://go.dev/)
- Database: [PostgreSQL](https://www.postgresql.org/)
- DB Driver: [pgx/v5](https://github.com/jackc/pgx)
- Migrations: [golang-migrate](https://github.com/golang-migrate/migrate)
- API
  Spec: [OpenAPI 3.0](https://swagger.io/specification/) + [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/Fiecher/pullrequest-inator.git
```

2. Start the service (application + DB + migrations):

```bash
docker-compose up --build
```

3. The service will be available at: `http://localhost:8080`

## Тестирование

The project includes E2E (End-to-End) tests, which spin up a separate Docker environment to verify real-world usage
scenarios.

- Run the [tests](https://github.com/Fiecher/pullrequest-inator/tree/main/test/e2e):

```bash
make test-e2e
```

## Makefile Commands

List of commands in the [Makefile](https://github.com/Fiecher/pullrequest-inator/blob/main/Makefile):

| Команда               | Категория | Описание                                                            |
|-----------------------|-----------|---------------------------------------------------------------------|
| `make run`            | Docker    | Starts the application and DB using Docker Compose.                 |
| `make stop`           | Docker    | Stops and removes the containers.                                   |
| `make test-e2e`       | Testing   | Runs the E2E tests (`./test/e2e`).                                  |
| `make generate`       | CodeGen   | Generates Go code based on `openapi.yaml`.                          |
| `make fmt`            | Linter    | Formats all project code (`go fmt`).                                |
| `make lint`           | Linter    | Runs the linter (`golangci-lint`)  to check for style and errors.   |
| `make deps`           | Tools     | Updates dependencies (`go mod tidy`, `go mod verify`).              |
| `make install-tools`  | Tools     | Installs necessary tools like `oapi-codegen`, `golangci-lint`, etc. |
| `make migrate-create` | Database  | Creates a new empty migration file.                                 |
| `make migrate-up`     | Database  | Applies all pending migrations to the local DB.                     |
| `make migrate-down`   | Database  | Rolls back the last migration step.                                 |
