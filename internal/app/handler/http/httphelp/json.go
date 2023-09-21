package httphelp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"web-studio-backend/internal/app/handler/http/httperr"
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
		slog.Error("Responding JSON data", slog.String("error", err.Error()))
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

	if appError := httperr.UnwrapAppError(err); appError != nil {
		httpError := httperr.ParseAppError(appError)
		SendJSON(httpError.HttpCode, httpError, w)
		return
	}

	// TODO: hide error message
	SendJSON(http.StatusInternalServerError, httperr.Error{Message: err.Error(), Type: httperr.ErrorTypeInternal}, w)
}
