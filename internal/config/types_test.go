package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumers_GetTopics(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		consumers Consumers
		expRes    []string
	}{
		{
			name: "возвращаем список топиков",
			consumers: Consumers{
				{Topic: "topic1"},
				{Topic: "topic2"},
			},
			expRes: []string{"topic1", "topic2"},
		},
		{
			name:      "возвращаем пустой список топиков",
			consumers: nil,
			expRes:    []string{},
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			res := s.consumers.GetTopics()
			require.Equal(t, s.expRes, res)
		})
	}
}

func TestConsumers_SetValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		conStr string
		expErr string
		expRes Consumers
	}{
		{
			name:   "устанавливаем корректное значение",
			conStr: `[{"topic": "user.update"}]`,
			expErr: "",
			expRes: Consumers{
				{Topic: "user.update"},
			},
		},
		{
			name:   "устанавливаем некорректное значение",
			conStr: `[{"topic": "topic1"}, {"topic": "topic2"}`,
			expErr: "парсинг consumers: config.Consumers: decode slice: expect ], but found \x00, error found in #10 byte of ...| \"topic2\"}|..., bigger context ...|[{\"topic\": \"topic1\"}, {\"topic\": \"topic2\"}|...",
			expRes: nil,
		},
		{
			name:   "устанавливаем некорректный топик",
			conStr: `[{"topic": "topic1"}, {"topic": "topic2"}]`,
			expErr: "валидация топиков: топик не найден: user.update",
			expRes: nil,
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			var consumers Consumers
			err := consumers.SetValue(s.conStr)
			if s.expErr != "" {
				assert.EqualError(t, err, s.expErr)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, s.expRes, consumers)
		})
	}
}

func TestProducers_GetTopics(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		producers Producers
		expRes    []string
	}{
		{
			name: "возвращаем список топиков",
			producers: Producers{
				{Topic: "topic1"},
				{Topic: "topic2"},
			},
			expRes: []string{"topic1", "topic2"},
		},
		{
			name:      "возвращаем пустой список топиков",
			producers: nil,
			expRes:    []string{},
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			res := s.producers.GetTopics()
			require.Equal(t, s.expRes, res)
		})
	}
}

func TestProducers_SetValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		prodStr string
		expErr  string
		expRes  Producers
	}{
		{
			name:    "устанавливаем корректное значение",
			prodStr: `[{"topic": "some.topic"}]`,
			expErr:  "",
			expRes: Producers{
				{Topic: "some.topic"},
			},
		},
		{
			name:    "устанавливаем некорректное значение",
			prodStr: `[{"topic": "topic1"}, {"topic": "topic2"}`,
			expErr:  "парсинг producers: config.Producers: decode slice: expect ], but found \x00, error found in #10 byte of ...| \"topic2\"}|..., bigger context ...|[{\"topic\": \"topic1\"}, {\"topic\": \"topic2\"}|...",
			expRes:  nil,
		},
		{
			name:    "устанавливаем некорректный топик",
			prodStr: `[{"topic": "topic1"}, {"topic": "topic2"}]`,
			expErr:  "валидация топиков: топик не найден: some.topic",
			expRes:  nil,
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			var producers Producers
			err := producers.SetValue(s.prodStr)
			if s.expErr != "" {
				assert.EqualError(t, err, s.expErr)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, s.expRes, producers)
		})
	}
}

func TestDuration_SetValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		durStr string
		expErr string
		expDur time.Duration
	}{
		{
			name:   "устанавливаем корректное значение",
			durStr: "3s",
			expErr: "",
			expDur: 3 * time.Second,
		},
		{
			name:   "устанавливаем корректное значение с двойными кавычками",
			durStr: "\"3s\"",
			expErr: "",
			expDur: 3 * time.Second,
		},
		{
			name:   "устанавливаем некорректное значение",
			durStr: "3",
			expErr: "time: missing unit in duration \"3\"",
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			var dur Duration
			err := dur.UnmarshalJSON([]byte(s.durStr))
			if s.expErr != "" {
				assert.EqualError(t, err, s.expErr)

				return
			}

			require.NoError(t, err)
			require.Equal(t, s.expDur, time.Duration(dur))
		})
	}
}
