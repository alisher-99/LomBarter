package service

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/example/gophers/libs/logger"
	"gitlab.com/example/gophers/libs/trace"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
	"github.com/alisher-99/LomBarter/internal/domain/form"
	"github.com/alisher-99/LomBarter/internal/domain/presenter"
	"github.com/alisher-99/LomBarter/internal/domain/repository"
)

// OrdersService представляет собой интерфейс сервиса для работы с заказами.
type OrdersService interface {
	// CreateOrder создает новый заказ.
	CreateOrder(ctx context.Context, createForm form.OrderCreate, currentTime time.Time) (presenter.CreatedOrder, error)
	// GetOrdersForClient возвращает список заказов для клиента.
	GetOrdersForClient(ctx context.Context, form form.OrdersGetForClient) (entity.Orders, error)
	// GetOrderForClient возвращает заказ для клиента.
	GetOrderForClient(ctx context.Context, form form.OrderGetForClient) (*entity.Order, error)
}

// orderService представляет сервис для работы с заказами.
type ordersService struct {
	ordersRepository repository.OrdersRepository // Репозиторий для работы с заказами
	tracer           trace.TracerProvider        // Отслеживает запросы между слоями и микросервисами.
	logger           logger.Logger               // Логирование запросов и ошибок сервиса.
}

// NewOrdersService создает новый экзмепляр сервиса для работы с заказами.
func NewOrdersService(ordersRepository repository.OrdersRepository, l logger.Logger, tracer trace.TracerProvider) OrdersService {
	return &ordersService{
		ordersRepository: ordersRepository,
		tracer:           tracer,
		logger:           l.WithFields(logger.Fields{"layer": "orders-service"}),
	}
}

// CreateOrder создает новый заказ.
func (s ordersService) CreateOrder(ctx context.Context, createForm form.OrderCreate, currentTime time.Time) (presenter.CreatedOrder, error) {
	ctx, span := s.tracer.Tracer(tracerName).Start(ctx, "OrdersService.CreateOrder")
	defer span.End()

	if err := createForm.Validate(); err != nil {
		return presenter.CreatedOrder{}, fmt.Errorf("валидация формы: %w", err)
	}

	// Создаем сущность заказа.
	order := entity.NewOrder(currentTime)

	if err := createForm.Fill(order); err != nil {
		return presenter.CreatedOrder{}, fmt.Errorf("заполнение сущности заказа: %w", err)
	}

	// Сохраняем заказ в репозитории.
	if err := s.ordersRepository.CreateOrder(ctx, order); err != nil {
		return presenter.CreatedOrder{}, fmt.Errorf("создание заказа: %w", err)
	}

	// Возвращаем информацию о созданном заказе.
	return presenter.NewCreatedOrder(order), nil
}

// GetOrdersForClient возвращает список заказов для клиента.
func (s ordersService) GetOrdersForClient(ctx context.Context, filter form.OrdersGetForClient) (entity.Orders, error) {
	ctx, span := s.tracer.Tracer(tracerName).Start(ctx, "OrdersService.GetOrdersForClient")
	defer span.End()

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("валидация формы: %w", err)
	}

	orders, err := s.ordersRepository.GetOrdersForClient(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("получение списка заказов: %w", err)
	}

	return orders, nil
}

// GetOrderForClient возвращает заказ для клиента.
func (s ordersService) GetOrderForClient(ctx context.Context, filter form.OrderGetForClient) (*entity.Order, error) {
	ctx, span := s.tracer.Tracer(tracerName).Start(ctx, "OrdersService.GetOrderForClient")
	defer span.End()

	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("валидация формы: %w", err)
	}

	order, err := s.ordersRepository.GetOrderForClient(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("получение заказа: %w", err)
	}

	return order, nil
}
