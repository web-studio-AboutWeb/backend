package apperror

import (
	"log/slog"
)

type CoreError struct {
	Message string        `json:"message"`
	Code    string        `json:"code,omitempty"`
	Type    CoreErrorType `json:"type,omitempty"`
}

var _ error = CoreError{}

func (e CoreError) Error() string {
	return e.Message
}

func NewInternalError(err error) *CoreError {
	slog.Error(err.Error())
	return &CoreError{
		Message: err.Error(),
		Type:    InternalType,
	}
}
