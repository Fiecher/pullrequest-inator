# pullrequest-inator

- [README EN ver](https://github.com/Fiecher/pullrequest-inator/blob/main/README.en.md)

Микросервис, автоматически назначающий ревьюеров на Pull Request’ы (PR).
Позволяет управлять командами и участниками. 
Взаимодействие происходит исключительно через HTTP API.

Разработан в рамках тестового задания [Avito-Tech Autumn-2025](https://github.com/avito-tech/tech-internship/tree/main/Tech%20Internships/Backend/Backend-trainee-assignment-autumn-2025).

## Стэк

- Язык: [Go 1.24+](https://go.dev/)
- База данных: [PostgreSQL](https://www.postgresql.org/)
- Работа с БД: [pgx/v5](https://github.com/jackc/pgx)
- Миграции: [golang-migrate](https://github.com/golang-migrate/migrate)
- API Spec: [OpenAPI 3.0](https://swagger.io/specification/) + [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) 

## Быстрый старт

1. Клонирование репозитория:

```bash
git clone https://github.com/Fiecher/pullrequest-inator.git
```

2. Запуск сервиса (приложение + БД + миграции):

```bash
docker-compose up --build
```

3. Сервис доступен по адресу: `http://localhost:8080`

## Тестирование

В проекте реализованы E2E (End-to-End) тесты, которые поднимают отдельное окружение в Docker и проверяют реальные сценарии использования.

- Запуск [тестов](https://github.com/Fiecher/pullrequest-inator/tree/main/test/e2e):


```bash
make test-e2e
```

## Команды Makefile

Список команд [Makefile](https://github.com/Fiecher/pullrequest-inator/blob/main/Makefile):

| Команда                       | Категория | Описание                                                        |
|-------------------------------|-----------|-----------------------------------------------------------------|
| `make run`                    | Docker    | Запускает приложение и БД через Docker Compose.                 |
| `make stop`                   | Docker    | Останавливает и удаляет контейнеры.                             |
| `make test-e2e`               | Testing   | Запускает E2E-тесты (`./test/e2e`).                             |
| `make generate`               | CodeGen   | Генерирует Go-код на основе `openapi.yaml`.                     |
| `make fmt`                    | Linter    | Форматирует весь код проекта (`go fmt`).                        |
| `make lint`                   | Linter    | Запускает линтер (`golangci-lint`) для проверки стиля и ошибок. |
| `make deps`                   | Tools     | Обновляет зависимости (`go mod tidy`, `go mod verify`).         |
| `make install-tools`          | Tools     | Установка `oapi-codegen`, `golangci-lint`....                    |
| `make migrate-create`         | Database  | Создает новый пустой файл миграции.                             |
| `make migrate-up`             | Database  | Применяет все незавершенные миграции к локальной БД.            |
| `make migrate-down`           | Database  | Откатывает последнюю миграцию на один шаг назад.                |
