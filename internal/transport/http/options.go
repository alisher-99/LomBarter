package http

import (
	"gitlab.com/example/gophers/libs/logger"
	"gitlab.com/example/gophers/libs/trace"

	"github.com/alisher-99/LomBarter/internal/service"
)

// Option определяет функцию для настройки HTTP сервера.
type Option func(*Server)

// WithUserService добавляет сервис пользователей в HTTP сервер.
func WithUserService(userService service.UserService) Option {
	return func(srv *Server) {
		srv.userService = userService
	}
}

// WithOrdersService добавляет сервис заказов в HTTP сервер.
func WithOrdersService(ordersService service.OrdersService) Option {
	return func(srv *Server) {
		srv.ordersService = ordersService
	}
}

// WithLogger добавляет логгер в HTTP сервер.
func WithLogger(log logger.Logger) Option {
	return func(srv *Server) {
		srv.logger = log
	}
}

// WithTracer добавляет трейсер в HTTP сервер.
func WithTracer(tracer trace.TracerProvider) Option {
	return func(srv *Server) {
		srv.tracer = tracer
	}
}
