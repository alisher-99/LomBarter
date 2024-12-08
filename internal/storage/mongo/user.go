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

// userRepository репозиторий пользователей.
type userRepository struct {
	collection *mongo.Collection    // Коллекция пользователей
	tracer     trace.TracerProvider // Отслеживает запросы между слоями и микросервисами
}

// NewUserRepository возвращает новый экземпляр репозитория пользователей.
func NewUserRepository(collection *mongo.Collection, tracer trace.TracerProvider) repository.UserRepository {
	return &userRepository{collection: collection, tracer: tracer}
}

// GetUsersByBio возвращает список пользователей по bio.
func (r userRepository) GetUsersByBio(ctx context.Context, filter form.UsersGetByBio) (entity.Users, error) {
	ctx, span := r.tracer.Tracer(tracerName).Start(ctx, "UserRepository.GetUsersByBio")
	defer span.End()

	match := bson.D{{Key: "bio", Value: filter.Bio}}

	cursor, err := r.collection.Find(ctx, match)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make(entity.Users, 0, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("получение списка пользователей: %w", err)
	}

	return users, nil
}

// GetUserByID возвращает пользователя по идентификатору.
func (r userRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	ctx, span := r.tracer.Tracer(tracerName).Start(ctx, "UserRepository.GetUserByID")
	defer span.End()

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrInvalidObjectID, err.Error())
	}

	match := bson.D{{Key: "_id", Value: idObj}}

	var user *entity.User
	if err = r.collection.FindOne(ctx, match).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrUserNotFound
		}

		return nil, fmt.Errorf("получение пользователя: %w", err)
	}

	return user, nil
}

// CreateUser сохраняет пользователя.
func (r userRepository) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	ctx, span := r.tracer.Tracer(tracerName).Start(ctx, "UserRepository.CreateUser")
	defer span.End()

	document := bson.D{
		{Key: "name", Value: user.Name},
		{Key: "bio", Value: user.Bio},
		{Key: "created_at", Value: user.CreatedAt},
		{Key: "updated_at", Value: user.UpdatedAt},
	}

	res, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("сохранение пользователя: %w", err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// UpdateUser обновляет пользователя.
func (r userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	ctx, span := r.tracer.Tracer(tracerName).Start(ctx, "UserRepository.UpdateUser")
	defer span.End()

	idObj, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("%w: %s", entity.ErrInvalidObjectID, err.Error())
	}

	match := bson.D{{Key: "_id", Value: idObj}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: user.Name},
		{Key: "bio", Value: user.Bio},
		{Key: "updated_at", Value: user.UpdatedAt},
	}}}

	_, err = r.collection.UpdateOne(ctx, match, update)

	return err
}
