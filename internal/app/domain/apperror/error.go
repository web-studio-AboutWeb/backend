package apperror

type ErrorType int8

const (
	_ ErrorType = iota
	NotFoundType
	InternalType
	UnauthorizedType
	ForbiddenType
	InvalidRequestType
	DisabledType
	DuplicateType
)

type Error struct {
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
	Type    ErrorType `json:"type"`
}

var _ error = Error{}

func (e Error) Error() string {
	return e.Message
}

func New(t ErrorType, msg, field string) *Error {
	return &Error{
		Message: msg,
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
		Type:    NotFoundType,
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
		Type:    DisabledType,
	}
}
