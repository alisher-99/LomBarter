package entity

import (
	"fmt"
	"time"

	"gitlab.com/example/gophers/grpcclients/template/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// User сущность пользователя.
type User struct {
	ID        string    `json:"id" db:"id" bson:"_id"`                       // Идентификатор пользователя
	Name      string    `json:"name" db:"name" bson:"name"`                  // Имя пользователя
	Bio       string    `json:"bio" db:"bio" bson:"bio"`                     // Биография пользователя
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at" bson:"updated_at"` // Дата обновления пользователя
	CreatedAt time.Time `json:"createdAt" db:"created_at" bson:"created_at"` // Дата создания пользователя
}

// NewUser возвращает нового пользователя.
func NewUser(currentTime time.Time) *User {
	return &User{
		UpdatedAt: currentTime,
		CreatedAt: currentTime,
	}
}

// Columns возвращает список колонок.
func (u *User) Columns() []string {
	return []string{"id", "name", "bio", "updated_at", "created_at"}
}

// GetUserCacheKey возвращает ключ для кеширования.
func GetUserCacheKey(id string) string {
	return fmt.Sprintf("data:%s", id)
}

// Users список пользователей.
type Users []User

// PROTOBUF

// FromProtoUser преобразует из proto в User.
func FromProtoUser(user *proto.User) *User {
	return &User{
		ID:        user.Id,
		Name:      user.Name,
		Bio:       user.Bio,
		UpdatedAt: user.UpdatedAt.AsTime(),
		CreatedAt: user.CreatedAt.AsTime(),
	}
}

// ToProto преобразует в proto.
func (u *User) ToProto() *proto.User {
	return &proto.User{
		Id:        u.ID,
		Name:      u.Name,
		Bio:       u.Bio,
		UpdatedAt: timestamppb.New(u.UpdatedAt),
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
}

// ToProto преобразует в proto.
func (u Users) ToProto() *proto.Users {
	var protoUsers []*proto.User
	for _, user := range u {
		protoUsers = append(protoUsers, user.ToProto())
	}

	return &proto.Users{Users: protoUsers}
}
