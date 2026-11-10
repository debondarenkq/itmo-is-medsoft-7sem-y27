# identity-service: аккаунты и безопасность

Самостоятельный Go-модуль со своей PostgreSQL-базой и миграциями Goose. Владеет аккаунтами, хэшами паролей, TOTP, refresh-токенами и аудитом. Проверяет права остальных сервисов через внутренний HTTP API.

```sh
go build ./cmd/identity-service
go test -race ./...
go vet ./...
docker build -t his-identity .
```

Для запуска нужны `DATABASE_URL`, `STAFF_URL`, `SIGNING_KEY_FILE`, `SERVICE_KEY_FILE`; `HTTP_ADDR` по умолчанию `:8080`. Миграции выполняются при старте. В runtime доступа к БД staff нет. Приватные ключи не входят в образ.

Публичные маршруты: `/api/v1/auth`, `/api/v1/accounts`, `/api/v1/audit`. Документация: `/security/openapi.json`. Внутренние `/internal/authorize` и `/internal/audit` требуют сервисный ключ и не публикуются Nginx. Справочник `/internal/staff` обслуживает staff-service. Проверки состояния: `/live`, `/ready`.

Источник публичного контракта — `contracts/security.yaml`; локальная копия и генерация — `api/security`, `go generate ./api/security`. Из корня репозитория используйте `make generate-contracts`. Авторизация и доступ к справочнику staff описаны в `contracts/internal-identity.yaml`; общий аудит и enum событий — в `contracts/audit.yaml`. Сервер аудита генерируется через `go generate ./api/audit`, входные модели — через `go generate ./internal/models`.

Тесты PostgreSQL требуют `TEST_DATABASE_URL`, создают и удаляют отдельные схемы. Единственная команда в `cmd/` запускает HTTP-сервис. Аккаунты создаются через публичный API; демонстрационное заполнение выполняет внешний `scripts/create-demo-accounts.py`.

Профили сотрудников не реплицируются: ФИО, роль и удаление принадлежат staff. Вход, refresh и каждая авторизация сотрудника проверяют его текущее состояние прямым HTTP-запросом. Смена роли из редактора аккаунта делегируется staff вне транзакции identity. Публичный контракт и оба приложения сохраняются.
