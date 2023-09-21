package apperror

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

func New(t ErrorType, msg, field string) *Error {
	return &Error{
		Message: msg,
		Field:   field,
		Type:    t,
	}
}

func NewInvalidRequest(msg, field string) error {
	return &Error{
		Message: msg,
		Field:   field,
		Type:    InvalidRequestType,
	}
}

func NewNotFound(field string) error {
	return &Error{
		Message: "Object not found.",
		Field:   field,
		Type:    NotFoundType,
	}
}

func NewUnauthorized(msg string) error {
	return &Error{
		Message: msg,
		Type:    UnauthorizedType,
	}
}

func NewForbidden(msg string) error {
	return &Error{
		Message: msg,
		Type:    ForbiddenType,
	}
}

func NewDisabled(msg string) error {
	return &Error{
		Message: msg,
		Type:    DisabledType,
	}
}

func NewDuplicate(msg, field string) error {
	return &Error{
		Message: msg,
		Field:   field,
		Type:    DuplicateType,
	}
}
