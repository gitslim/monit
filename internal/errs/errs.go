// Package errs определяет ошибки приложения.
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

// Error определяет сигнальную ошибку.
type Error struct {
	Code    int
	Message string
}

// Error реализацует интерфейс error.
func (e *Error) Error() string {
	return e.Message
}

// NewError создает новую ошибку Error.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
