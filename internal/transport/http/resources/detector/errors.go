package detector

import (
	"errors"

	"github.com/go-chi/render"
	"gitlab.com/example/gophers/libs/errors/httperrors"
	"gitlab.com/example/gophers/libs/validate"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
)

// Error обрабатывает ошибки, возникающие при работе с ресурсами.
func Error(err error) render.Renderer {
	if err == nil {
		return nil
	}

	validationErr := validate.ValidationError{}
	if errors.As(err, &validationErr) {
		return httperrors.BadRequest(err, validationErr.Code)
	}

	renderer := pageDetect(err)
	if renderer != nil {
		return renderer
	}

	renderer = userDetect(err)
	if renderer != nil {
		return renderer
	}

	renderer = ordersDetect(err)
	if renderer != nil {
		return renderer
	}

	return httperrors.Internal(err, entity.InternalCode)
}

// pageDetect обрабатывает ошибки, возникающие при работе с пагинацией.
func pageDetect(err error) render.Renderer {
	switch {
	case errors.Is(err, entity.ErrPageInvalidLimit):
		return httperrors.BadRequest(err, entity.PageInvalidLimitCode)
	case errors.Is(err, entity.ErrPageInvalidState):
		return httperrors.BadRequest(err, entity.PageInvalidStateCode)
	default:
		return nil
	}
}

// userDetect обрабатывает ошибки, возникающие при работе с пользователями.
func userDetect(err error) render.Renderer {
	switch {
	case errors.Is(err, entity.ErrUserNotFound), errors.Is(err, entity.ErrInvalidObjectID):
		return httperrors.ResourceNotFound(err, entity.UserNotFoundCode)
	case errors.Is(err, entity.ErrUserIDEmpty):
		return httperrors.BadRequest(err, entity.UserIDEmptyCode)
	case errors.Is(err, entity.ErrUserDecode):
		return httperrors.BadRequest(err, entity.UserDecodeCode)
	default:
		return nil
	}
}

// ordersDetect обрабатывает ошибки, возникающие при работе с заказами.
func ordersDetect(err error) render.Renderer {
	switch {
	case errors.Is(err, entity.ErrOrderDecode):
		return httperrors.BadRequest(err, entity.OrderDecodeCode)
	case errors.Is(err, entity.ErrOrderNotFound):
		return httperrors.BadRequest(err, entity.OrderNotFoundCode)
	default:
		return nil
	}
}
