package http

import (
	"net/http"
	"web-studio-backend/internal/app/domain"
)

func (s *server) SignIn(w http.ResponseWriter, r *http.Request) {
	var req domain.SignInRequest
	if err := s.readJSON(&req, r); err != nil {
		s.sendError(err, w)
		return
	}

	response, err := s.core.SignIn(r.Context(), &req)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

func (s *server) SignUp(w http.ResponseWriter, r *http.Request) {
	var req domain.SignUpRequest
	if err := s.readJSON(&req, r); err != nil {
		s.sendError(err, w)
		return
	}

	response, err := s.core.SignUp(r.Context(), &req)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}
