package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/example/gophers/libs/logger"
)

// WithLogger добавляет логгер в Prometheus сервер.
func WithLogger(log logger.Logger) Option {
	return func(srv *Server) {
		srv.logger = log
	}
}

// WithRegistry добавляет реестр метрик в Prometheus сервер.
func WithRegistry(registry *prometheus.Registry) Option {
	return func(s *Server) {
		s.registry = registry
	}
}
