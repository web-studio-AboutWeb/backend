package httphelp

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseParamString(name string, r *http.Request) string {
	return chi.URLParam(r, name)
}

func ParseParamInt(name string, r *http.Request) int {
	param := chi.URLParam(r, name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return val
}

func ParseParamInt16(name string, r *http.Request) int16 {
	param := chi.URLParam(r, name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return int16(val)
}

func ParseParamInt32(name string, r *http.Request) int32 {
	param := chi.URLParam(r, name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return int32(val)
}

func ParseParamInt64(name string, r *http.Request) int64 {
	param := chi.URLParam(r, name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return int64(val)
}
