package entity

import "time"

// Order сущность заказа.
type Order struct {
	ID        string    `json:"id" db:"id" bson:"_id"`                       // Идентификатор заказа
	UserID    string    `json:"userID" db:"user_id" bson:"user_id"`          // Идентификатор пользователя
	Cost      int       `json:"cost" db:"cost" bson:"cost"`                  // Стоимость заказа
	CreatedAt time.Time `json:"createdAt" db:"created_at" bson:"created_at"` // Дата создания заказа
}

// NewOrder создает заказ.
func NewOrder(currentTime time.Time) *Order {
	return &Order{
		CreatedAt: currentTime,
	}
}

// Orders список заказов.
type Orders []*Order
