package config

import (
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
)

// Duration вспомогательный тип для конфигурации.
type Duration time.Duration

// UnmarshalJSON - установка значения. Необходимо для работы с переменными окружения.
func (c *Duration) UnmarshalJSON(s []byte) error {
	dur := strings.Trim(string(s), "\"")

	d, err := time.ParseDuration(dur)
	if err != nil {
		return err
	}

	*c = Duration(d)

	return nil
}

// UnmarshalYAML - установка значения. Необходимо для работы с переменными окружения.
func (c *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var durationStr string
	if err := unmarshal(&durationStr); err != nil {
		return err
	}

	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return err
	}

	*c = Duration(d)

	return nil
}

// Consumers конфигурация консюмеров.
type Consumers []struct {
	Topic        string `json:"topic" yaml:"topic" env-required:"true" env-description:"Топик Kafka"`
	Group        string `json:"group" yaml:"group" env-description:"Консюмер группа Kafka"`
	AsyncCommits bool   `json:"asyncCommits" yaml:"asyncCommits" env-description:"Асинхронные коммиты сообщений в Kafka"`
}

// GetTopics возвращает список топиков. Реализует интерфейс entity.KafkaConfig для валидации топиков.
func (c *Consumers) GetTopics() []string {
	if c == nil {
		return nil
	}

	topics := make([]string, 0, len(*c))

	for _, consumer := range *c {
		topics = append(topics, consumer.Topic)
	}

	return topics
}

// SetValue - установка значения. Необходимо для работы с переменными окружения через cleanenv.
func (c *Consumers) SetValue(s string) error {
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(s), &c)
	if err != nil {
		return fmt.Errorf("парсинг consumers: %w", err)
	}

	if err = entity.ValidateConsumerTopics(c); err != nil {
		return fmt.Errorf("валидация топиков: %w", err)
	}

	return nil
}

// Producers конфигурация продюсеров.
type Producers []struct {
	Topic                     string   `json:"topic" yaml:"topic" env-required:"true" env-description:"Топик Kafka"`
	NumPartitions             int      `json:"numPartitions" yaml:"numPartitions" env-required:"true" env-description:"Количество партиций"`
	ReplicationFactor         int      `json:"replicationFactor" yaml:"replicationFactor" env-description:"Фактор репликации"`
	Balancer                  string   `json:"balancer" yaml:"balancer" env-description:"Тип балансировщика"`
	Async                     bool     `json:"async" yaml:"async" env-description:"Асинхронный продюсер"`
	BatchBytes                int      `json:"batchBytes" yaml:"batchBytes" env-description:"Максимальный размер батча сообщений"`
	CompressionCodec          string   `json:"compressionCodec" yaml:"compressionCodec" env-description:"Тип компрессии сообщений"`
	DisallowAutoTopicCreation bool     `json:"disallowAutoTopicCreation" yaml:"disallowAutoTopicCreation" env-description:"Запрещать авто-создание топиков"`
	MessageRetention          Duration `json:"messageRetention" yaml:"messageRetention" env-description:"Длительность хранения сообщений в топике"`
}

// GetTopics возвращает список топиков. Реализует интерфейс entity.KafkaConfig для валидации топиков.
func (c *Producers) GetTopics() []string {
	if c == nil {
		return nil
	}

	topics := make([]string, 0, len(*c))

	for _, producer := range *c {
		topics = append(topics, producer.Topic)
	}

	return topics
}

// SetValue - установка значения. Необходимо для работы с переменными окружения через cleanenv.
func (c *Producers) SetValue(s string) error {
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(s), &c)
	if err != nil {
		return fmt.Errorf("парсинг producers: %w", err)
	}

	if err = entity.ValidateProducerTopics(c); err != nil {
		return fmt.Errorf("валидация топиков: %w", err)
	}

	return nil
}
