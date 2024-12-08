package entity

import (
	"errors"
)

// Ошибки базы данных.

var (
	ErrInvalidDatastoreHosts = errors.New("неверный хост")
	ErrInvalidDatabaseName   = errors.New("неверное название базы данных")
	ErrInvalidDatabaseURL    = errors.New("неверный URL базы данных")
	ErrInvalidObjectID       = errors.New("неверный идентификатор объекта")
	ErrGenerateSQL           = errors.New("ошибка генерации sql")
)

// Сервисные ошибки.

var (
	ErrTopicNotFound = errors.New("топик не найден")
	ErrTopicsLength  = errors.New("неверное количество топиков")

	ErrNilPointer   = errors.New("значение не может быть nil")
	ErrUserNotFound = errors.New("пользователь не найден")
	ErrUserIDEmpty  = errors.New("идентификатор пуст")
	ErrUserDecode   = errors.New("ошибка декодирования пользователя")

	ErrOrderDecode   = errors.New("ошибка декодирования заказа")
	ErrOrderNotFound = errors.New("заказ не найден")

	ErrPageInvalidLimit = errors.New("неверное значение лимита")
	ErrPageInvalidPage  = errors.New("неверное значение страницы")
	ErrPageInvalidState = errors.New("неверное состояние страницы")
)

// Клиентские ошибки.

// Структура кода ошибки
// TMP - краткое название сервиса
// USER - название сущности
// NOT_FOUND - описание ошибки

const (
	InternalCode = "TMP_INTERNAL" // Внутренняя ошибка сервера

	UserNotFoundCode = "TMP_USER_NOT_FOUND" // Пользователь не найден
	UserIDEmptyCode  = "TMP_USER_ID_EMPTY"  // Идентификатор пуст
	UserDecodeCode   = "TMP_USER_DECODE"    // Ошибка декодирования пользователя

	OrderDecodeCode   = "TMP_ORDER_DECODE"    // Ошибка декодирования заказа
	OrderNotFoundCode = "TMP_ORDER_NOT_FOUND" // Ошибка декодирования заказа

	PageInvalidLimitCode = "TMP_PAGE_INVALID_LIMIT" // Неверное значение лимита
	PageInvalidStateCode = "TMP_PAGE_INVALID_STATE" // Неверное состояние страницы
)
