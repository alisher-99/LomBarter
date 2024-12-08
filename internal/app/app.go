package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/example/gophers/libs/logger"
	"golang.org/x/sync/errgroup"

	"github.com/alisher-99/LomBarter/internal/config"
	"github.com/alisher-99/LomBarter/internal/domain/entity"
	"github.com/alisher-99/LomBarter/internal/service"
	"github.com/alisher-99/LomBarter/internal/storage"
	"github.com/alisher-99/LomBarter/internal/transport/http"
)

const (
	fileVersion    = "version" // Файл с версией приложения
	defaultVersion = "unknown" // Версия по умолчанию
)

// gracefulShutdownTimeout представляет собою время, за которое должно сработать корректное завершение программы.
const gracefulShutdownTimeout = 5 * time.Second

// Run запускает приложение.
func Run(cfg *config.Config) error {
	// Инициализация контекста для корректного завершения программы.
	// При получении сигнала SIGTERM или SIGINT, будет вызвана функция cancel.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg.Version = readVersion(fileVersion)

	// Инициализация логгера и graylog.
	log, err := logger.New(cfg.LogLevel, cfg.ServiceName)
	if err != nil {
		return fmt.Errorf("инициализация логгера: %w", err)
	}

	log.WithFields(logger.Fields{
		"version":       cfg.Version,
		"kafka_brokers": cfg.Kafka.Brokers,
		"database_url":  cfg.DSURL,
		"cache_addr":    cfg.CacheAddr,
		"jaeger_url":    cfg.JaegerURL,
	}).Info("Запуск приложения...")

	// Инициализация базы данных.
	ds, err := storage.NewDatabase(&cfg.Database, log, tracer)
	if err != nil {
		return fmt.Errorf("инициализация базы данных: %w", err)
	}

	if err = ds.Connect(); err != nil {
		return fmt.Errorf("подключение к базе данных: %w", err)
	}

	defer func() {
		// Даем таймаут корректному завершению программы
		// Если через определенное время, трейсер не закроется, то закрыть принудительно.
		sCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()

		if cErr := ds.Close(sCtx); cErr != nil {
			log.Errorf("закрытие базы данных: %s", cErr)
		}
	}()

	log.Infof("Подключение к базе данных %s успешно", cfg.DSName)

	// Инициализация сервисов.
	userService := service.NewUserService(ds.UserRepository(), cacheData, log, tracer, producers[entity.SomeTopic], promMetrics)
	orderService := service.NewOrdersService(ds.OrdersRepository(), log, tracer)

	g, gCtx := errgroup.WithContext(ctx)

	// HTTP Сервер.
	g.Go(func() error {
		httpOpts := []http.Option{
			http.WithUserService(userService),
			http.WithOrdersService(orderService),
			http.WithTracer(tracer),
			http.WithLogger(log),
		}

		httpServer := http.NewServer(cfg, httpOpts...)

		return httpServer.Run(gCtx)
	})

	if err = g.Wait(); err != nil {
		return fmt.Errorf("работа основных горутин: %w", err)
	}

	return nil
}
