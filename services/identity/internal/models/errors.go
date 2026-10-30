package models

import "errors"

var (
	ErrUnavailable  = errors.New("Сервис медработников недоступен")
	ErrInvalidInput = errors.New("Некорректные данные")
	ErrNotFound     = errors.New("Объект не найден")
	ErrConflict     = errors.New("Конфликт данных")
)
