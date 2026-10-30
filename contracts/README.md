# HTTP-контракты HIS

Редактируемые OpenAPI 3.0.3 схемы всех микросервисов находятся здесь:

| Контракт | Владелец | Потребители |
|---|---|---|
| [staff.yaml](staff.yaml) | staff-service | clinical-service, SwiftUI |
| [catalog.yaml](catalog.yaml) | catalog-service | clinical-service, SwiftUI |
| [clinical.yaml](clinical.yaml) | clinical-service | SwiftUI |

После изменения схемы выполните из корня репозитория:

```sh
make generate-contracts
make integration check macos-build
```

Генерация обновляет локальные копии схем в сервисах, Go DTO и HTTP-клиенты, JSON-документацию, объединённую спецификацию gateway и Swift-клиент. Локальные копии в `services/*/api/openapi.yaml` и `services/clinical/internal/pkg/adapters/*/openapi.yaml` не редактируются вручную: они нужны для самостоятельной сборки и генерации внутри отдельного сервиса.

Сервисы не импортируют Go-код друг друга. Общий контракт — HTTP/JSON; взаимодействие остаётся через сгенерированные клиенты. Полный порядок работы описан в [docs/contracts.md](../docs/contracts.md).

Безопасность лабораторной 2 описана в `security.yaml`: вход/настройка администратора, обновление токенов, аккаунты и аудит. Go binding генерируется в `services/identity/api/security/`, Swift-клиент включается в общий HISAPI.

Внутренние операции синхронизации сотрудников описаны в `internal-identity.yaml`; они не публикуются gateway и требуют `X-Service-Key`.
