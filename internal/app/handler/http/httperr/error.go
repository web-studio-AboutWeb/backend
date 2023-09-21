package httperr

import (
	"errors"
	"net/http"

	"web-studio-backend/internal/app/domain/apperr"
)

type ErrorType string

const (
	ErrorTypeInternal       ErrorType = "INTERNAL"
	ErrorTypeInvalidRequest ErrorType = "INVALID_REQUEST"
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeConflict       ErrorType = "CONFLICT"
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN"
	ErrorTypeValidation     ErrorType = "VALIDATION"
)

type (
	ValidationError struct {
		Field   string `json:"field" validate:"required" example:"title"`
		Message string `json:"message" validate:"required" example:"Title cannot be empty."`
	}

	Error struct {
		Message          string            `json:"message" validate:"required" example:"Project not found."`
		Type             ErrorType         `json:"type" validate:"required"`
		Target           string            `json:"target,omitempty" example:"user_id"` // Target represents error reason object like field
		ValidationErrors []ValidationError `json:"validationErrors,omitempty"`
		HttpCode         int               `json:"-"`
	} //@name Error
)

// ParseAppError parses apperror.Error and returns Error with appropriate fields.
func ParseAppError(ae *apperr.Error) *Error {
	if ae == nil {
		return &Error{
			HttpCode: http.StatusInternalServerError,
			Message:  "Unknown error.",
			Type:     ErrorTypeInternal,
		}
	}

	httpError := Error{
		Message: ae.Message,
		Target:  ae.Field,
	}

	switch ae.Type {
	case apperr.NotFoundType:
		httpError.Type = ErrorTypeNotFound
		httpError.HttpCode = http.StatusNotFound
	case apperr.InvalidRequestType:
		httpError.Type = ErrorTypeInvalidRequest
		httpError.HttpCode = http.StatusBadRequest
	case apperr.UnauthorizedType:
		httpError.Type = ErrorTypeUnauthorized
		httpError.HttpCode = http.StatusUnauthorized
	case apperr.DuplicateType, apperr.DisabledType:
		httpError.Type = ErrorTypeConflict
		httpError.HttpCode = http.StatusConflict
	case apperr.ForbiddenType:
		httpError.Type = ErrorTypeForbidden
		httpError.HttpCode = http.StatusForbidden
	default:
		// TODO: hide message
		httpError.Type = ErrorTypeInternal
		httpError.HttpCode = http.StatusInternalServerError
	}

	if len(ae.ValidationErrors) > 0 {
		httpError.Type = ErrorTypeValidation
		httpError.ValidationErrors = make([]ValidationError, 0, len(ae.ValidationErrors))

		for _, ve := range ae.ValidationErrors {
			httpError.ValidationErrors = append(httpError.ValidationErrors, ValidationError{
				Field:   ve.Field,
				Message: ve.Message,
			})
		}
	}

	return &httpError
}

func UnwrapAppError(err error) *apperr.Error {
	var (
		ae   *apperr.Error
		temp = err
	)

	for {
		if temp == nil {
			return nil
		}

		if errors.As(temp, &ae) {
			return ae
		}

		temp = errors.Unwrap(temp)
	}
}
