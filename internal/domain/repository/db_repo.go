package repository

import (
	"context"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
	"github.com/alisher-99/LomBarter/internal/domain/form"
)

type DataStore interface {
	// TxStarter интерфейс для работы с транзакциями
	// TxStarter

	// Base базовый интерфейс для работы с DataStore
	Base
	// UserRepository возвращает репозиторий пользователей.
	UserRepository() UserRepository
	// OrdersRepository возвращает репозиторий заказов.
	OrdersRepository() OrdersRepository
}

// Base представляет базовый интерфейс для работы с DataStore.
type Base interface {
	// Name возвращает название DataStore
	Name() string

	// Close закрывает соединение с DataStore
	Close(ctx context.Context) error

	// Connect устанавливает соединение с DataStore
	Connect() error
}

// UserRepository представляет интерфейс для работы с репозиторием пользователей.
type UserRepository interface {
	// GetUsersByBio возвращает список пользователей по bio.
	GetUsersByBio(ctx context.Context, filter form.UsersGetByBio) (entity.Users, error)
	// GetUserByID возвращает пользователя по идентификатору.
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	// CreateUser сохраняет пользователя.
	CreateUser(ctx context.Context, user *entity.User) (string, error)
	// UpdateUser обновляет пользователя.
	UpdateUser(ctx context.Context, user *entity.User) error
}

// OrdersRepository представляет интерфейс для работы с репозиторием заказов.
type OrdersRepository interface {
	// CreateOrder создает новый заказ.
	CreateOrder(ctx context.Context, order *entity.Order) error
	// GetOrdersForClient возвращает список заказов для клиента.
	GetOrdersForClient(ctx context.Context, filter form.OrdersGetForClient) (entity.Orders, error)
	// GetOrderForClient возвращает заказ для клиента.
	GetOrderForClient(ctx context.Context, filter form.OrderGetForClient) (*entity.Order, error)
}

// TxCallback представляет функцию обратного вызова для обработки результатов транзакции.
type TxCallback func(context.Context, error) error

// TxStarter определяет интерфейс для запуска транзакций.
type TxStarter interface {
	// StartSession создает сессию для транзакции
	StartSession(ctx context.Context) (context.Context, TxCallback, error)
}
