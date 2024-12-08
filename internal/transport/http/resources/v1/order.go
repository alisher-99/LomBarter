package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/example/gophers/libs/errors/httperrors"
	"gitlab.com/example/gophers/libs/logger"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
	"github.com/alisher-99/LomBarter/internal/domain/form"
	"github.com/alisher-99/LomBarter/internal/service"
	"github.com/alisher-99/LomBarter/internal/transport/http/resources/detector"
)

const HeaderXUserID = "X-User-Id" // Идентификатор пользователя

// OrdersResource представляет собой обработчик для заказов.
type OrdersResource struct {
	ordersService service.OrdersService // Сервис для работы с пользователями
	logger        logger.Logger         // Логирование запросов и ошибок обработчиков
	json          jsoniter.API          // JSON-парсер
}

// NewOrdersHandler создает новый экземпляр OrdersResource.
func NewOrdersHandler(orderService service.OrdersService, log logger.Logger) *OrdersResource {
	return &OrdersResource{
		ordersService: orderService,
		logger:        log,
		json:          jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

// Routes возвращает роутер для обработчика пользователей.
func (vr OrdersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", vr.createOrder)
	r.Get("/", vr.getOrderList)
	r.Get("/{orderID}", vr.getOrderInfo)

	return r
}

// createOrder создает новый заказ.
// @Summary Создание заказа
// @Description Создание заказа
// @Tags orders
// @Accept json
// @Produce json
// @Param X-User-Id header string true "Идентификатор пользователя"
// @Param order body form.OrderCreate true "Заказ"
// @Success 200 {object} presenter.CreatedOrder
// @Failure 400 {object} swagger.HTTPResponse400 "Код ошибки"
// @Failure 500 {object} swagger.HTTPResponse500 "Внутренняя ошибка сервера"
// @Router /v1/orders [post]
func (vr OrdersResource) createOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var order form.OrderCreate
	if err := vr.json.NewDecoder(r.Body).Decode(&order); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err, entity.OrderDecodeCode))

		return
	}

	order.UserID = r.Header.Get(HeaderXUserID)

	createdOrder, err := vr.ordersService.CreateOrder(ctx, order, time.Now().UTC())
	if err != nil {
		vr.logger.Error("ошибка создания заказа", err)
		_ = render.Render(w, r, detector.Error(err))

		return
	}

	render.JSON(w, r, createdOrder)
}

// getOrderList возвращает список заказов.
// @Summary Список заказов
// @Description Список заказов
// @Tags orders
// @Accept json
// @Produce json
// @Param X-User-Id header string true "Идентификатор пользователя"
// @Success 200 {object} entity.List{items=entity.Orders}
// @Failure 400 {object} swagger.HTTPResponse400 "Код ошибки"
// @Failure 500 {object} swagger.HTTPResponse500 "Внутренняя ошибка сервера"
// @Router /v1/orders [get]
func (vr OrdersResource) getOrderList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := form.OrdersGetForClient{
		UserID: r.Header.Get(HeaderXUserID),
	}

	if err := filter.Validate(); err != nil {
		_ = render.Render(w, r, detector.Error(err))

		return
	}

	orders, err := vr.ordersService.GetOrdersForClient(ctx, filter)
	if err != nil {
		vr.logger.Error("ошибка получения списка заказов", err)
		_ = render.Render(w, r, detector.Error(err))

		return
	}

	render.JSON(w, r, orders)
}

// getOrderInfo возвращает информацию о заказе.
// @Summary Информация о заказе
// @Description Информация о заказе
// @Tags orders
// @Accept json
// @Produce json
// @Param X-User-Id header string true "Идентификатор пользователя"
// @Param id path string true "Идентификатор заказа"
// @Success 200 {object} entity.Order
// @Failure 400 {object} swagger.HTTPResponse400 "Код ошибки"
// @Failure 500 {object} swagger.HTTPResponse500 "Внутренняя ошибка сервера"
// @Router /v1/orders/{orderID} [get]
func (vr OrdersResource) getOrderInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := form.OrderGetForClient{
		OrderID: chi.URLParam(r, "orderID"),
		UserID:  r.Header.Get(HeaderXUserID),
	}

	if err := filter.Validate(); err != nil {
		_ = render.Render(w, r, detector.Error(err))

		return
	}

	order, err := vr.ordersService.GetOrderForClient(ctx, filter)
	if err != nil {
		vr.logger.Error("ошибка получения информации о заказе", err)
		_ = render.Render(w, r, detector.Error(err))

		return
	}

	render.JSON(w, r, order)
}
