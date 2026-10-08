# clinical-service: Пациенты и ЭМК

Самостоятельный Go-модуль. Все исходники, зависимости, миграции и Dockerfile находятся в этом каталоге; общий репозиторий для сборки не требуется.

```sh
go build ./cmd/clinical-service
go test -race ./...
go vet ./...
docker build -t his-clinical .
```

Перед запуском задайте `DATABASE_URL`, указывающий на собственную PostgreSQL-базу сервиса. Примените миграции `go run ./cmd/clinical-service migrate`, затем запустите `go run ./cmd/clinical-service`. Адрес HTTP-сервера задаётся `HTTP_ADDR` (по умолчанию `:8080`).
Для клинического сервиса также задаются `STAFF_URL` и `CATALOG_URL` — базовые HTTP-адреса справочных сервисов.


API: `/api/v1/patients`, `/api/v1/patients/{id}/record`, `/api/v1/records/{id}`. Спецификация — `api/openapi.json`, работающий сервис отдаёт её через `GET /openapi.json`. Диагностика: `GET /live`, `GET /ready`.

Интеграционные тесты выполняются при наличии `TEST_DATABASE_URL`: `go test -race -count=1 ./...`. Они используют собственные временные схемы и удаляют их после завершения. При отсутствии переменной выполняются доступные тесты без БД, PostgreSQL-тесты пропускаются.
