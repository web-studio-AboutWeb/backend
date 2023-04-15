package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	errcore "web-studio-backend/internal/app/core/shared/errors"
)

func (s *server) sendJSON(code int, data any, w http.ResponseWriter) {
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

func (s *server) readJSON(to any, r *http.Request) error {
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

func (s *server) sendError(err error, w http.ResponseWriter) {
	var coreError *errcore.CoreError

	if errors.As(err, &coreError) {
		var code int

		switch coreError.Type {
		case errcore.NotFoundType:
			code = http.StatusNotFound
		case errcore.InvalidRequestType:
			code = http.StatusBadRequest
		case errcore.UnauthorizedType:
			code = http.StatusUnauthorized
		case errcore.ObjectDuplicateType, errcore.ObjectDisabledType:
			code = http.StatusConflict
		case errcore.AccessDeniedType:
			code = http.StatusForbidden
		default:
			code = http.StatusInternalServerError
		}

		s.sendJSON(code, coreError, w)
		return
	}

	s.sendJSON(http.StatusInternalServerError, errcore.CoreError{Message: err.Error(), Type: errcore.InternalType}, w)
}
