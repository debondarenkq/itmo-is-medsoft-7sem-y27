package models

import (
	"errors"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
	identity "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/identityapi"
)

var ErrConflict = errors.New("Анализ уже зафиксирован; доступен только просмотр")
var ErrNotFound = errors.New("Запись не найдена")
var ErrInvalid = errors.New("Некорректные данные")
var ErrUnconfigured = errors.New("Введите API-ключ сервисного аккаунта LIS")

type Connection struct {
	Session identity.Session
}
type History = contract.LISHistory
