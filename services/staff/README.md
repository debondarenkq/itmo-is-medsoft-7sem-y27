# staff-service: Медработники

Самостоятельный Go-модуль. Все исходники, зависимости, миграции и Dockerfile находятся в этом каталоге; общий репозиторий для сборки не требуется.

```sh
go build ./cmd/staff-service
go test -race ./...
go vet ./...
docker build -t his-staff .
```

Перед запуском задайте `DATABASE_URL`, указывающий на собственную PostgreSQL-базу сервиса. Примените миграции `go run ./cmd/staff-service migrate`, затем запустите `go run ./cmd/staff-service`. Адрес HTTP-сервера задаётся `HTTP_ADDR` (по умолчанию `:8080`).

API: `/api/v1/staff`. Источник контракта — `api/openapi.yaml`; JSON-документация — сгенерированный `api/openapi.json`, работающий сервис отдаёт её через `GET /openapi.json`. Диагностика: `GET /live`, `GET /ready`.

Интеграционные тесты выполняются при наличии `TEST_DATABASE_URL`: `go test -race -count=1 ./...`. Они используют собственные временные схемы и удаляют их после завершения. При отсутствии переменной выполняются доступные тесты без БД, PostgreSQL-тесты пропускаются.

После изменения контракта выполните `go generate ./api`. `api/generated.gen.go` содержит DTO, типизированные request/response, `StrictServerInterface` и gorilla/mux adapter. Сгенерированные файлы не редактируются вручную; `internal/app/api` реализует интерфейс и преобразует DTO в бизнес-модели.
