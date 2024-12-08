package entity

// List список сущностей.
type List struct {
	Items interface{} `json:"items"`           // Список сущностей
	Count int64       `json:"count,omitempty"` // Количество сущностей
	State string      `json:"state,omitempty"` // Состояние пагинации
}

// Response ответ сервера.
type Response struct {
	Detail string `json:"detail"` // Детальное описание ответа
}
