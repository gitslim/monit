package errs

import "net/http"

var (
	ErrInternal           = NewError(http.StatusInternalServerError, "internal error")
	ErrMetricNotFound     = NewError(http.StatusNotFound, "metric not found")
	ErrInvalidMetricType  = NewError(http.StatusBadRequest, "invalid metric type")
	ErrInvalidMetricValue = NewError(http.StatusBadRequest, "invalid metric value")
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
