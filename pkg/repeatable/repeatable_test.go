package repeatable

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoWithTries_SuccessOnFirstAttempt(t *testing.T) {
	t.Parallel()

	attempts := 3
	delay := time.Millisecond

	// Мок функция, которая всегда успешна
	fn := func() error {
		return nil
	}

	err := DoWithTries(fn, attempts, delay)
	assert.NoError(t, err)
}

func TestDoWithTries_SuccessAfterRetries(t *testing.T) {
	t.Parallel()

	attempts := 3
	delay := time.Millisecond

	// Мок функция, которая сначала возвращает ошибку 2 раза, а потом успешна
	var count int

	fn := func() error {
		if count < 2 {
			count++

			return assert.AnError
		}

		return nil
	}

	err := DoWithTries(fn, attempts, delay)
	assert.NoError(t, err)
}

func TestDoWithTries_AllAttemptsFail(t *testing.T) {
	t.Parallel()

	attempts := 3
	delay := time.Millisecond

	// Мок функция, которая всегда возвращает ошибку
	fn := func() error {
		return assert.AnError
	}

	err := DoWithTries(fn, attempts, delay)
	assert.Equal(t, assert.AnError, err)
}
