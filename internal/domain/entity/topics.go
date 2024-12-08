package entity

import (
	"fmt"
)

// Топики для консюмеров.
const (
	// UserUpdateTopic топик для обновления пользователя.
	UserUpdateTopic = "user.update"
)

// Топики для продюсера.
const (
	// SomeTopic тестовый топик для продюсера.
	SomeTopic = "some.topic"
)

// KafkaConfig интерфейс для работы с конфигурацией Kafka.
type KafkaConfig interface {
	// GetTopics возвращает список топиков.
	GetTopics() []string
}

// ValidateConsumerTopics проверяет топики на валидность.
func ValidateConsumerTopics(cfgs KafkaConfig) error {
	return validateTopics(cfgs, []string{UserUpdateTopic})
}

// ValidateProducerTopics проверяет топики на валидность.
func ValidateProducerTopics(cfgs KafkaConfig) error {
	return validateTopics(cfgs, []string{SomeTopic})
}

// validateTopics проверяет, что все переданные топики присутствуют в конфигурации.
func validateTopics(cfgs KafkaConfig, availableTopics []string) error {
	// создаем мапу топиков для быстрой проверки наличия
	cfgTopics := cfgs.GetTopics()

	topicMap := make(map[string]struct{}, len(cfgTopics))
	for _, topic := range cfgTopics {
		topicMap[topic] = struct{}{}
	}

	// Проверяем, что все описанные топики присутствуют в cfgs
	for _, topic := range availableTopics {
		_, ok := topicMap[topic]
		if !ok {
			return fmt.Errorf("%w: %s", ErrTopicNotFound, topic)
		}
	}

	if len(availableTopics) != len(topicMap) {
		return fmt.Errorf("%w: %d!=%d", ErrTopicsLength, len(availableTopics), len(topicMap))
	}

	return nil
}
