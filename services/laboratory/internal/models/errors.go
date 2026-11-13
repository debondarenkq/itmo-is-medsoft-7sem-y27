package models

import "errors"

var (
	ErrInvalidInput = errors.New("Некорректные данные")
	ErrNotFound     = errors.New("Объект не найден")
	ErrConflict     = errors.New("Конфликт данных")
	ErrUnavailable  = errors.New("Справочный сервис недоступен или вернул некорректный ответ")
)
