// Package swagger содержит иноформативную мета информацию, которая нужна для корректности генерируемых swagger документации для tmp.
package swagger

// HTTPResponse400 структура, которая отображается как тело ответа при 400 коде возврата от HTTP.
type HTTPResponse400 struct {
	Code string `json:"code" example:"TMP_INVALID_USER"` // Код ошибки.
}

// HTTPResponse500 структура, которая отображается как тело ответа при 500 коде возврата от HTTP.
type HTTPResponse500 struct {
	Code string `json:"code" example:"TMP_INTERNAL"` // Код ошибки.
}
