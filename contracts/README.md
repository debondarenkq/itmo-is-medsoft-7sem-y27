# HTTP-контракты HIS

Редактируемые OpenAPI 3.0.3 схемы всех микросервисов находятся здесь:

| Контракт | Владелец | Потребители |
|---|---|---|
| [staff.yaml](staff.yaml) | staff-service | clinical-service, Flutter |
| [catalog.yaml](catalog.yaml) | catalog-service | clinical-service, Flutter |
| [clinical.yaml](clinical.yaml) | clinical-service | Flutter |

После изменения схемы выполните из корня репозитория:

```sh
make generate-contracts
make integration check desktop-test
```

Генерация обновляет локальные копии схем в сервисах, Go DTO и HTTP-клиенты, JSON-документацию, объединённую спецификацию gateway и Dart-клиент. Локальные копии в `services/*/api/openapi.yaml` и `services/clinical/internal/pkg/adapters/*/openapi.yaml` не редактируются вручную: они нужны для самостоятельной сборки и генерации внутри отдельного сервиса.

Сервисы не импортируют Go-код друг друга. Общий контракт — HTTP/JSON; взаимодействие остаётся через сгенерированные клиенты. Полный порядок работы описан в [docs/contracts.md](../docs/contracts.md).
