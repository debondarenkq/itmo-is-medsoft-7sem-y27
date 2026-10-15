package models

import "errors"

var (
	ErrInvalidInput = errors.New("Некорректные данные")
	ErrNotFound     = errors.New("Объект не найден")
	ErrConflict     = errors.New("Конфликт данных")
)
