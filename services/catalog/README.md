# catalog-service: Справочник диагнозов

Самостоятельный Go-модуль. Все исходники, зависимости, миграции и Dockerfile находятся в этом каталоге; общий репозиторий для сборки не требуется.

```sh
go build ./cmd/catalog-service
go test -race ./...
go vet ./...
docker build -t his-catalog .
```

Перед запуском задайте `DATABASE_URL`, указывающий на собственную PostgreSQL-базу сервиса. Примените миграции `go run ./cmd/catalog-service migrate`, затем запустите `go run ./cmd/catalog-service`. Адрес HTTP-сервера задаётся `HTTP_ADDR` (по умолчанию `:8080`).

API: `/api/v1/diagnoses` и `/api/v1/diagnoses/seed`. Спецификация — `api/openapi.json`, работающий сервис отдаёт её через `GET /openapi.json`. Диагностика: `GET /live`, `GET /ready`.

Интеграционные тесты выполняются при наличии `TEST_DATABASE_URL`: `go test -race -count=1 ./...`. Они используют собственные временные схемы и удаляют их после завершения. При отсутствии переменной выполняются доступные тесты без БД, PostgreSQL-тесты пропускаются.
