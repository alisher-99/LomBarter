package form

import (
	"gitlab.com/example/gophers/libs/validate"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
)

// OrderCreate форма создания заказа.
type OrderCreate struct {
	UserID string `json:"-" validate:"required" example:"655d8a4d3afea534e56b570e"` // Идентификатор пользователя. Передается в заголовке X-User-Id
	Cost   int    `json:"cost" validate:"required,gt=0" example:"39900"`            // Стоимость заказа
}

// Validate валидирует форму создания заказа.
func (f *OrderCreate) Validate() error {
	if f == nil {
		return entity.ErrNilPointer
	}

	return validate.New(shortServiceName).Validate(f)
}

// Fill заполняет сущность заказа.
func (f *OrderCreate) Fill(order *entity.Order) error {
	if f == nil || order == nil {
		return entity.ErrNilPointer
	}

	order.UserID = f.UserID
	order.Cost = f.Cost

	return nil
}

// OrdersGetForClient форма получения списка заказов для клиента.
type OrdersGetForClient struct {
	UserID string `json:"userID" validate:"required" example:"655d8a4d3afea534e56b570e"` // Идентификатор пользователя
}

// Validate валидирует форму получения списка заказов для клиента.
func (f OrdersGetForClient) Validate() error {
	return validate.New(shortServiceName).Validate(f)
}

// OrderGetForClient форма получения заказа для клиента.
type OrderGetForClient struct {
	OrderID string `json:"orderID" validate:"required,mongodb" example:"5f8b9b1b3afea534e56b570e"` // Идентификатор заказа
	UserID  string `json:"-" validate:"required" example:"655d8a4d3afea534e56b570e"`               // Идентификатор пользователя. Передается в заголовке X-User-Id
}

// Validate валидирует форму получения заказа для клиента.
func (f OrderGetForClient) Validate() error {
	return validate.New(shortServiceName).Validate(f)
}
