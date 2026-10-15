# clinical-service: Пациенты и ЭМК

Самостоятельный Go-модуль. Все исходники, зависимости, миграции и Dockerfile находятся в этом каталоге; общий репозиторий для сборки не требуется.

```sh
go build ./cmd/clinical-service
go test -race ./...
go vet ./...
docker build -t his-clinical .
```

Перед запуском задайте `DATABASE_URL`, указывающий на собственную PostgreSQL-базу сервиса. Примените миграции `go run ./cmd/clinical-service migrate`, затем запустите `go run ./cmd/clinical-service`. Адрес HTTP-сервера задаётся `HTTP_ADDR` (по умолчанию `:8080`).
Для клинического сервиса обязателен `GATEWAY_URL` — базовый HTTP(S)-адрес Nginx. Клиенты медработников и диагнозов используют этот общий gateway; адреса конечных сервисов находятся в `deploy/nginx.conf`. В Compose значение по умолчанию — `http://nginx`; его можно переопределить окружением без пересборки. При отсутствии или неверном адресе запуск завершается ошибкой конфигурации.


API: `/api/v1/patients`, `/api/v1/patients/{id}/record`, `/api/v1/records/{id}`. Источник контракта в репозитории — `contracts/clinical.yaml`; `api/openapi.yaml` — локальная копия для самостоятельной генерации; JSON-документация — сгенерированный `api/openapi.json`, работающий сервис отдаёт её через `GET /openapi.json`. Диагностика: `GET /live`, `GET /ready`.

Интеграционные тесты выполняются при наличии `TEST_DATABASE_URL`: `go test -race -count=1 ./...`. Они используют собственные временные схемы и удаляют их после завершения. При отсутствии переменной выполняются доступные тесты без БД, PostgreSQL-тесты пропускаются.

Контракт изменяется в корневом `contracts/`, затем из корня выполняется `make generate-contracts`. Для самостоятельной генерации из сохранённой локальной копии выполните `go generate ./api`. `api/generated.gen.go` содержит DTO, типизированные request/response, `StrictServerInterface` и gorilla/mux adapter. Сгенерированные файлы не редактируются вручную; `internal/app/api` реализует интерфейс и преобразует DTO в бизнес-модели.
