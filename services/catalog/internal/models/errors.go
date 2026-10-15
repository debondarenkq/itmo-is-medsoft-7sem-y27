package models

import "errors"

var (
	ErrDuplicateCode = errors.New("Повторяющийся код")
	ErrDuplicateName = errors.New("Повторяющееся наименование")
	ErrInvalidInput  = errors.New("Некорректные данные")
	ErrNotFound      = errors.New("Объект не найден")
	ErrConflict      = errors.New("Конфликт данных")
)
