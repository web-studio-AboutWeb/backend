package errcore

var (
	UserNotFoundError = &CoreError{
		Message: "User not found.",
		Code:    "user.not_found",
		Type:    NotFoundType,
	}
)
