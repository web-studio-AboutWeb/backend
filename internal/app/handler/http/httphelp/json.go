package httphelp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"web-studio-backend/internal/app/domain/apperror"
)

func SendJSON(code int, data any, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if data == nil {
		return
	}

	buf, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	_, err = w.Write(buf)
	if err != nil {
		panic(err)
	}
}

func ReadJSON(to any, r *http.Request) error {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, to)
	if err != nil {
		return err
	}

	return nil
}

func SendError(err error, w http.ResponseWriter) {
	var coreError *apperror.CoreError

	if errors.As(err, &coreError) {
		var code int

		switch coreError.Type {
		case apperror.NotFoundType:
			code = http.StatusNotFound
		case apperror.InvalidRequestType:
			code = http.StatusBadRequest
		case apperror.UnauthorizedType:
			code = http.StatusUnauthorized
		case apperror.ObjectDuplicateType, apperror.ObjectDisabledType:
			code = http.StatusConflict
		case apperror.AccessDeniedType:
			code = http.StatusForbidden
		default:
			code = http.StatusInternalServerError
		}

		SendJSON(code, coreError, w)
		return
	}

	SendJSON(http.StatusInternalServerError, apperror.CoreError{Message: err.Error(), Type: apperror.InternalType}, w)
}
