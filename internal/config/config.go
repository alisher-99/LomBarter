package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gitlab.com/example/gophers/libs/kafka/knet"
	envoy "gitlab.com/example/gophers/libs/route-registrator"
)

const (
	productionEnvironment = "production" // Production окружение.

	serviceRoutesYaml = "./config/service-routes.yml" // Путь для регистрации роутов.
)

type (
	// Config конфигурация приложения.
	Config struct {
		Server      `yaml:"server"`
		Log         `yaml:"logger"`
		Database    `yaml:"database"`
		Kafka       `yaml:"kafka"`
		Cache       `yaml:"cache"`
		Tracing     `yaml:"tracing"`
		ServiceMesh `yaml:"service_mesh"`
		Environment `yaml:"environment"`
		ServiceName string `env:"SERVICE_NAME" yaml:"service_name" env-default:"tmp" env-description:"Название сервиса"`
		Version     string `env:"APP_VERSION" yaml:"version" env-default:"unknown" env-description:"Версия приложения"`
	}

	// Environment окружение.
	Environment struct {
		Name string `env:"ENVIRONMENT_NAME" yaml:"name" env-default:"production" env-description:"Название окружения"`
	}

	// Server сервер.
	Server struct {
		Host           string `env:"SERVER_HOST" yaml:"host" env-default:"0.0.0.0" env-description:"Хост HTTP сервиса"`
		HTTPListenAddr int    `env:"SERVER_PORT" yaml:"http_listen_addr" env-default:"8000" env-description:"Адрес HTTP сервера"`
		GrpcListenAddr int    `env:"GRPC_LISTEN" yaml:"grpc_listen_addr" env-default:"4040" env-description:"Адрес GRPC сервера"`
		PromListenAddr int    `env:"PROM_LISTEN" yaml:"prom_listen_addr" env-default:"9090" env-description:"Адрес Prometheus сервера"`
		BasePath       string `env:"BASE_PATH" yaml:"base_path" env-default:"/" env-description:"Базовый путь сервиса"`
		FilesDir       string `env:"FILES_DIR" yaml:"files_dir" env-default:"/swagger" env-description:"Директория с файлами"`
	}

	// Log логирование.
	Log struct {
		LogLevel string `env:"LOG_LEVEL" yaml:"level" env-default:"info" env-description:"Уровень логирования"`
	}

	// Database база данных.
	Database struct {
		DSName string `env:"DATASTORE_NAME" yaml:"name" env-default:"mongo" env-description:"Название БД"`

		// CASSANDRA
		DSPassword string   `env:"DATASTORE_PASSWORD" env-description:"Пароль БД"`
		DSUsername string   `env:"DATASTORE_USER" env-description:"Пользователь БД"`
		DSHosts    []string `env:"DATASTORE_HOSTS" yaml:"hosts" env-description:"Хосты БД"`
		DSKeyspace string   `env:"DATASTORE_KEYSPACE" yaml:"keyspace" env-description:"Пространство БД"`

		// Mongo
		DSDB  string `env:"DATASTORE_DB" yaml:"db" env-description:"DataStore database name (format: fcm)" env-default:"tmp"`
		DSURL string `env:"MONGO_CONTACT_POINTS" yaml:"url" env-required:"true" env-description:"DataStore URL (format: mongodb://localhost:27017)"`
	}

	// Kafka конфигурация Kafka.
	Kafka struct {
		Brokers                    string             `yaml:"brokers" env:"KAFKA_BOOTSTRAP_SERVERS" env-required:"true" env-description:"Брокеры Kafka"`
		Username                   string             `yaml:"username" env:"KAFKA_USERNAME" env-description:"Логин Kafka"`
		Password                   string             `yaml:"password" env:"KAFKA_PASSWORD" env-description:"Пароль Kafka"`
		AuthMechanism              knet.AuthMechanism `yaml:"authMechanism" env:"KAFKA_AUTH_MECHANISM" env-description:"Механизм аутентификации для Kafka"`
		ConsumeTimeout             time.Duration      `yaml:"consumeTimeout" env:"KAFKA_CONSUME_TIMEOUT" env-default:"5s" env-description:"Таймаут при получении сообщений"`
		IsManualCommitAfterProcess bool               `yaml:"isManualCommitAfterProcess" env:"KAFKA_IS_MANUAL_COMMIT_AFTER_PROCESS" env-description:"Режим работы консюмера"`

		Consumers Consumers `yaml:"consumers" env:"KAFKA_CONSUMERS" env-description:"Консюмеры Kafka"`
		Producers Producers `yaml:"producers" env:"KAFKA_PRODUCERS" env-description:"Продюсер Kafka"`
	}

	// Cache конфигурация кэша.
	Cache struct {
		CacheName     string `env:"CACHE_NAME" yaml:"name" env-default:"dragonfly" env-description:"Название кеша"`
		CacheAddr     string `env:"CACHE_ADDR" yaml:"addr" env-description:"Адрес кэша"`
		CacheUsername string `env:"CACHE_USERNAME" yaml:"username" env-description:"Имя пользователя кэша"`
		CachePassword string `env:"CACHE_PASSWORD" yaml:"password" env-description:"Пароль кэша"`
	}

	// Tracing конфигурация трейсинга.
	Tracing struct {
		JaegerEnabled bool   // Включен ли jaeger.
		JaegerURL     string `env:"JAEGER_URL" yaml:"addr" env-default:"" env-description:"Адрес трейсинга"`
	}

	// ServiceMesh конфигурация для envoy service mesh.
	ServiceMesh struct {
		ServiceMeshHost                string `env:"SERVICE_MASH_HOST" yaml:"service_mesh_host" env-default:"127.0.0.1" env-description:"Хост Envoy куда мы отправляем запросы"`
		ServiceMeshPort                int    `env:"SERVICE_MASH_PORT" yaml:"service_mesh_port" env-default:"9200" env-description:"Порт Envoy куда мы отправляем запросы"`
		IsSendEnabledToServiceMesh     bool   `env:"ADV_ENABLED" yaml:"adv_enabled" env-default:"true" env-description:"Активна ли отправка нами запросов в Envoy"`
		SendPeriodToServiceMeshSeconds int    `env:"ADV_PERIOD" yaml:"adv_period" env-default:"7" env-description:"С какой периодичностью мы можем отправлять запрос в Envoy"`
		InternalTmpContainerHost       string `env:"ADV_HOST" yaml:"adv_host" env-default:"127.0.0.1" env-description:"Внутренний хост контейнера микросервиса куда Envoy переотправляет запросы из внешней сети"`
		InternalTmpContainerPort       int    `env:"ADV_PORT" yaml:"adv_port" env-default:"8080" env-description:"Внутренний порт контейнера микросервиса куда Envoy переотправляет запросы из внешней сети"`
		HTTPClientTimeoutSecond        int    `env:"HTTP_CLIENT_TIMEOUT_SECOND" yaml:"http_client_timeout_second" env-default:"5" env-description:"Время ожидания при выполнении HTTP-запросов в секундах"`
	}
)

// NewConfig - конструктор конфига.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("парсинг config.yml: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("парсинг переменных окружения: %w", err)
	}

	cfg.Tracing.defineState()

	return cfg, nil
}

// GetHTTPDomain возвращает домен для HTTP сервера.
func (s Server) GetHTTPDomain() string {
	return fmt.Sprintf("%s:%d", s.Host, s.HTTPListenAddr)
}

// IsProduction является ли прод окружением.
func (e *Environment) IsProduction() bool {
	return e.Name == productionEnvironment
}

// defineState определяет будет ли включен jaeger-a.
func (t *Tracing) defineState() {
	if t.JaegerURL != "" {
		t.JaegerEnabled = true

		return
	}

	t.JaegerEnabled = false
}

// ToEnvoyConfig возвращает envoy конфиг.
func (c *Config) ToEnvoyConfig() *envoy.Config {
	return &envoy.Config{
		ClusterID:               c.ServiceName,
		ServiceMeshHost:         c.ServiceMesh.ServiceMeshHost,
		AdvHost:                 c.ServiceMesh.InternalTmpContainerHost,
		AdvPeriodSecond:         c.ServiceMesh.SendPeriodToServiceMeshSeconds,
		ServiceMeshPort:         c.ServiceMesh.ServiceMeshPort,
		AdvPort:                 c.ServiceMesh.InternalTmpContainerPort,
		HTTPClientTimeoutSecond: c.ServiceMesh.HTTPClientTimeoutSecond,
		AdvEnabled:              c.ServiceMesh.IsSendEnabledToServiceMesh,
		ServiceRoutesYAML:       serviceRoutesYaml,
		Version:                 c.Version,
	}
}
