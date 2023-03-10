package errcore

import "web-studio-backend/internal/app/infrastructure/logger"

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
	logger.Logger.Err(err)
	return &CoreError{
		Message: err.Error(),
		Type:    InternalType,
	}
}
