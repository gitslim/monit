// Модуль retry содержит функции для управления повторными попытками при ошибках.
package retry

import (
	"log"
	"time"
)

// RetryableFunc определяет функцию, которая может быть повторена.
type RetryableFunc func() error

// IRetriableError определяет интерфейс для ошибок, которые можно повторить.
type IRetriableError interface {
	error
	IsRetriable() bool
}

// Retry выполняет функцию с повторными попытками.
func Retry(operation RetryableFunc, maxRetries int) error {
	var err error
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for i := 0; i <= maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		if retriableErr, ok := err.(IRetriableError); ok && retriableErr.IsRetriable() {
			if i < maxRetries {
				log.Printf("Ошибка: %v. Повтор попытки %d через %v...", err, i+1, retryIntervals[i])
				time.Sleep(retryIntervals[i])
			}
		} else {
			return err
		}
	}

	return err
}
