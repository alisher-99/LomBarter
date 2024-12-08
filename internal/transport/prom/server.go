package prom

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/example/gophers/libs/logger"

	"github.com/alisher-99/LomBarter/internal/config"
)

const (
	readTimeout  = 5 * time.Second  // Время ожидания запроса
	writeTimeout = 30 * time.Second // Время ожидания ответа
)

// Option определяет функцию для настройки PROMETHEUS сервера.
type Option func(*Server)

// Server представляет собой сервер для метрик.
type Server struct {
	Address           string  // Адрес сервера
	CertFile, KeyFile *string // Файлы сертификатов

	registry        *prometheus.Registry // Реестр метрик
	logger          logger.Logger        // Логирование запросов и ошибок сервера
	idleConnsClosed chan struct{}        // Способ определить незавершенные соединения
}

// NewServer создает новый экземпляр сервера для метрик.
func NewServer(cfg *config.Config, options ...Option) *Server {
	srv := &Server{
		Address:         fmt.Sprintf(":%d", cfg.PromListenAddr),
		idleConnsClosed: make(chan struct{}),
	}

	for _, opt := range options {
		opt(srv)
	}

	return srv
}

// setupRouter настраивает роутер.
func (srv *Server) setupRouter() http.Handler {
	handler := http.NewServeMux()

	if srv.registry == nil {
		srv.registry = prometheus.NewRegistry()
	}

	cs := []prometheus.Collector{
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	}

	srv.registry.MustRegister(cs...)

	handler.Handle("/metrics", promhttp.InstrumentMetricHandler(srv.registry, promhttp.HandlerFor(srv.registry, promhttp.HandlerOpts{})))

	return handler
}

// Run запускает сервер.
func (srv *Server) Run(ctx context.Context) error {
	s := http.Server{
		Addr:         srv.Address,
		Handler:      srv.setupRouter(),
		ReadTimeout:  readTimeout,  // wait() + tls handshake + req.headers + req.body
		WriteTimeout: writeTimeout, // все что выше + response
	}

	go func() {
		<-ctx.Done()

		const timeout = 5 * time.Second

		sCtx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()

		defer close(srv.idleConnsClosed)

		srv.logger.Info("PROM сервер остановлен")

		if err := s.Shutdown(sCtx); err != nil {
			srv.logger.Errorf("PROM сервер не остановлен: %v", err)
		}
	}()

	srv.logger.Infof("PROM сервер запущен на %s", srv.Address)

	var err error
	if srv.CertFile != nil && srv.KeyFile != nil {
		err = s.ListenAndServeTLS(*srv.CertFile, *srv.KeyFile)
	} else {
		err = s.ListenAndServe()
	}

	if err != nil {
		srv.Wait()

		return err
	}

	return nil
}

// Wait ожидает момента завершения обработки всех соединений.
func (srv *Server) Wait() {
	<-srv.idleConnsClosed
}
