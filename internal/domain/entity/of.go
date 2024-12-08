package entity

// Pointer возвращает указатель на значение.
func Pointer[T any](v T) *T {
	return &v
}
