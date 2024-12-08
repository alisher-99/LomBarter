package mongo

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/example/gophers/libs/logger"
	"gitlab.com/example/gophers/libs/trace"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/alisher-99/LomBarter/internal/config"
	"github.com/alisher-99/LomBarter/internal/domain/entity"
	"github.com/alisher-99/LomBarter/internal/domain/repository"
)

const (
	// connectionTimeout время ожидания подключения к MongoDB.
	connectionTimeout = 3 * time.Second
	// ensureIdxTimeout время ожидания создания индексов.
	ensureIdxTimeout = 10 * time.Second
)

const (
	// tracerName название трейса.
	tracerName = "mongo"

	// userCollection коллекция пользователей.
	userCollection = "user"
	// ordersCollection коллекция заказов.
	ordersCollection = "orders"
)

// Mongo реализация DataStore для MongoDB.
type Mongo struct {
	connURL string               // URL подключения к базе данных
	dbName  string               // Название базы данных
	logger  logger.Logger        // Логирование запросов и ошибок базы
	tracer  trace.TracerProvider // Отслеживает запросы между слоями и микросервисами

	client *mongo.Client   // Клиент для работы с MongoDB
	DB     *mongo.Database // База данных

	connectionTimeout time.Duration // Время ожидания подключения к MongoDB
	ensureIdxTimeout  time.Duration // Время ожидания создания индексов

	userRepo   repository.UserRepository   // Репозиторий пользователей
	ordersRepo repository.OrdersRepository // Репозиторий заказов
}

// Name возвращает название DataStore.
func (m *Mongo) Name() string { return "mongo" }

// New создание нового datastore.
func New(conf *config.Database, log logger.Logger, tracer trace.TracerProvider) (repository.DataStore, error) {
	if conf.DSURL == "" {
		return nil, entity.ErrInvalidDatabaseURL
	}

	if conf.DSDB == "" {
		return nil, entity.ErrInvalidDatabaseName
	}

	return &Mongo{
		connURL:           conf.DSURL,
		dbName:            conf.DSDB,
		logger:            log,
		tracer:            tracer,
		connectionTimeout: connectionTimeout,
		ensureIdxTimeout:  ensureIdxTimeout,
	}, nil
}

// Connect подключение к MongoDB.
func (m *Mongo) Connect() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.connURL))
	if err != nil {
		return fmt.Errorf("коннект к MongoDB: %w", err)
	}

	if err = m.Ping(); err != nil {
		return fmt.Errorf("пинг MongoDB: %w", err)
	}

	m.DB = m.client.Database(m.dbName)

	// убеждаемся что созданы все необходимые индексы.
	return m.ensureIndexes()
}

// Ping проверяет что соединение с MongoDB установлено.
func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
}

// Close закрывает соединение с MongoDB.
func (m *Mongo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// UserRepository возвращает репозиторий пользователей.
func (m *Mongo) UserRepository() repository.UserRepository {
	if m.userRepo == nil {
		m.userRepo = NewUserRepository(m.DB.Collection(userCollection), m.tracer)
	}

	return m.userRepo
}

// OrdersRepository возвращает репозиторий заказов.
func (m *Mongo) OrdersRepository() repository.OrdersRepository {
	if m.ordersRepo == nil {
		m.ordersRepo = NewOrdersRepository(m.DB.Collection(ordersCollection), m.tracer)
	}

	return m.ordersRepo
}

// ensureIndexes убеждается что все индексы построены.
func (m *Mongo) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	if err := m.ensureUserIndexes(ctx); err != nil {
		return fmt.Errorf("построение индексов для пользователей: %w", err)
	}

	if err := m.ensureOrdersIndexes(ctx); err != nil {
		return fmt.Errorf("построение индексов для заказов: %w", err)
	}

	return nil
}

// ensureUserIndexes убеждается что все индексы построены для коллекции пользователей.
func (m *Mongo) ensureUserIndexes(_ context.Context) error {
	return nil
}

// ensureOrdersIndexes убеждается что все индексы построены для коллекции заказов.
func (m *Mongo) ensureOrdersIndexes(_ context.Context) error {
	return nil
}

// StartSession создает сессию для транзакции.
func (m *Mongo) StartSession(ctx context.Context) (context.Context, repository.TxCallback, error) {
	wc := writeconcern.Majority()
	rc := readconcern.Snapshot()
	txOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := m.client.StartSession()
	if err != nil {
		return nil, nil, fmt.Errorf("начало сессии: %w", err)
	}

	if err = session.StartTransaction(txOpts); err != nil {
		return nil, nil, fmt.Errorf("начало транзакции: %w", err)
	}

	return mongo.NewSessionContext(ctx, session), callback(session), nil
}

// callback для отката или коммита транзакции.
func callback(session mongo.Session) func(ctx context.Context, err error) error {
	return func(ctx context.Context, err error) error {
		defer session.EndSession(ctx)

		if err == nil {
			err = session.CommitTransaction(ctx)
		}

		if err != nil {
			if abortErr := session.AbortTransaction(ctx); abortErr != nil {
				err = fmt.Errorf("%w, tx abortErr: %s", err, abortErr.Error())
			}
		}

		return err
	}
}
