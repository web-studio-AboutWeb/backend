package errcore

var (
	UserNotFoundError = &CoreError{
		Message: "User not found.",
		Code:    "user.not_found",
		Type:    NotFoundType,
	}
	LoginAlreadyTakenError = &CoreError{
		Message: "Login already taken.",
		Code:    "user.login_already_taken",
		Type:    ObjectDuplicateType,
	}
)
