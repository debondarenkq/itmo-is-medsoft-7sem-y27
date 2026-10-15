# Учебная госпитальная информационная система

Backend на Go: три самостоятельных микросервиса, PostgreSQL и единая точка входа Nginx. Нативный клиент реализован на Swift 6 и SwiftUI для macOS 14+: интерфейсы администратора и медработника, ЭМК, история и импорт справочника. Все запросы клиента проходят через Nginx.

## Запуск

Нужны Docker с Docker Compose. Из корня проекта:

```sh
docker compose up --build -d --wait
```

API доступен по `http://localhost:8080/api/v1`. Compose создаёт три базы с отдельными пользователями, выполняет SQL-миграции Goose отдельными задачами и запускает сервисы. Данные сохраняются в Docker volume между перезапусками. Новая база содержит пустой справочник диагнозов. Записи добавляются вручную или импортируются из выбранного JSON-файла через `POST /api/v1/diagnoses/import`. Пример файла — `examples/diagnoses.json`; он не загружается автоматически.

Если установлен отдельный `docker-compose`, можно использовать `./scripts/compose.sh up --build -d --wait`. Скрипт сам выбирает доступную команду. На macOS/Linux также доступен `make up`.

Остановка с сохранением данных:

```sh
docker compose down
```

Для выбора HTTP-порта измените опубликованный порт Nginx в `compose.yaml`. PostgreSQL и сервисы наружу по умолчанию не публикуются. Используются только синтетические учебные данные, без авторизации и TLS согласно заданию.

## Структура

```text
contracts/        исходные OpenAPI-схемы всех сервисов
services/
  staff/          медработники
  catalog/        справочник диагнозов
  clinical/       пациенты, ЭМК, назначения, история
    cmd/clinical-service/
    internal/app/api/         HTTP-обработчики
    internal/models/          модели, запросы и обычные ошибки Go
    internal/config/          конфигурация зависимостей
    internal/pkg/modules/     бизнес-логика и интерфейсы зависимостей
    internal/pkg/storage/     PostgreSQL
    internal/pkg/adapters/    HTTP-клиенты справочников
    internal/app/server/      запуск HTTP-сервера, БД и миграций
    internal/validation/      проверки UUID и текстовых полей
    api/                      OpenAPI-контракт и сгенерированные DTO
    migrations/               SQL-миграции Goose
    go.mod
    go.sum
    Dockerfile
macos/             основной нативный клиент Swift + SwiftUI
desktop/           предыдущий Flutter-клиент, сохранён для сравнения
deploy/            Nginx, инициализация и тестовая конфигурация PostgreSQL
examples/          примеры файлов для ручного импорта
scripts/           запуск Compose и сквозные проверки API
docs/              архитектура и описание API
```

Каждый сервис имеет собственные `go.mod`, `go.sum`, модели, утилиты, миграции и Dockerfile. Общих Go-пакетов, `go.work`, локальных `replace` и импортов другого сервиса нет. Для сборки достаточно его каталога:

```sh
cd services/clinical
go build ./cmd/clinical-service
docker build -t his-clinical .
```

Миграции включены в исполняемый файл. Для самостоятельного запуска задайте `DATABASE_URL`, выполните `go run ./cmd/clinical-service migrate`, затем `go run ./cmd/clinical-service`. Для остальных сервисов используется соответствующее имя команды. На Windows переменные окружения задаются средствами PowerShell.

Конфигурация:

| Переменная | Где нужна | Значение |
|---|---|---|
| `DATABASE_URL` | Все сервисы | Обязательный адрес собственной PostgreSQL-базы |
| `HTTP_ADDR` | Все сервисы | По умолчанию `:8080` |
| `GATEWAY_URL` | Клинический сервис | Обязательный URL Nginx; в Compose — `http://nginx`, можно переопределить окружением |

Клинический сервис обращается к справочникам через Nginx: `clinical → nginx → staff/catalog`. Адреса конечных сервисов находятся в upstream-блоках `deploy/nginx.conf`. При смене адреса staff/catalog обновляется конфигурация Nginx; клинический сервис не нужно менять или пересобирать. Смену IP контейнера за прежним DNS-именем Nginx отслеживает через встроенный DNS Docker.

Если меняется адрес самого gateway, задайте `GATEWAY_URL` и пересоздайте клинический контейнер без пересборки:

```sh
GATEWAY_URL=http://gateway.internal docker compose up -d --force-recreate clinical-service
```

## API и история ЭМК

[Описание запросов и команд](docs/api.md) и [архитектура](docs/architecture.md).

Спецификации работающего приложения:

- `GET /docs/staff/openapi.json`
- `GET /docs/catalog/openapi.json`
- `GET /docs/clinical/openapi.json`

Каждый сервис также отдаёт `/openapi.json`, `/live` и `/ready` на своём порту. По gateway `/live` проверяет сам Nginx; готовность отдельных сервисов проверяется healthchecks Compose.

При создании или изменении ЭМК обязателен `staff_id`. Блокировка редактирования реализована только в UI. Сохранение карты атомарно: команды, текущее состояние и история записываются одной транзакцией. При изменении назначения создаются события удаления старого и добавления нового. История хранит полные снимки данных и позволяет восстановить карту по времени или версии. Пациента с ЭМК нельзя удалить даже прямым вызовом API.

## Проверки

Требуется Go 1.26 или новее. Каждый модуль проверяется отдельно; из корня на macOS/Linux:

```sh
make test
make check
make integration
```

`make integration` запускает отдельный PostgreSQL на `127.0.0.1:55433`, выполняет тесты с `-race` и останавливает тестовый Compose-проект. Каждый тест создаёт собственную схему, применяет миграции и удаляет схему после завершения. Проверяются конкурентные операции, ограничения БД, мягкое удаление, загрузка диагнозов, атомарность и восстановление ЭМК. Обычный `make test` запускает unit-тесты; тесты PostgreSQL пропускаются, если не задан `TEST_DATABASE_URL`.

Для отдельного сервиса можно передать `TEST_DATABASE_URL` и выполнить `go test -race -count=1 ./...` внутри его каталога. Используйте отдельную тестовую базу, доступную пользователю с правом создания схем.

Сквозная проверка через Nginx после запуска backend:

```sh
python3 scripts/smoke.py
```

Она выполняет сценарии задания, включая CRUD, изменения ЭМК, исторические снимки, восстановление и ошибки API. Проверка создаёт учебные записи, в том числе ЭМК; карты намеренно не удаляются через API. Для чистой проверки используйте отдельный Compose-проект.

## Приложение для macOS

Инструкция — [macos/README.md](macos/README.md). После запуска backend выполните:

```sh
make macos-build
open macos/build/HIS.app
```

Нужны macOS 14+ и Xcode 16+. Для запуска из Xcode откройте `macos/HIS.xcodeproj`, выберите **HIS → My Mac** и нажмите **⌘R**. Адрес Nginx задаётся в меню HIS → Настройки… (`⌘,`). Workflow [macOS](https://github.com/debondarenkq/itmo-is-medsoft-7sem-y27/actions/workflows/desktop.yml) собирает `.app` и публикует ZIP. Старый Flutter-клиент сохранён в `desktop/` для сравнения.

## Контракты

HTTP API описывается в корневом каталоге [`contracts/`](contracts/README.md). Локальные схемы в сервисах обновляются генератором и не редактируются вручную. Go DTO, интерфейсы ручек, binding для gorilla/mux и HTTP-клиенты генерируются через `oapi-codegen`. SwiftUI-приложение использует сгенерированный Swift-пакет `macos/Generated` для единственного gateway. Подробности — [docs/contracts.md](docs/contracts.md). После изменения YAML выполните `make generate-contracts`; CI проверяет отсутствие расхождений между схемами и кодом.
