package form

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

const (
	// defaultLimit лимит по умолчанию.
	defaultLimit = 10
	// defaultOrderBy сортировка по умолчанию.
	defaultOrderBy = ASC
)

// Pagination представляет форму пагинации.
type Pagination struct {
	Page    uint64 `json:"page" validate:"omitempty,min=1"`              // Номер страницы. Используется для пагинации в mongo
	Limit   uint64 `json:"limit" validate:"omitempty,min=1,max=100"`     // Количество элементов на странице
	OrderBy string `json:"order_by" validate:"omitempty,oneof=asc desc"` // Сортировка. asc - по возрастанию, desc - по убыванию

	PageState      string `json:"page_state,omitempty"` // Состояние страницы, строка в base64. Используется для пагинации в кассандре
	PageStateBytes []byte `json:"-"`                    // Состояние страницы в байтах
}

// NewPagination возвращает новую пагинацию.
func NewPagination(limit uint64, pageState, orderBy string) (Pagination, error) {
	p := Pagination{
		Limit:   defaultLimit,
		OrderBy: defaultOrderBy,
	}

	if limit > 0 {
		p.Limit = limit
	}

	if orderBy != "" {
		p.OrderBy = orderBy
	}

	if pageState != "" {
		err := p.SetPageState(pageState)
		if err != nil {
			return Pagination{}, fmt.Errorf("установка состояния страницы: %w", err)
		}
	}

	return p, nil
}

// ParsePagination парсит пагинацию из url.
func ParsePagination(values url.Values) (Pagination, error) {
	pagination, err := NewPagination(0, "", "")
	if err != nil {
		return Pagination{}, fmt.Errorf("создание новой пагинации: %w", err)
	}

	if str := values.Get("limit"); str != "" {
		limit, pErr := strconv.ParseUint(str, 10, 64)
		if pErr != nil {
			return Pagination{}, fmt.Errorf("%w: %s", entity.ErrPageInvalidLimit, pErr.Error())
		}

		pagination.Limit = limit
	}

	if str := values.Get("order_by"); str != "" {
		pagination.OrderBy = str
	}

	if str := values.Get("page_state"); str != "" {
		err = pagination.SetPageState(str)
		if err != nil {
			return Pagination{}, fmt.Errorf("установка состояния страницы: %w", err)
		}
	}

	if str := values.Get("page"); str != "" {
		page, pErr := strconv.ParseUint(str, 10, 64)
		if pErr != nil {
			return Pagination{}, fmt.Errorf("%w: %s", entity.ErrPageInvalidPage, pErr.Error())
		}

		pagination.Page = page
	}

	return pagination, nil
}

// Offset возвращает смещение.
func (p *Pagination) Offset() uint64 {
	if p.Page == 0 {
		p.Page = 1
	}

	return (p.Page - 1) * p.Limit
}

// SortToInt возвращает сортировку в виде числа для Mongo.
func (p *Pagination) SortToInt() int {
	if p.OrderBy == DESC {
		return -1
	}

	return 1
}

// SortToBool возвращает сортировку в виде булевого значения для Cassandra.
func (p *Pagination) SortToBool() bool {
	return p.OrderBy != DESC
}

// SetPageState устанавливает состояние страницы.
func (p *Pagination) SetPageState(state interface{}) error {
	switch v := state.(type) {
	case []byte:
		p.PageStateBytes = v
		p.PageState = base64.StdEncoding.EncodeToString(v)
	case string:
		bytes, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return fmt.Errorf("%w: %s", entity.ErrPageInvalidState, err.Error())
		}

		p.PageState = v
		p.PageStateBytes = bytes
	}

	return nil
}
