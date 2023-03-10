package errcore

type CoreErrorType string

const (
	NotFoundType        CoreErrorType = "not_found"
	InternalType        CoreErrorType = "internal"
	UnauthorizedType    CoreErrorType = "unauthorized"
	AccessDeniedType    CoreErrorType = "access_denied"
	InvalidRequestType  CoreErrorType = "invalid_request"
	ObjectDisabledType  CoreErrorType = `object_disabled`
	ObjectDuplicateType CoreErrorType = `object_duplicate`
)
