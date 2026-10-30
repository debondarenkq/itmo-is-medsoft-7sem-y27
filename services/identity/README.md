# identity-service: аккаунты и безопасность

Самостоятельный Go-модуль со своей PostgreSQL-базой и миграциями Goose. Владеет аккаунтами, хэшами паролей, TOTP, refresh-токенами и аудитом. Проверяет права остальных сервисов через внутренний HTTP API.

```sh
go build ./cmd/identity-service
go test -race ./...
go vet ./...
docker build -t his-identity .
```

Для запуска нужны `DATABASE_URL`, `STAFF_URL`, `SIGNING_KEY_FILE`, `SERVICE_KEY_FILE`; `HTTP_ADDR` по умолчанию `:8080`. Миграции выполняются при старте. В runtime доступа к БД staff нет. Приватные ключи не входят в образ.

Публичные маршруты: `/api/v1/auth`, `/api/v1/accounts`, `/api/v1/audit`. Документация: `/security/openapi.json`. Внутренние `/internal/authorize`, `/internal/audit`, `/internal/staff` требуют сервисный ключ и не публикуются Nginx. Проверки состояния: `/live`, `/ready`.

Источник публичного контракта — `contracts/security.yaml`; локальная копия и генерация — `api/security`, `go generate ./api/security`. Из корня репозитория используйте `make generate-contracts`. Внутренние операции синхронизации описаны в `contracts/internal-identity.yaml`.

Тесты PostgreSQL требуют `TEST_DATABASE_URL`, создают и удаляют отдельные схемы. Для перехода со старой версии используйте `make migrate-identity` из корня проекта. Подкоманда `identity-service migrate` предназначена только для офлайн-переноса и дополнительно требует `LEGACY_DATABASE_URL`; она не сбрасывает существующие аккаунты.
