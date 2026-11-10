# HTTP-контракты HIS

Редактируемые OpenAPI 3.0.3 схемы всех микросервисов находятся здесь:

| Контракт | Владелец | Потребители |
|---|---|---|
| [staff.yaml](staff.yaml) | staff-service | clinical-service, SwiftUI |
| [catalog.yaml](catalog.yaml) | catalog-service | clinical-service, SwiftUI |
| [clinical.yaml](clinical.yaml) | clinical-service | SwiftUI |
| [security.yaml](security.yaml) | identity-service | SwiftUI, внешний демо-скрипт |
| [audit.yaml](audit.yaml) | identity-service | staff-service, catalog-service, clinical-service |
| [internal-identity.yaml](internal-identity.yaml) | identity-service, staff-service | внутренние HTTP-клиенты сервисов |

После изменения схемы выполните из корня репозитория:

```sh
make generate-contracts
make integration check macos-build
```

Генерация обновляет локальные копии схем в сервисах, Go DTO и HTTP-клиенты, JSON-документацию, объединённую спецификацию gateway и Swift-клиент. Локальные копии в `services/*/api/openapi.yaml` и `services/clinical/internal/pkg/adapters/*/openapi.yaml` не редактируются вручную: они нужны для самостоятельной сборки и генерации внутри отдельного сервиса.

Сервисы не импортируют Go-код друг друга. Общий контракт — HTTP/JSON; взаимодействие остаётся через сгенерированные клиенты. Полный порядок работы описан в [docs/contracts.md](../docs/contracts.md).

Безопасность лабораторной 2 описана в `security.yaml`: вход/настройка администратора, обновление токенов, аккаунты и аудит. Go binding генерируется в `services/identity/api/security/`, Swift-клиент включается в общий HISAPI.

Внутренние операции авторизации и прямого чтения/изменения роли сотрудников описаны в `internal-identity.yaml`; они не публикуются gateway и требуют `X-Service-Key`.

Общий межсервисный аудит описан в `audit.yaml`: HTTP-запрос и перечисление событий используются всеми сервисами. Публичный контракт безопасности ссылается на то же перечисление. Каждый потребитель генерирует нужные серверные типы или клиент через `oapi-codegen` со своей локальной конфигурацией. Метаданные `x-audit` HTTP-операций читаются стандартным OpenAPI-router из сгенерированной спецификации; отдельного генератора маршрутов нет.
