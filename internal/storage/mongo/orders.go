package mongo

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/example/gophers/libs/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
	"github.com/alisher-99/LomBarter/internal/domain/form"
	"github.com/alisher-99/LomBarter/internal/domain/repository"
)

// ordersRepository репозиторий заказов.
type ordersRepository struct {
	collection *mongo.Collection    // Коллекция пользователей
	tracer     trace.TracerProvider // Отслеживает запросы между слоями и микросервисами
}

// NewOrdersRepository возвращает новый экземпляр репозитория заказов.
func NewOrdersRepository(collection *mongo.Collection, tracer trace.TracerProvider) repository.OrdersRepository {
	return ordersRepository{collection: collection, tracer: tracer}
}

// CreateOrder создает новый заказ.
func (o ordersRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	ctx, span := o.tracer.Tracer(tracerName).Start(ctx, "OrdersRepository.CreateOrder")
	defer span.End()

	document := bson.D{
		{Key: "user_id", Value: order.UserID},
		{Key: "cost", Value: order.Cost},
		{Key: "created_at", Value: order.CreatedAt},
	}

	res, err := o.collection.InsertOne(ctx, document)
	if err != nil {
		return fmt.Errorf("добавление документа в коллекцию: %w", err)
	}

	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("%w: %v", entity.ErrInvalidObjectID, res.InsertedID)
	}

	order.ID = objID.Hex()

	return nil
}

// GetOrdersForClient возвращает список заказов для клиента.
func (o ordersRepository) GetOrdersForClient(ctx context.Context, filter form.OrdersGetForClient) (entity.Orders, error) {
	ctx, span := o.tracer.Tracer(tracerName).Start(ctx, "OrdersRepository.GetOrdersForClient")
	defer span.End()

	match := bson.D{
		{Key: "user_id", Value: filter.UserID},
	}

	cur, err := o.collection.Find(ctx, match)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("получение списка заказов: %w", entity.ErrOrderNotFound)
		}

		return nil, fmt.Errorf("получение списка заказов: %w", err)
	}

	orders := make(entity.Orders, 0, cur.RemainingBatchLength())
	if err = cur.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("декодирование списка заказов: %w", err)
	}

	return orders, nil
}

// GetOrderForClient возвращает заказ для клиента.
func (o ordersRepository) GetOrderForClient(ctx context.Context, filter form.OrderGetForClient) (*entity.Order, error) {
	ctx, span := o.tracer.Tracer(tracerName).Start(ctx, "OrdersRepository.GetOrderForClient")
	defer span.End()

	idObj, err := primitive.ObjectIDFromHex(filter.OrderID)
	if err != nil {
		return nil, fmt.Errorf("получение идентификатора заказа: %w", entity.ErrInvalidObjectID)
	}

	match := bson.D{
		{Key: "_id", Value: idObj},
		{Key: "user_id", Value: filter.UserID},
	}

	var order entity.Order
	if err = o.collection.FindOne(ctx, match).Decode(&order); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("получение списка заказов: %w", entity.ErrOrderNotFound)
		}

		return nil, fmt.Errorf("получение заказа: %w", err)
	}

	return &order, nil
}
