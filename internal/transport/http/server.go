package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	authMiddleware "gitlab.com/example/gophers/libs/auth/middleware"
	"gitlab.com/example/gophers/libs/logger"
	loggerMiddleware "gitlab.com/example/gophers/libs/logger/middleware"
	multiLangMiddleware "gitlab.com/example/gophers/libs/multi-lang/middleware/http"
	"gitlab.com/example/gophers/libs/trace"
	traceMiddleware "gitlab.com/example/gophers/libs/trace/middleware/http"

	"github.com/alisher-99/LomBarter/internal/config"
	"github.com/alisher-99/LomBarter/internal/service"
	"github.com/alisher-99/LomBarter/internal/transport/http/resources"
	v1 "github.com/alisher-99/LomBarter/internal/transport/http/resources/v1"
)

const (
	compressLevel = 5                // Уровень сжатия.
	readTimeout   = 5 * time.Second  // Время ожидания запроса.
	writeTimeout  = 30 * time.Second // Время ожидания ответа.
	maxAge        = 300              // Время жизни C.O.R.S. заголовков.
)

// Server представляет собой HTTP сервер.
type Server struct {
	Address     string             // Адрес сервера
	BasePath    string             // Базовый путь
	FilesDir    string             // Директория с файлами
	Environment config.Environment // Окружение

	logger          logger.Logger        // Логирование запросов и ошибок сервера
	tracer          trace.TracerProvider // Отслеживает запросы между слоями и микросервисами
	idleConnsClosed chan struct{}        // Способ определить незавершенные соединения
	version         string               // Версия приложения

	userService   service.UserService   // Сервис пользователей
	ordersService service.OrdersService // Сервис заказов
}

// NewServer создает новый HTTP сервер.
func NewServer(cfg *config.Config, options ...Option) *Server {
	srv := &Server{
		Address:     cfg.GetHTTPDomain(),
		FilesDir:    cfg.FilesDir,
		BasePath:    cfg.BasePath,
		Environment: cfg.Environment,

		idleConnsClosed: make(chan struct{}),
		version:         cfg.Version,
	}

	for _, opt := range options {
		opt(srv)
	}

	return srv
}

// setupRouter инициализирует HTTP роутер. Функция используется для подключения middleware и маппинга ресурсов.
func (srv *Server) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)                  // no-cache
	r.Use(middleware.RequestID)                // вставляет request ID в контекст каждого запроса
	r.Use(authMiddleware.ParseUserHeaders)     // добавляет кастомные заголовки
	r.Use(loggerMiddleware.Logger(srv.logger)) // логирует начало и окончание каждого запроса с указанием времени обработки
	r.Use(middleware.Recoverer)                // управляемо обрабатывает паники и выдает stack trace при их возникновении
	r.Use(middleware.RealIP)                   // устанавливает RemoteAddr для каждого запроса с заголовками X-Forwarded-For или X-Real-IP
	r.Use(middleware.NewCompressor(compressLevel).Handler)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(srv.Environment),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           maxAge, // Максимальное время жизни C.O.R.S. заголовков.
	}))
	r.Use(multiLangMiddleware.LanguageMiddleware)
	r.Use(multiLangMiddleware.SetLanguageResponseHeader)

	tm := traceMiddleware.New(srv.tracer)
	r.Use(tm.OpenTelemetryMiddleware)

	// монтируем дополнительные ресурсы
	r.Mount("/version", resources.VersionResource{Version: srv.version}.Routes())
	r.Mount("/api/v1/users", v1.NewUserHandler(srv.userService, srv.logger).Routes())
	r.Mount("/api/v1/orders", v1.NewOrdersHandler(srv.ordersService, srv.logger).Routes())

	if !srv.Environment.IsProduction() {
		r.Mount("/files", resources.FilesResource{FilesDir: srv.FilesDir}.Routes())
		r.Mount("/swagger", resources.SwaggerResource{FilesPath: "/files", BasePath: srv.BasePath}.Routes())
	}

	return r
}

// getAllowedOrigins возвращает список хостов для C.O.R.S.
func allowedOrigins(environment config.Environment) []string {
	if environment.IsProduction() {
		return []string{"*"}
	}

	return []string{"*"}
}

// Run запускает HTTP или HTTPS листенер в зависимости от того как заполнена
// структура Server{}.
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

		srv.logger.Info("HTTP сервер остановлен")

		if err := s.Shutdown(sCtx); err != nil {
			srv.logger.Errorf("HTTP сервер не остановлен: %v", err)
		}
	}()

	srv.logger.Infof("HTTP сервер запущен на %s", srv.Address)

	if err := s.ListenAndServe(); err != nil {
		srv.Wait()

		return err
	}

	return nil
}

// Wait ожидает момента завершения обработки всех соединений.
func (srv *Server) Wait() {
	<-srv.idleConnsClosed
}
