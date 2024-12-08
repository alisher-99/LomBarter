package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.com/example/gophers/libs/kafka/producer"
	"gitlab.com/example/gophers/libs/logger"
	"gitlab.com/example/gophers/libs/trace"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
	"github.com/alisher-99/LomBarter/internal/domain/form"
	"github.com/alisher-99/LomBarter/internal/domain/presenter"
	"github.com/alisher-99/LomBarter/internal/domain/repository"
	"github.com/alisher-99/LomBarter/pkg/metrics"
)

// tracerName название трейса.
const tracerName = "service"

// UserService представляет интерфейс для работы с сервисом пользователей.
type UserService interface {
	// GetUsersByBio возвращает список пользователей по bio.
	GetUsersByBio(ctx context.Context, filter form.UsersGetByBio) (entity.Users, error)
	// GetUserByID возвращает пользователя по идентификатору.
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	// CreateUser сохраняет пользователя.
	CreateUser(ctx context.Context, createForm form.UserCreate, currentTime time.Time) (presenter.CreatedUser, error)
	// UpdateUser обновляет пользователя.
	UpdateUser(ctx context.Context, user form.UserUpdate, currentTime time.Time) error
}

// userService представляет сервис для работы с пользователей.
type userService struct {
	userRepo  repository.UserRepository // Репозиторий для работы с пользователями
	cacheData repository.CacheStore     // Кэш для хранения данных о пользователях
	tracer    trace.TracerProvider      // Отслеживает запросы между слоями и микросервисами
	logger    logger.Logger             // Логирование запросов и ошибок сервиса
	producer  producer.MessageProducer  // Продюсер в топик Кафки
	metrics   metrics.UserMetrics       // Метрики пользователей
	json      jsoniter.API              // JSON-парсер
}

// NewUserService создает новый экземпляр сервиса для работы с пользователями.
func NewUserService(
	repo repository.UserRepository,
	cacheData repository.CacheStore,
	l logger.Logger,
	tracer trace.TracerProvider,
	kafkaProducer producer.MessageProducer,
	userMetrics metrics.UserMetrics,
) UserService {
	return &userService{
		userRepo:  repo,
		cacheData: cacheData,
		logger:    l.WithFields(logger.Fields{"layer": "updateForm-service"}),
		tracer:    tracer,
		producer:  kafkaProducer,
		metrics:   userMetrics,
		json:      jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

// GetUsersByBio возвращает список пользователей, использует кэш.
func (u *userService) GetUsersByBio(ctx context.Context, filter form.UsersGetByBio) (entity.Users, error) {
	ctx, span := u.tracer.Tracer(tracerName).Start(ctx, "UserService.GetUsersByBio")
	defer span.End()

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("валидация фильтра: %w", err)
	}

	users, err := u.userRepo.GetUsersByBio(ctx, filter)
	if err != nil {
		u.metrics.IncFailedReceivingUsers()

		return nil, fmt.Errorf("получение пользователей: %w", err)
	}

	u.metrics.IncSuccessfulReceivingUsers()

	return users, nil
}

// GetUserByID возвращает пользователя по идентификатору.
func (u *userService) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	ctx, span := u.tracer.Tracer(tracerName).Start(ctx, "UserService.GetUserByID")
	defer span.End()

	user, cErr := u.cacheData.UserCache().GetUserByID(ctx, id)
	if cErr == nil {
		return user, nil
	}

	if !errors.Is(cErr, entity.ErrUserNotFound) {
		u.logger.WithFields(logger.Fields{"id": id}).Errorf("получение пользователя из кэша: %v", cErr)
	}

	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return &entity.User{}, fmt.Errorf("получение пользователя: %w", err)
	}

	if cErr = u.cacheData.UserCache().SetUser(ctx, user); cErr != nil {
		u.logger.WithFields(logger.Fields{"id": id}).Errorf("установка кэша: %v", cErr)
	}

	return user, nil
}

// CreateUser сохраняет пользователя.
func (u *userService) CreateUser(ctx context.Context, createForm form.UserCreate, currentTime time.Time) (presenter.CreatedUser, error) {
	ctx, span := u.tracer.Tracer(tracerName).Start(ctx, "UserService.CreateUser")
	defer span.End()

	if err := createForm.Validate(); err != nil {
		return presenter.CreatedUser{}, fmt.Errorf("валидация формы: %w", err)
	}

	// Создаем сущность пользователя.
	user := entity.NewUser(currentTime)

	if err := createForm.Fill(user); err != nil {
		return presenter.CreatedUser{}, fmt.Errorf("заполнение формы: %w", err)
	}

	id, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return presenter.CreatedUser{}, fmt.Errorf("создание пользователя: %w", err)
	}

	return presenter.NewCreatedUser(id), nil
}

// UpdateUser обновляет пользователя.
func (u *userService) UpdateUser(ctx context.Context, updateForm form.UserUpdate, currentTime time.Time) error {
	ctx, span := u.tracer.Tracer(tracerName).Start(ctx, "UserService.UpdateUser")
	defer span.End()

	if err := updateForm.Validate(); err != nil {
		return fmt.Errorf("валидация формы: %w", err)
	}

	// Получаем пользователя.
	user, err := u.userRepo.GetUserByID(ctx, updateForm.ID)
	if err != nil {
		return fmt.Errorf("получение пользователя: %w", err)
	}

	// Заполняем сущность пользователя обновленными данными.
	if err = updateForm.Fill(user, currentTime); err != nil {
		return fmt.Errorf("заполнение формы: %w", err)
	}

	// Обновляем кэш.
	if err = u.cacheData.UserCache().SetUser(ctx, user); err != nil {
		return fmt.Errorf("удаление кэша: %w", err)
	}

	// Обновляем пользователя.
	if err = u.userRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("обновление пользователя: %w", err)
	}

	u.logger.Debugf("отправка уведомления: %v", user)

	// Пример отправки сообщения в топик Кафки
	err = u.producer.Write(ctx, producer.Message{Value: []byte(fmt.Sprintf("User %q updated", user.ID))})
	if err != nil {
		return fmt.Errorf("запись в kafka: %w", err)
	}

	return nil
}
