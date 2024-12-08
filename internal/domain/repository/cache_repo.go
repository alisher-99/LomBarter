package repository

import (
	"context"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
)

// CacheStore представляет интерфейс для работы с кэшем.
type CacheStore interface {
	// UserCache возвращает репозиторий пользователей.
	UserCache() UserCache
}

// UserCache представляет интерфейс для работы с кэшем пользователей.
type UserCache interface {
	// GetUserByID возвращает пользователя по идентификатору.
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	// SetUser сохраняет пользователя.
	SetUser(ctx context.Context, user *entity.User) error
}
