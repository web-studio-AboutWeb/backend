package user

import (
	"web-studio-backend/internal/app/core/shared/errors"
)

var (
	UserNotFoundError = &errors.CoreError{
		Message: "User not found.",
		Code:    "user.not_found",
		Type:    errors.NotFoundType,
	}
	LoginAlreadyTakenError = &errors.CoreError{
		Message: "Login already taken.",
		Code:    "user.login_already_taken",
		Type:    errors.ObjectDuplicateType,
	}
)
