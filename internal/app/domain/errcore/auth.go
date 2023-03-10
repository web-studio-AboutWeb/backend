package errcore

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
)

var (
	InvalidCredentialsError = &CoreError{
		Message: "Invalid login or password.",
		Code:    "auth.invalid_credentials",
		Type:    UnauthorizedType,
	}
	EmailTakenError = &CoreError{
		Message: "Email already taken.",
		Code:    "auth.email_taken",
		Type:    ObjectDuplicateType,
	}
	UsernameTakenError = &CoreError{
		Message: "Username already taken.",
		Code:    "auth.username_taken",
		Type:    ObjectDuplicateType,
	}
)
