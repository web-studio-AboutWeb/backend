package apperr

type ErrorType int8

const (
	NotFoundType ErrorType = iota + 1
	InternalType
	UnauthorizedType
	ForbiddenType
	InvalidRequestType
	DisabledType
	DuplicateType
)

type (
	ValidationError struct {
		Message string
		Field   string
	}

	Error struct {
		Message          string
		Field            string
		Type             ErrorType
		ValidationErrors []ValidationError
	}
)

var _ error = Error{}

func (e Error) Error() string {
	return e.Message
}

func New(t ErrorType, msg, field string) error {
	return &Error{
		Message: msg,
		Field:   field,
		Type:    t,
	}
}

func NewValidationError(verrors []ValidationError, field string) error {
	return &Error{
		Message:          "Validation error.",
		Type:             InvalidRequestType,
		ValidationErrors: verrors,
		Field:            field,
	}
}

func NewInvalidRequest(msg, field string) error {
	return New(InvalidRequestType, msg, field)
}

func NewNotFound(field string) error {
	return New(NotFoundType, "Object not found.", field)
}

func NewUnauthorized(msg string) error {
	return New(UnauthorizedType, msg, "")
}

func NewForbidden(msg string) error {
	return New(ForbiddenType, msg, "")
}

func NewDisabled(msg string) error {
	return New(DisabledType, msg, "")
}

func NewDuplicate(msg, field string) error {
	return New(DuplicateType, msg, field)
}

func NewInternal(msg string) error {
	return New(InternalType, msg, "")
}
