package presenter

import "github.com/alisher-99/LomBarter/internal/domain/entity"

// CreatedOrder информация о созданном заказе.
type CreatedOrder struct {
	ID     string `json:"id" example:"655d8a3577a0a79c69a7cdfc"`     // Идентификатор заказа
	UserID string `json:"userID" example:"655d8a4d3afea534e56b570e"` // Идентификатор пользователя
	Cost   int    `json:"cost" example:"39900"`                      // Стоимость заказа
}

// NewCreatedOrder создает информацию о созданном заказе.
func NewCreatedOrder(order *entity.Order) CreatedOrder {
	if order == nil {
		return CreatedOrder{}
	}

	return CreatedOrder{
		ID:     order.ID,
		UserID: order.UserID,
		Cost:   order.Cost,
	}
}
