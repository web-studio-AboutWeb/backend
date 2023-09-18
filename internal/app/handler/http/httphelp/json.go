package httphelp

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
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
	slog.Error(err.Error())

	if appError := unwrapAppError(err); appError != nil {
		var code int

		switch appError.Type {
		case apperror.NotFoundType:
			code = http.StatusNotFound
		case apperror.InvalidRequestType:
			code = http.StatusBadRequest
		case apperror.UnauthorizedType:
			code = http.StatusUnauthorized
		case apperror.DuplicateType, apperror.DisabledType:
			code = http.StatusConflict
		case apperror.ForbiddenType:
			code = http.StatusForbidden
		default:
			code = http.StatusInternalServerError
		}

		SendJSON(code, appError, w)
		return
	}

	SendJSON(http.StatusInternalServerError, apperror.Error{Message: err.Error(), Type: apperror.InternalType}, w)
}

func unwrapAppError(err error) *apperror.Error {
	var (
		ae   *apperror.Error
		temp = err
	)

	for {
		if temp == nil {
			return nil
		}

		if errors.As(temp, &ae) {
			return ae
		}

		temp = errors.Unwrap(temp)
	}
}
