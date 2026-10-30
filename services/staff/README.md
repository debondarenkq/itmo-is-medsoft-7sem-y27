# staff-service: Медработники

Самостоятельный Go-модуль. Все исходники, зависимости, миграции и Dockerfile находятся в этом каталоге; общий репозиторий для сборки не требуется.

```sh
go build ./cmd/staff-service
go test -race ./...
go vet ./...
docker build -t his-staff .
```

Перед запуском задайте `DATABASE_URL`, указывающий на собственную PostgreSQL-базу сервиса. Задайте также `AUTH_URL` (identity-service) и `SERVICE_KEY_FILE`, затем запустите `go run ./cmd/staff-service`; миграции применяются при старте. При обновлении существующей лабораторной сначала выполните `make migrate-identity` из корня проекта. Адрес HTTP-сервера задаётся `HTTP_ADDR` (по умолчанию `:8080`).

API: `/api/v1/staff`. Источник контракта в репозитории — `contracts/staff.yaml`; `api/openapi.yaml` — локальная копия для самостоятельной генерации; JSON-документация — сгенерированный `api/openapi.json`, работающий сервис отдаёт её через `GET /openapi.json`. Диагностика: `GET /live`, `GET /ready`.

Интеграционные тесты выполняются при наличии `TEST_DATABASE_URL`: `go test -race -count=1 ./...`. Они используют собственные временные схемы и удаляют их после завершения. При отсутствии переменной выполняются доступные тесты без БД, PostgreSQL-тесты пропускаются.

Контракт изменяется в корневом `contracts/`, затем из корня выполняется `make generate-contracts`. Для самостоятельной генерации из сохранённой локальной копии выполните `go generate ./api`. `api/generated.gen.go` содержит DTO, типизированные request/response, `StrictServerInterface` и gorilla/mux adapter. Сгенерированные файлы не редактируются вручную; `internal/app/api` реализует интерфейс и преобразует DTO в бизнес-модели.
