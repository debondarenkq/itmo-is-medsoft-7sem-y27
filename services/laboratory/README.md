# Laboratory service

Собственная БД `laboratory`: справочник LOINC, назначения и результаты FHIR R5. Миграции Goose выполняются перед HTTP-сервером. Внешние маршруты и права описаны в `contracts/laboratory.yaml`; типы и обработчики генерируются `oapi-codegen`.

Обязательные переменные: `DATABASE_URL`, `AUTH_URL`, `CLINICAL_URL`, `STAFF_URL`, `SERVICE_KEY_FILE`. `HTTP_ADDR` по умолчанию `:8080`. Доступ к исходным Patient/Practitioner выполняется по защищённым внутренним API clinical/staff. Снаружи эти внутренние маршруты не публикуются.

`go test ./...`, `go vet ./...`. Интеграционные тесты включаются через `TEST_DATABASE_URL`; из корня — `make integration`. Сценарий работы и развёртывание: [лабораторная 3](../../docs/lab3.md).

`DATABASE_PASSWORD_FILE` загружает пароль БД из локального игнорируемого файла. В Compose новые пароли не хранятся в исходниках; `make certs` создаёт их один раз.
