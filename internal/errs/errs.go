// Модуль errs определяет ошибки приложения.
package errs

import "net/http"

// Сигнальные ошибки приложения.
var (
	ErrInternal           = NewError(http.StatusInternalServerError, "internal error")
	ErrBadRequest         = NewError(http.StatusBadRequest, "bad request")
	ErrMetricNotFound     = NewError(http.StatusNotFound, "metric not found")
	ErrInvalidMetricType  = NewError(http.StatusBadRequest, "invalid metric type")
	ErrInvalidMetricValue = NewError(http.StatusBadRequest, "invalid metric value")
)

// Сигнальная ошибка.
type Error struct {
	Code    int
	Message string
}

// Реализация интерфейса error.
func (e *Error) Error() string {
	return e.Message
}

// Создает новую ошибку Error.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
