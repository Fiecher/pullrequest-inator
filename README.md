# pullrequest-inator

Микросервис, автоматически назначающий ревьюеров на Pull Request’ы (PR).
Позволяет управлять командами и участниками. 
Взаимодействие происходит исключительно через HTTP API.

Разработан в рамках [тестового задания](https://github.com/avito-tech/tech-internship/tree/main/Tech%20Internships/Backend/Backend-trainee-assignment-autumn-2025)

## Стэк

- Язык: [Go 1.24+](https://go.dev/)
- База данных: [PostgreSQL](https://www.postgresql.org/)
- Работа с БД: [pgx/v5](https://github.com/jackc/pgx)
- Миграции: [golang-migrate](https://github.com/golang-migrate/migrate)
- API Spec: [OpenAPI](https://swagger.io/specification/) 3.0 + [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) 

## Команды Makefile

- `make generate-api` ([Makefile](https://github.com/alphameo/pr-reviewnager/blob/main/Makefile))

## Быстрый старт

Для запуска требуется только Docker и Docker Compose.

1. Клонируйте репозиторий:

```bash
git clone https://github.com/Fiecher/pullrequest-inator.git
```

2. Запустите сервис (приложение + БД + миграции):

```bash
docker-compose up --build
```

3. Сервис доступен по адресу: `http://localhost:8080`
   - Health Check: `GET /health`

## Тестирование

В проекте реализованы E2E (End-to-End) тесты, которые поднимают отдельное окружение в Docker и проверяют реальные сценарии использования.

- Запуск тестов:


```bash
make test-e2e
```

## Команды Makefile

| Команда              | Категория | Описание                                                                 |
|----------------------|-----------|--------------------------------------------------------------------------|
| `make run`           | Docker    | Запускает приложение и БД через Docker Compose.                          |
| `make stop`          | Docker    | Останавливает и удаляет контейнеры.                                      |
| `make generate`      | CodeGen   | Генерирует Go код (`types`, `server`, `spec`) на основе `openapi.yaml`.  |
| `make deps`          | Tools     | Обновляет зависимости (`go mod tidy`, `go mod verify`).                  |
| `make fmt`           | Linter    | Форматирует весь код проекта (`go fmt`).                                 |
| `make lint`          | Linter    | Запускает линтер (`golangci-lint`) для проверки стиля и ошибок.          |
| `make test`          | Testing   | Запускает только Unit-тесты (`./internal/...`).                          |
| `make test-e2e`      | Testing   | Запускает E2E-тесты (`./test/e2e`).                                      |
| `make covergge`      | Testing   | Собирает отчет о покрытии кода и открывает его в браузере.               |
| `make migrate-create`| Database  | Создает новый пустой файл миграции.                                      |
| `make migrate-up`    | Database  | Применяет все незавершенные миграции к локальной БД.                     |
| `make migrate-down`  | Database  | Откатывает последнюю миграцию на один шаг назад.                         |
