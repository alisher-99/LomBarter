package storage

import (
	"fmt"
	"strings"

	"gitlab.com/example/gophers/libs/logger"
	"gitlab.com/example/gophers/libs/trace"

	"github.com/alisher-99/LomBarter/internal/config"
	"github.com/alisher-99/LomBarter/internal/domain/repository"
	"github.com/alisher-99/LomBarter/internal/storage/cassandra"
	"github.com/alisher-99/LomBarter/internal/storage/mongo"
)

// ErrInvalidDataStoreName ошибка неверного названия datastore.
type ErrInvalidDataStoreName []string

// Error реализация интерфейса error.
func (ds ErrInvalidDataStoreName) Error() string {
	return fmt.Sprintf("неверное название datastore, доступные: %s", strings.Join(ds, ", "))
}

// dataStoreFactory фабрика для datastore.
type dataStoreFactory func(conf *config.Database, logger logger.Logger, tracer trace.TracerProvider) (repository.DataStore, error)

// newDataStoreFactories создание фабрик для datastore.
func newDataStoreFactories() map[string]dataStoreFactory {
	return map[string]dataStoreFactory{
		"cassandra": cassandra.New,
		"mongo":     mongo.New,
	}
}

// NewDatabase создание нового datastore.
func NewDatabase(conf *config.Database, log logger.Logger, tracer trace.TracerProvider) (repository.DataStore, error) {
	dataStoreFactories := newDataStoreFactories()

	engineFactory, ok := dataStoreFactories[conf.DSName]
	if !ok {
		// выбран неверный datastore, получаем список доступных и отдаем его
		// пользователю
		availableDataStores := make([]string, 0, len(dataStoreFactories)-1)
		for k := range dataStoreFactories {
			availableDataStores = append(availableDataStores, k)
		}

		return nil, ErrInvalidDataStoreName(availableDataStores)
	}

	return engineFactory(conf, log, tracer)
}
