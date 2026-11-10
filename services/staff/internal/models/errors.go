package models

import "errors"

var (
	ErrForbidden    = errors.New("Недостаточно прав")
	ErrInvalidInput = errors.New("Некорректные данные")
	ErrNotFound     = errors.New("Объект не найден")
	ErrConflict     = errors.New("Конфликт данных")
)
