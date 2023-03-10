package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (s *server) parseParamString(name string, r *http.Request) string {
	return chi.URLParam(r, name)
}

func (s *server) parseParamInt64(name string, r *http.Request) int64 {
	param := chi.URLParam(r, name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return int64(val)
}

func (s *server) parseParamInt(name string, r *http.Request) int {
	param := chi.URLParam(r, name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return val
}

func (s *server) parseParamInt16(name string, r *http.Request) int16 {
	param := chi.URLParam(r, name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return int16(val)
}
