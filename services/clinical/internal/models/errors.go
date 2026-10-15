package models

import "errors"

var (
	ErrReferenced      = errors.New("Запись связана с другими данными")
	ErrAlreadyExists   = errors.New("Запись уже существует")
	ErrInvalidInput    = errors.New("Некорректные данные")
	ErrNotFound        = errors.New("Объект не найден")
	ErrConflict        = errors.New("Конфликт данных")
	ErrVersionConflict = errors.New("ЭМК изменена другим пользователем; обновите карту")
	ErrUnavailable     = errors.New("Справочный сервис недоступен или вернул некорректный ответ")
)
