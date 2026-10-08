# Контракты и генерация кода

Источник HTTP-контракта каждого сервиса — его `api/openapi.yaml` в формате OpenAPI 3.0.3. Поля запросов/ответов, типы, UUID, календарные даты, timestamp, enum, обязательность полей, ограничения и маршруты описываются здесь.

| Сервис | Источник | Сгенерированный код |
|---|---|---|
| Медработники | `services/staff/api/openapi.yaml` | `services/staff/api/generated.gen.go` |
| Диагнозы | `services/catalog/api/openapi.yaml` | `services/catalog/api/generated.gen.go` |
| Пациенты и ЭМК | `services/clinical/api/openapi.yaml` | `services/clinical/api/generated.gen.go` |

`openapi.json` рядом с YAML — сгенерированная JSON-документация. Сервис отдаёт её через `/openapi.json`; вручную этот файл не редактируется. Копия схемы также встраивается генератором в Go-код.

## Go-обработчики

`oapi-codegen` v2.8.0 генерирует:

- структуры DTO и enum;
- `RequestObject` с типизированными path/query/body;
- `ResponseObject` и типы допустимых успешных/ошибочных ответов;
- `StrictServerInterface`, который реализует API-слой;
- адаптер маршрутов и JSON-binding для `gorilla/mux`;
- HTTP-клиенты для описанных методов.

Например, ID пациента определяется контрактом как `type: string`, `format: uuid` и превращается в `uuid.UUID`:

```go
func (i *Implementation) GetPatient(
    ctx context.Context,
    request contract.GetPatientRequestObject,
) (contract.GetPatientResponseObject, error) {
    if err := platform.RequireID(request.ID); err != nil {
        return nil, err
    }
    patient, err := i.Module.GetPatient(ctx, request.ID.String())
    if err != nil {
        return nil, err
    }
    return contract.GetPatient200JSONResponse(toPatient(patient)), nil
}
```

ID и query-параметры разбирает сгенерированный adapter; тело декодируется в `request.Body`. В ручках нет `mux.Vars`, ручного `strconv` и `json.Decoder`. `RequireID` проверяет только бизнес-ограничение «UUID не нулевой», а не парсит строку.

Перед binding middleware проверяет запрос по схеме: обязательные поля, UUID/даты, enum, длины, диапазоны, массивы и неизвестные поля. Бизнес-модуль отдельно проверяет смысл операций: активность медработника, соответствие полей типу команды, наличие карты, конфликт версии и ограничения удаления.

В `mapping.go` транспортные DTO явно преобразуются в бизнес-модели и обратно. Бизнес-модели и данные БД не заменяются HTTP-структурами. История сохраняет собственные долговечные снимки.

## Межсервисные вызовы

Клинический сервис использует сгенерированные HTTP-клиенты:

- `internal/pkg/adapters/staffapi`: контракт медработников;
- `internal/pkg/adapters/catalogapi`: контракт справочника.

YAML этих клиентов копируется из контракта владельца при общей генерации, а Go-код генерируется внутри клинического сервиса. Это локальные копии внешних контрактов, которые не редактируются отдельно. Благодаря этому сервис собирается из своего каталога без импортов соседних Go-модулей, общего `go.mod`, `go.work` или `replace`. Межсервисный таймаут остаётся 3 секунды.

## Flutter

`desktop/api/openapi.json` — объединённое описание публичного gateway, автоматически собранное из трёх сгенерированных спецификаций. Это производный файл, а не отдельный редактируемый контракт.

OpenAPI Generator v7.26.0 создаёт пакет `desktop/packages/his_api` с DTO, enum, сериализацией и HTTP-методами. В клиенте один `contract.ApiClient` для одного адреса Nginx; `StaffApi`, `DiagnosesApi`, `PatientsApi`, `RecordsApi` — группы публичных ресурсов, не адреса микросервисов.

Формы создают типизированные `StaffInput`, `DiagnosisInput`, `PatientInput`, а черновик ЭМК — `Command`. Ручной сборки URL и JSON для обычных операций нет. UI-модели отдельно отображают полученные контрактные объекты.

Для необязательных полей используется `Optional`: отсутствие поля отличается от явно переданного null. Это важно для команд ЭМК, чтобы не отправлять `id: null` в команде добавления назначения. При импорте сохраняется исходный JSON выбранного файла, включая неизвестные поля, чтобы backend мог проверить весь документ без скрытого удаления данных сериализатором.

## Изменение и проверка

Редактируйте YAML владельца сервиса, затем из корня:

```sh
make generate-contracts
```

Нужны Go 1.26, Docker и Python 3. Генератор Go закреплён как Go tool в каждом `go.mod`; версия генератора Dart закреплена в скрипте. Если локального Dart нет, форматирование выполняется через SDK в Docker.

Для отдельного сервиса можно выполнить `go generate ./api` внутри его каталога. Для клинического сервиса внешние клиенты также генерируются через `go generate ./internal/pkg/adapters/staffapi ./internal/pkg/adapters/catalogapi` из сохранённых локальных схем.

Сгенерированные файлы сохраняются в Git и не правятся вручную. После генерации проверьте изменённые mapping и бизнес-логику, затем выполните тесты. При добавлении или изменении методов компилятор проверяет соответствие `StrictServerInterface`. CI workflow Contracts повторно генерирует код и требует пустой `git diff`, чтобы обнаруживать несогласованные правки контрактов и DTO.
