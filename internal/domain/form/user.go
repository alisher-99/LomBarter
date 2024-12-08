package form

import (
	"time"

	"gitlab.com/example/gophers/grpcclients/template/proto"
	"gitlab.com/example/gophers/libs/validate"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
)

const shortServiceName = "TMP" // Короткое название сервиса

// UserCreate форма создания пользователя.
type UserCreate struct {
	Name string `json:"name" validate:"required,min=3,max=255"` // Имя пользователя
	Bio  string `json:"bio" validate:"omitempty,min=3,max=500"` // Биография пользователя
}

// Validate валидирует форму.
func (c *UserCreate) Validate() error {
	if c == nil {
		return entity.ErrNilPointer
	}

	return validate.New(shortServiceName).Validate(c)
}

// Fill заполняет сущность пользователя.
func (c *UserCreate) Fill(user *entity.User) error {
	if c == nil || user == nil {
		return entity.ErrNilPointer
	}

	user.Name = c.Name
	user.Bio = c.Bio

	return nil
}

// GetUserCreateFromProto преобразует из proto в UserCreate.
func GetUserCreateFromProto(protoUser *proto.User) UserCreate {
	if protoUser == nil {
		return UserCreate{}
	}

	return UserCreate{
		Name: protoUser.Name,
		Bio:  protoUser.Bio,
	}
}

// UsersGetByBio форма получения пользователя по bio.
type UsersGetByBio struct {
	Bio string `json:"bio" db:"bio" validate:"required,min=3,max=255"` // Имя пользователя
}

// Validate валидирует форму.
func (f *UsersGetByBio) Validate() error {
	if f == nil {
		return nil
	}

	return validate.New(shortServiceName).Validate(*f)
}

// UserUpdate форма обновления пользователя.
type UserUpdate struct {
	ID   string  `json:"id" validate:"required" example:"655d8a4d3afea534e56b570e"`   // Идентификатор пользователя
	Name *string `json:"name" validate:"omitempty,min=3,max=255" example:"John"`      // Имя пользователя
	Bio  *string `json:"bio" validate:"omitempty,min=3,max=500" example:"Programmer"` // Биография пользователя
}

// Validate валидирует форму.
func (u UserUpdate) Validate() error {
	return validate.New(shortServiceName).Validate(u)
}

// Fill заполняет сущность пользователя.
func (u UserUpdate) Fill(user *entity.User, currentTime time.Time) error {
	if user == nil {
		return entity.ErrNilPointer
	}

	if u.Name != nil {
		user.Name = *u.Name
	}

	if u.Bio != nil {
		user.Bio = *u.Bio
	}

	user.UpdatedAt = currentTime

	return nil
}
