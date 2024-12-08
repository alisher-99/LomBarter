package presenter

// CreatedUser информация о созданном пользователе.
type CreatedUser struct {
	ID string `json:"id" example:"5f8b9b1b3afea534e56b570e"` // Идентификатор пользователя
}

// NewCreatedUser возвращает информацию о созданном пользователе.
func NewCreatedUser(id string) CreatedUser {
	return CreatedUser{ID: id}
}
