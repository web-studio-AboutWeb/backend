package errcore

var (
	ProjectNotFoundError = &CoreError{
		Message: "Project not found.",
		Code:    "project.not_found",
		Type:    NotFoundType,
	}
)
